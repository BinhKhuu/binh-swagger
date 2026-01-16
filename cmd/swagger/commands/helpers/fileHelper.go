package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func validateFileInfo(fileInfo os.FileInfo) error {
	if fileInfo.IsDir() {
		return fmt.Errorf("expected a file but got directory: %s", fileInfo.Name())
	}
	return nil
}

// getSanitiseFilePath constructs a sanitized absolute file path
func getSanitiseFilePath(filelocation string, filename string) string {
	cleanPath := filepath.Clean(filelocation)
	cleanFile := filepath.Clean(filename)
	cleanFilePath := filepath.Join(cleanPath, cleanFile)
	if !filepath.IsAbs(cleanFilePath) {
		cleanFilePath, _ = filepath.Abs(cleanFilePath)
	}
	return cleanFilePath
}

func checkSymlinks(absFilePath string) error {
	info, err := os.Lstat(absFilePath)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return errors.New("symlinks not allowed")
	}
	return nil
}
