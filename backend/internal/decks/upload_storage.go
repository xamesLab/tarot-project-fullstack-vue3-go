package decks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var defaultUploadStorage = NewUploadStorage(resolveUploadsRoot())

type UploadStorage struct {
	root string
}

func NewUploadStorage(root string) *UploadStorage {
	return &UploadStorage{root: filepath.Clean(root)}
}

func (s *UploadStorage) Save(relativePath string, payload []byte) (string, error) {
	cleanRelative := filepath.Clean(relativePath)
	if strings.Contains(cleanRelative, "..") {
		return "", ErrInvalidInput
	}

	fsPath := filepath.Join(s.root, cleanRelative)
	if err := os.MkdirAll(filepath.Dir(fsPath), 0o755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}

	dst, err := os.Create(fsPath)
	if err != nil {
		return "", fmt.Errorf("create upload file: %w", err)
	}
	defer dst.Close()

	if _, err := dst.Write(payload); err != nil {
		return "", fmt.Errorf("write upload file: %w", err)
	}

	return "/uploads/" + filepath.ToSlash(cleanRelative), nil
}

func (s *UploadStorage) RootDir() string {
	return s.root
}

func UploadsRootDir() string {
	return defaultUploadStorage.RootDir()
}

func resolveUploadsRoot() string {
	if raw := strings.TrimSpace(os.Getenv("UPLOADS_DIR")); raw != "" {
		return filepath.Clean(raw)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return filepath.Join("storage", "uploads")
	}

	if isBackendModuleRoot(cwd) {
		return filepath.Join(cwd, "storage", "uploads")
	}

	backendDir := filepath.Join(cwd, "backend")
	if isBackendModuleRoot(backendDir) {
		return filepath.Join(backendDir, "storage", "uploads")
	}

	return filepath.Join(cwd, "storage", "uploads")
}

func isBackendModuleRoot(dir string) bool {
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "module tarotdaily/backend")
}
