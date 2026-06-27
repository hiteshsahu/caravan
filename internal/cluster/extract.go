package cluster

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// The GPU Slurm cluster scaffold travels inside the binary.
//
//go:embed assets/*
var assets embed.FS

// scaffoldDir is where the embedded assets get written before compose runs.
// Override with CARAVAN_DIR.
func scaffoldDir() (string, error) {
	if d := os.Getenv("CARAVAN_DIR"); d != "" {
		return d, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".caravan", "cluster"), nil
}

// Extract writes the embedded scaffold into dir.
func Extract(dir string) error {
	return fs.WalkDir(assets, "assets", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == "assets" {
			return nil
		}
		rel := strings.TrimPrefix(p, "assets/")
		target := filepath.Join(dir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		b, err := assets.ReadFile(p)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		mode := os.FileMode(0o644)
		if strings.HasSuffix(rel, ".sh") {
			mode = 0o755
		}
		return os.WriteFile(target, b, mode)
	})
}
