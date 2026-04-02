package utils

import (
	"archive/zip"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	_ "image/jpeg"
	_ "image/png"
)

const (
	deckImageMaxBytes      = 2 * 1024 * 1024
	deckAspectRatioTarget  = 2.0 / 3.0
	deckAspectRatioAllowed = 0.05
)

var minorSuitFolders = map[string]struct{}{
	"wands":     {},
	"cups":      {},
	"swords":    {},
	"pentacles": {},
}

type DeckUploadManifest struct {
	CardKeys     []string
	CardZipPaths map[string]string
}

func ValidateAndExtractDeckZIP(zipPath, destPath string, expectedKeys []string) (*DeckUploadManifest, error) {
	manifest, err := ValidateDeckZIP(zipPath, expectedKeys)
	if err != nil {
		return nil, err
	}

	if err := ExtractDeckZIP(zipPath, destPath, manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}

func ValidateDeckZIP(zipPath string, expectedKeys []string) (*DeckUploadManifest, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}
	defer reader.Close()

	expectedSet := make(map[string]struct{}, len(expectedKeys))
	for _, key := range expectedKeys {
		expectedSet[key] = struct{}{}
	}

	cardZipPaths := make(map[string]string, len(expectedKeys))
	majorCount := 0
	minorCounts := map[string]int{
		"wands":     0,
		"cups":      0,
		"swords":    0,
		"pentacles": 0,
	}

	for _, file := range reader.File {
		cleanedPath, skip, err := normalizeZipEntryPath(file)
		if err != nil {
			return nil, err
		}
		if skip {
			continue
		}

		_, cardKey, err := classifyDeckZipEntry(cleanedPath, expectedSet)
		if err != nil {
			return nil, err
		}

		if err := validateZipImage(file); err != nil {
			return nil, fmt.Errorf("invalid card image %s: %w", cleanedPath, err)
		}
		if _, exists := cardZipPaths[cardKey]; exists {
			return nil, fmt.Errorf("duplicate card image for key %s", cardKey)
		}
		cardZipPaths[cardKey] = cleanedPath

		if strings.HasPrefix(cleanedPath, "major/") {
			majorCount++
			continue
		}

		parts := strings.Split(cleanedPath, "/")
		if len(parts) == 3 {
			minorCounts[parts[1]]++
		}
	}

	if majorCount != 22 {
		return nil, fmt.Errorf("major arcana must contain exactly 22 cards")
	}
	for suit, count := range minorCounts {
		if count != 14 {
			return nil, fmt.Errorf("minor/%s must contain exactly 14 cards", suit)
		}
	}
	if len(cardZipPaths) != len(expectedSet) {
		return nil, fmt.Errorf("zip must contain exactly %d unique card images", len(expectedSet))
	}

	missingKeys := make([]string, 0)
	for key := range expectedSet {
		if _, ok := cardZipPaths[key]; !ok {
			missingKeys = append(missingKeys, key)
		}
	}
	if len(missingKeys) > 0 {
		sort.Strings(missingKeys)
		return nil, fmt.Errorf("missing card images: %s", strings.Join(missingKeys[:min(5, len(missingKeys))], ", "))
	}

	keys := make([]string, 0, len(cardZipPaths))
	for key := range cardZipPaths {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return &DeckUploadManifest{
		CardKeys:     keys,
		CardZipPaths: cardZipPaths,
	}, nil
}

func ExtractDeckZIP(zipPath, destPath string, manifest *DeckUploadManifest) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}
	defer reader.Close()

	if err := os.RemoveAll(destPath); err != nil {
		return fmt.Errorf("reset destination: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(destPath, "cards"), 0o755); err != nil {
		return fmt.Errorf("prepare destination: %w", err)
	}

	entries := make(map[string]*zip.File, len(reader.File))
	for _, file := range reader.File {
		cleanedPath, skip, err := normalizeZipEntryPath(file)
		if err != nil {
			return err
		}
		if skip {
			continue
		}
		entries[cleanedPath] = file
	}

	for _, key := range manifest.CardKeys {
		sourcePath, ok := manifest.CardZipPaths[key]
		if !ok {
			_ = os.RemoveAll(destPath)
			return fmt.Errorf("card %s missing from manifest", key)
		}
		if err := extractZipImageAsPNG(entries[sourcePath], filepath.Join(destPath, "cards", key+".png")); err != nil {
			_ = os.RemoveAll(destPath)
			return fmt.Errorf("extract card %s: %w", key, err)
		}
	}

	return nil
}

