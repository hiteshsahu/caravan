package cluster

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// The GPU Slurm cluster scaffold travels inside the binary.
//
//go:embed assets/*
var assets embed.FS

const project = "caravan"

// scaffoldDir is where the embedded assets get written before `docker compose`
// runs. Override with CARAVAN_DIR.
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

func composeFile(dir string) string { return filepath.Join(dir, "docker-compose.yml") }

// Extract writes the embedded scaffold (Dockerfile, compose, slurm.conf, …)
// into dir.
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

// Up extracts the scaffold and brings the cluster online.
func Up() error {
	dir, err := scaffoldDir()
	if err != nil {
		return err
	}
	fmt.Printf("→ scaffolding GPU Slurm cluster in %s\n", dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if err := Extract(dir); err != nil {
		return fmt.Errorf("extract scaffold: %w", err)
	}
	fmt.Println("→ docker compose up (first run builds the image)…")
	return runIn(dir, "docker", "compose", "-p", project, "-f", composeFile(dir), "up", "-d", "--build")
}

// Down stops the cluster; volumes also wipes its state.
func Down(volumes bool) error {
	dir, err := scaffoldDir()
	if err != nil {
		return err
	}
	args := []string{"compose", "-p", project, "-f", composeFile(dir), "down"}
	if volumes {
		args = append(args, "--volumes")
	}
	return runIn(dir, "docker", args...)
}

// Status prints container state, then Slurm node state.
func Status() error {
	dir, err := scaffoldDir()
	if err != nil {
		return err
	}
	_ = runIn(dir, "docker", "compose", "-p", project, "-f", composeFile(dir), "ps")
	fmt.Println()
	return run("docker", "exec", "slurmctld", "sinfo")
}

func run(name string, args ...string) error { return runIn("", name, args...) }

func runIn(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
