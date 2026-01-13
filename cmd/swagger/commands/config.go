package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type ConfigCommand struct {
	Out  io.Writer
	Args ConfigArgs `positional-args:"true"`
}

type ConfigArgs struct {
	Filepath string `positional-arg-name:"filepath"`
	Filename string `positional-arg-name:"filename"`
}

type Config struct {
	Name    string `yaml:"name"`
	Version int    `yaml:"version"`
	Enabled bool   `yaml:"enabled"`

	Metadata struct {
		Author  string   `yaml:"author"`
		Created string   `yaml:"created"`
		Tags    []string `yaml:"tags"`
	} `yaml:"metadata"`

	Config struct {
		Retries        int `yaml:"retries"`
		TimeoutSeconds int `yaml:"timeout_seconds"`
		Paths          struct {
			Input  string `yaml:"input"`
			Output string `yaml:"output"`
		} `yaml:"paths"`
	} `yaml:"config"`

	Items []struct {
		ID    int    `yaml:"id"`
		Value string `yaml:"value"`
	} `yaml:"items"`
}

func (c *ConfigCommand) Execute(_ []string) error {
	f, err := validateAndOpenFile(c)
	if err != nil {
		fmt.Fprintf(c.Out, "there was an error retrieving the file: %v", err)
		return err
	}
	defer f.Close()

	config, err := parseFile(f)
	if err != nil {
		fmt.Fprintf(c.Out, "there was an error prasing the file: %v", err)
	}
	fmt.Fprintf(c.Out, "%+v\n", config)
	return nil
}

func parseFile(file *os.File) (*Config, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func validateAndOpenFile(c *ConfigCommand) (*os.File, error) {
	sanitisedFilePath := getSanitiseFilePath(c.Args)

	err := checkSymlinks(sanitisedFilePath)
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

func getSanitiseFilePath(args ConfigArgs) string {
	cleanPath := filepath.Clean(args.Filepath)
	cleanFile := filepath.Clean(args.Filename)
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
