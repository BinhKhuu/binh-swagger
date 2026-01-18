package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	filePermOwnerReadWrite = 0o600
	fileModeExecutable     = 0o755
)

func ValidateFileInfo(fileInfo os.FileInfo) error {
	if fileInfo.IsDir() {
		return fmt.Errorf("expected a file but got directory: %s", fileInfo.Name())
	}
	return nil
}

// GetAbsoluteSanitiseFilePath constructs a sanitized absolute file path from the given location and filename.
func GetAbsoluteSanitiseFilePath(filelocation string, filename string) string {
	cleanPath := filepath.Clean(filelocation)
	cleanFile := filepath.Clean(filename)
	cleanFilePath := filepath.Join(cleanPath, cleanFile)
	if !filepath.IsAbs(cleanFilePath) {
		cleanFilePath, _ = filepath.Abs(cleanFilePath)
	}
	return cleanFilePath
}

func CheckSymlinks(absFilePath string) error {
	info, err := os.Lstat(absFilePath)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return errors.New("symlinks not allowed")
	}
	return nil
}

func EnsureOutputDirectoryExists(outputPath string, output io.Writer, input io.Reader) (bool, error) {
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if !promptCreateDirectory(outputPath, output, input) {
			return true, errors.New("directory creation cancelled by user")
		}
		if err := os.MkdirAll(outputPath, fileModeExecutable); err != nil {
			return true, fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return false, nil
}

func promptCreateDirectory(path string, output io.Writer, input io.Reader) bool {
	reader := bufio.NewReader(input)
	fmt.Fprintf(output, "\nDirectory %s does not exist. Create it? it will be made in the directory relative to where you ran the command [y/N]: ", path)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
