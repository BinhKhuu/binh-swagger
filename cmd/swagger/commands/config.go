package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ConfigCommand struct {
	Out  io.Writer
	Args ConfigArgs `positional-args:"true"`
}

type ConfigArgs struct {
	Filepath string `positional-arg-name:"filepath"`
	Filename string `positional-arg-name:"filename"`
}

func (c *ConfigCommand) Execute(_ []string) error {
	fileInfo, err := getFileInfo(c)
	if err != nil {
		fmt.Fprintf(c.Out, "error file does not exist: %v", err)
		return err
	}

	fmt.Fprintf(c.Out, "fileinfo %v", fileInfo)

	return nil
}

func getFileInfo(c *ConfigCommand) (os.FileInfo, error) {
	sanitisedFilePath, err := getSantiseFilePath(c.Args)
	if err != nil {
		fmt.Fprintf(c.Out, "Error sanitising file path: %v", err)
		return nil, err
	}

	err = checkSymlinks(sanitisedFilePath)
	if err != nil {
		fmt.Fprintf(c.Out, "Error : %v", err)
	}

	fileInfo, err := os.Stat(sanitisedFilePath)
	if err != nil {
		message := "unable to access config file: %v"
		if os.IsNotExist(err) {
			message = "config file not found %v"
		}
		fmt.Fprintf(c.Out, message, err)
		return nil, err
	}

	err = validateFileInfo(fileInfo)
	if err != nil {
		fmt.Fprintf(c.Out, "file info validation error: %v", err)
		return nil, err
	}

	fmt.Fprintf(c.Out, "found file: %s", fileInfo.Name())
	return fileInfo, nil
}

func validateFileInfo(fileInfo os.FileInfo) error {
	if fileInfo.IsDir() {
		return fmt.Errorf("expected a file but got directory: %s", fileInfo.Name())
	}
	return nil
}

func getSantiseFilePath(args ConfigArgs) (string, error) {
	cleanPath := filepath.Clean(args.Filepath)
	cleanFile := filepath.Clean(args.Filename)
	cleanFilePath := filepath.Join(cleanPath, cleanFile)
	if !filepath.IsAbs(cleanFilePath) {
		cleanFilePath, _ = filepath.Abs(cleanFilePath)
	}
	return cleanFilePath, nil
}

func checkSymlinks(absFilePath string) error {
	info, err := os.Lstat(absFilePath)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("symlinks not allowed")
	}
	return nil
}
