package helpers

import (
	"binh-swagger/cmd/swagger/commands/internal/pkg"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrNoGoModFile    = errors.New("no go.mod file found in the current directory")
	ErrEmptyGoModFile = errors.New("go.mod file is empty")
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

func EnsureOutputDirectoryExists(dirPath string, output io.Writer, input io.Reader) (bool, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if !promptCreateDirectory(dirPath, output, input) {
			return true, errors.New("directory creation cancelled by user")
		}
		if _, err := CreateDirectory(dirPath); err != nil {
			return true, err
		}

		if err := os.MkdirAll(dirPath, pkg.FileModeExecutable); err != nil {
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

func CreateDirectory(rootPath string) (string, error) {
	dirPath := filepath.Clean(rootPath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, pkg.FileModeExecutable); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return dirPath, nil
}

func ReadAllChildDirectoriesRecursive(path string) ([]string, error) {
	var allDirs []string

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			childPath := path + "/" + entry.Name()
			allDirs = append(allDirs, childPath)

			// Recursively read child directories
			subDirs, err := ReadAllChildDirectoriesRecursive(childPath)
			if err != nil {
				return nil, err
			}
			allDirs = append(allDirs, subDirs...)
		}
	}

	return allDirs, nil
}

func HasGoModFile() error {
	_, err := os.Stat("go.mod")
	if os.IsNotExist(err) {
		return fmt.Errorf("%w", ErrNoGoModFile)
	}
	return nil
}

func GetGoModImportPath2() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	return "", fmt.Errorf("module path not found in go.mod")
}

func GetGoModImportPath() (string, error) {
	data, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	reader := bufio.NewReader(data)
	for {
		line, err := reader.ReadString('\n')
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
		if err != nil && err == io.EOF {
			return "", fmt.Errorf("%w: End of file reached", ErrEmptyGoModFile)
		}
	}
}
