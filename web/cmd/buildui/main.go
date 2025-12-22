package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	guiDir, err := locateGUIDir()
	if err != nil {
		log.Fatal(err)
	}

	if err := run(guiDir, "bun", "install"); err != nil {
		log.Fatalf("bun install failed: %v", err)
	}

	if err := run(guiDir, "bun", "run", "generate"); err != nil {
		log.Fatalf("houdini generate failed: %v", err)
	}

	if err := run(guiDir, "bun", "run", "build"); err != nil {
		log.Fatalf("bun run build failed: %v", err)
	}
}

func run(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %v: %w", name, args, err)
	}
	return nil
}

func locateGUIDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	original := wd
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(wd, "gui")
		if ok, _ := exists(filepath.Join(candidate, "package.json")); ok {
			return candidate, nil
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	return "", errors.New("gui directory not found from " + original)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