func normalizeZipEntryPath(file *zip.File) (string, bool, error) {
	name := strings.ReplaceAll(file.Name, "\\", "/")
	cleaned := path.Clean(strings.TrimSpace(name))

	if cleaned == "." || cleaned == "" {
		return "", true, nil
	}
	if strings.HasPrefix(cleaned, "../") || strings.Contains(cleaned, "/../") || path.IsAbs(cleaned) {
		return "", false, fmt.Errorf("zip contains invalid path: %s", file.Name)
	}
	if file.FileInfo().IsDir() {
		return "", true, nil
	}
	if isIgnorableZipEntry(cleaned) {
		return "", true, nil
	}

	return cleaned, false, nil
}

func isIgnorableZipEntry(cleaned string) bool {
	base := path.Base(cleaned)
	return strings.HasPrefix(cleaned, "__MACOSX/") || base == ".DS_Store"
}

func classifyDeckZipEntry(cleanedPath string, expectedKeys map[string]struct{}) (string, string, error) {
	parts := strings.Split(cleanedPath, "/")
	base := strings.ToLower(path.Base(cleanedPath))

	switch {
	case len(parts) == 2 && parts[0] == "major":
		key, err := normalizeDeckCardKey(strings.TrimSuffix(base, path.Ext(base)), "", expectedKeys)
		return "card", key, err
	case len(parts) == 3 && parts[0] == "minor":
		suit := parts[1]
		if _, ok := minorSuitFolders[suit]; !ok {
			return "", "", fmt.Errorf("unexpected minor arcana suit folder: %s", suit)
		}
		key, err := normalizeDeckCardKey(strings.TrimSuffix(base, path.Ext(base)), suit, expectedKeys)
		return "card", key, err
	default:
		return "", "", fmt.Errorf("unexpected file structure in deck zip: %s", cleanedPath)
	}
}

func normalizeDeckCardKey(rawName, suit string, expectedKeys map[string]struct{}) (string, error) {
	name := strings.ToLower(strings.TrimSpace(rawName))
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, " ", "_")

	for len(name) > 0 && name[0] >= '0' && name[0] <= '9' {
		name = name[1:]
	}
	name = strings.TrimLeft(name, "_")
	name = collapseUnderscores(name)

	if _, ok := expectedKeys[name]; ok {
		return name, nil
	}

	if suit != "" && !strings.Contains(name, "_of_") {
		withSuit := collapseUnderscores(name + "_of_" + suit)
		if _, ok := expectedKeys[withSuit]; ok {
			return withSuit, nil
		}
	}

	return "", fmt.Errorf("file name %q does not map to a known tarot card key", rawName)
}

func collapseUnderscores(value string) string {
	for strings.Contains(value, "__") {
		value = strings.ReplaceAll(value, "__", "_")
	}
	return strings.Trim(value, "_")
}

func validateZipImage(file *zip.File) error {
	ext := strings.ToLower(path.Ext(file.Name))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return fmt.Errorf("unsupported file format")
	}
	if file.UncompressedSize64 > deckImageMaxBytes {
		return fmt.Errorf("file exceeds 2MB limit")
	}

	reader, err := file.Open()
	if err != nil {
		return fmt.Errorf("open image: %w", err)
	}
	defer reader.Close()

	return validateImageReader(reader)
}

func ValidateAndSaveBackImage(sourcePath, destPath string) error {
	if err := validateImagePath(sourcePath); err != nil {
		return fmt.Errorf("invalid back image: %w", err)
	}

	if err := saveImageAsPNG(sourcePath, destPath); err != nil {
		return fmt.Errorf("save back image: %w", err)
	}

	return nil
}

func validateImagePath(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("stat image: %w", err)
	}
	if info.Size() > deckImageMaxBytes {
		return fmt.Errorf("file exceeds 2MB limit")
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return fmt.Errorf("unsupported file format")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open image: %w", err)
	}
	defer file.Close()

	return validateImageReader(file)
}

func validateImageReader(reader io.Reader) error {
	config, format, err := image.DecodeConfig(reader)
	if err != nil {
		return fmt.Errorf("decode image config: %w", err)
	}
	if format != "png" && format != "jpeg" {
		return fmt.Errorf("unsupported image format")
	}
	if config.Width <= 0 || config.Height <= 0 {
		return fmt.Errorf("invalid image dimensions")
	}

	ratio := float64(config.Width) / float64(config.Height)
	diff := ratio / deckAspectRatioTarget
	if diff < 1-deckAspectRatioAllowed || diff > 1+deckAspectRatioAllowed {
		return fmt.Errorf("image aspect ratio must be close to 2:3")
	}

	return nil
}

func saveImageAsPNG(sourcePath, destPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return err
	}

	output, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer output.Close()

	return png.Encode(output, img)
}

func extractZipImageAsPNG(file *zip.File, destPath string) error {
	if file == nil {
		return fmt.Errorf("zip entry missing")
	}

	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return err
	}

	output, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer output.Close()

	if err := png.Encode(output, img); err != nil {
		return err
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
