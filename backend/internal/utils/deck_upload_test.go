package utils

import (
	"archive/zip"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var testMajorKeys = []string{
	"the_fool", "the_magician", "the_high_priestess", "the_empress", "the_emperor", "the_hierophant",
	"the_lovers", "the_chariot", "strength", "the_hermit", "wheel_of_fortune", "justice",
	"the_hanged_man", "death", "temperance", "the_devil", "the_tower", "the_star",
	"the_moon", "the_sun", "judgement", "the_world",
}

var testMinorSuitKeys = map[string][]string{
	"wands": {
		"ace_of_wands", "two_of_wands", "three_of_wands", "four_of_wands", "five_of_wands", "six_of_wands", "seven_of_wands",
		"eight_of_wands", "nine_of_wands", "ten_of_wands", "page_of_wands", "knight_of_wands", "queen_of_wands", "king_of_wands",
	},
	"cups": {
		"ace_of_cups", "two_of_cups", "three_of_cups", "four_of_cups", "five_of_cups", "six_of_cups", "seven_of_cups",
		"eight_of_cups", "nine_of_cups", "ten_of_cups", "page_of_cups", "knight_of_cups", "queen_of_cups", "king_of_cups",
	},
	"swords": {
		"ace_of_swords", "two_of_swords", "three_of_swords", "four_of_swords", "five_of_swords", "six_of_swords", "seven_of_swords",
		"eight_of_swords", "nine_of_swords", "ten_of_swords", "page_of_swords", "knight_of_swords", "queen_of_swords", "king_of_swords",
	},
	"pentacles": {
		"ace_of_pentacles", "two_of_pentacles", "three_of_pentacles", "four_of_pentacles", "five_of_pentacles", "six_of_pentacles", "seven_of_pentacles",
		"eight_of_pentacles", "nine_of_pentacles", "ten_of_pentacles", "page_of_pentacles", "knight_of_pentacles", "queen_of_pentacles", "king_of_pentacles",
	},
}

func TestValidateAndExtractDeckZIP(t *testing.T) {
	t.Parallel()

	expectedKeys := buildExpectedTestKeys()
	zipPath := filepath.Join(t.TempDir(), "valid-deck.zip")
	if err := writeTestDeckZip(zipPath, false); err != nil {
		t.Fatalf("write zip: %v", err)
	}

	destDir := filepath.Join(t.TempDir(), "deck")
	manifest, err := ValidateAndExtractDeckZIP(zipPath, destDir, expectedKeys)
	if err != nil {
		t.Fatalf("ValidateAndExtractDeckZIP returned error: %v", err)
	}

	if len(manifest.CardKeys) != 78 {
		t.Fatalf("expected 78 card keys, got %d", len(manifest.CardKeys))
	}

	if _, err := os.Stat(filepath.Join(destDir, "cards", "the_fool.png")); err != nil {
		t.Fatalf("expected extracted card image: %v", err)
	}
	if _, err := os.Stat(filepath.Join(destDir, "cards", "king_of_pentacles.png")); err != nil {
		t.Fatalf("expected extracted final card image: %v", err)
	}
}

func TestValidateDeckZIPFailsOnMissingCard(t *testing.T) {
	t.Parallel()

	expectedKeys := buildExpectedTestKeys()
	zipPath := filepath.Join(t.TempDir(), "invalid-deck.zip")
	if err := writeTestDeckZip(zipPath, true); err != nil {
		t.Fatalf("write zip: %v", err)
	}

	_, err := ValidateDeckZIP(zipPath, expectedKeys)
	if err == nil {
		t.Fatal("expected validation error for missing card")
	}
	if !strings.Contains(err.Error(), "missing card images") && !strings.Contains(err.Error(), "exactly 78 unique card images") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateDeckZIPRejectsNumericCardIDs(t *testing.T) {
	t.Parallel()

	expectedKeys := buildExpectedTestKeys()
	zipPath := filepath.Join(t.TempDir(), "numeric-deck.zip")
	if err := writeNumericTestDeckZip(zipPath); err != nil {
		t.Fatalf("write zip: %v", err)
	}

	_, err := ValidateDeckZIP(zipPath, expectedKeys)
	if err == nil {
		t.Fatal("expected validation error for numeric card ids")
	}
	if !strings.Contains(err.Error(), "does not map to a known tarot card key") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateAndSaveBackImage(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "back.jpg")
	if err := writeImageFile(sourcePath); err != nil {
		t.Fatalf("write back image: %v", err)
	}

	destPath := filepath.Join(tempDir, "deck", "back.png")
	if err := ValidateAndSaveBackImage(sourcePath, destPath); err != nil {
		t.Fatalf("ValidateAndSaveBackImage returned error: %v", err)
	}

	if _, err := os.Stat(destPath); err != nil {
		t.Fatalf("expected saved back image: %v", err)
	}
}

func buildExpectedTestKeys() []string {
	keys := make([]string, 0, 78)
	keys = append(keys, testMajorKeys...)
	for _, suit := range []string{"wands", "cups", "swords", "pentacles"} {
		keys = append(keys, testMinorSuitKeys[suit]...)
	}
	return keys
}

func writeTestDeckZip(zipPath string, omitLastMinor bool) error {
	file, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := zip.NewWriter(file)

	for index, key := range testMajorKeys {
		name := strings.ReplaceAll(key, "_", "-")
		path := filepath.ToSlash(filepath.Join("major", formatMajorName(index, name)+".png"))
		if err := addPNGToZip(writer, path); err != nil {
			return err
		}
	}

	for _, suit := range []string{"wands", "cups", "swords", "pentacles"} {
		cards := testMinorSuitKeys[suit]
		limit := len(cards)
		if omitLastMinor && suit == "pentacles" {
			limit--
		}
		for _, key := range cards[:limit] {
			base := strings.ReplaceAll(key, "_", "-")
			if strings.HasSuffix(base, "-of-"+suit) {
				base = strings.TrimSuffix(base, "-of-"+suit)
			}
			path := filepath.ToSlash(filepath.Join("minor", suit, base+".png"))
			if err := addPNGToZip(writer, path); err != nil {
				return err
			}
		}
	}

	return writer.Close()
}

func writeNumericTestDeckZip(zipPath string) error {
	file, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := zip.NewWriter(file)

	for id := 22; id <= 121; id++ {
		entry, err := writer.Create(filepath.ToSlash(filepath.Join("card", fmt.Sprintf("%d.png", id))))
		if err != nil {
			return err
		}
		if err := writeImage(entry); err != nil {
			return err
		}
	}

	return writer.Close()
}

func formatMajorName(index int, name string) string {
	return fmt.Sprintf("%02d-%s", index, name)
}

func addPNGToZip(writer *zip.Writer, name string) error {
	entry, err := writer.Create(name)
	if err != nil {
		return err
	}

	return writeImage(entry)
}

func writeImageFile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return writeImage(file)
}

func writeImage(writer io.Writer) error {
	img := image.NewRGBA(image.Rect(0, 0, 200, 300))
	fill := color.RGBA{R: 212, G: 175, B: 55, A: 255}
	for y := 0; y < 300; y++ {
		for x := 0; x < 200; x++ {
			img.Set(x, y, fill)
		}
	}

	return png.Encode(writer, img)
}
