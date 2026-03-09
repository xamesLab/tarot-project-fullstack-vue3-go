package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvMap(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "# comment\nMONGO_URI=mongodb://localhost:27017\n INVALID\nJWT_SECRET = secret-value \nEMPTY_KEY=\n"

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := LoadEnvMap(path)
	if err != nil {
		t.Fatalf("LoadEnvMap() error = %v", err)
	}

	if got["MONGO_URI"] != "mongodb://localhost:27017" {
		t.Fatalf("MONGO_URI = %q", got["MONGO_URI"])
	}
	if got["JWT_SECRET"] != "secret-value" {
		t.Fatalf("JWT_SECRET = %q", got["JWT_SECRET"])
	}
	if _, ok := got["INVALID"]; ok {
		t.Fatalf("unexpected INVALID key in env map")
	}
}

func TestLoadEnvMapMissingFile(t *testing.T) {
	got, err := LoadEnvMap(filepath.Join(t.TempDir(), "missing.env"))
	if err != nil {
		t.Fatalf("LoadEnvMap(missing) error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("LoadEnvMap(missing) len = %d, want 0", len(got))
	}
}
