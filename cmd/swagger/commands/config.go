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
	f, err := getFile(c)
	if err != nil {
		fmt.Fprintf(c.Out, "there was an error retrieving the file: %v", err)
		return err
	}
	defer f.Close()

	return nil
}

// probably should rename to validate file since the file info object is not going to be used passed this point
func getFile(c *ConfigCommand) (*os.File, error) {
	sanitisedFilePath, err := getSanitiseFilePath(c.Args)
	if err != nil {
		return nil, err
	}

	err = checkSymlinks(sanitisedFilePath)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(sanitisedFilePath)
	if err != nil {
		message := "unable to access config file: %v"
		if os.IsNotExist(err) {
			message = "config file not found %v"
		}
		return nil, fmt.Errorf(message, err)
	}

	err = validateFileInfo(fileInfo)
	if err != nil {
		fmt.Fprintf(c.Out, "file info validation error: %v", err)
		return nil, err
	}

	fmt.Fprintf(c.Out, "found file: %s", fileInfo.Name())

	f, err := os.Open(sanitisedFilePath)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func validateFileInfo(fileInfo os.FileInfo) error {
	if fileInfo.IsDir() {
		return fmt.Errorf("expected a file but got directory: %s", fileInfo.Name())
	}
	return nil
}

func getSanitiseFilePath(args ConfigArgs) (string, error) {
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
