package helper

import (
	"errors"
	"os"
	"path/filepath"
)

func FolderExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func LocateMygitPath() (string, error) {
	path := ".mygit"
	dir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	for {
		candidate := filepath.Join(dir, path)

		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}

		parent := filepath.Dir(dir)

		if parent == dir {
			return "", errors.New("fatal: not a git repository (or any of the parent directories): .git")
		}

		dir = parent
	}

}
