package commands

import (
	"binh-swagger/cmd/swagger/commands/generate"
	"binh-swagger/cmd/swagger/commands/helpers"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type ConfigCommand struct {
	BaseCommand
	Helpers

	API  bool       `description:"validate the configuration against the api model" long:"api" optional:"true"`
	CMD  bool       `description:"validate against the command model"               long:"cmd"`
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
	err := validateConfigFlags(c)
	if err != nil {
		fmt.Fprintf(c.Out, "there was an error with the provided flags: %v\n", err)
		return err
	}

	f, err := validateAndOpenFile(c)
	if err != nil {
		fmt.Fprintf(c.Out, "there was an error retrieving the file: %v\n", err)
		return err
	}
	defer f.Close()

	if c.API {
		config, err := parseFile[spec.APIConfig](f)
		if err != nil {
			fmt.Fprintf(c.Out, "there was an error prasing the file: %v\n", err)
		}
		err = generateFromAPIConfig(config, c)
		if err != nil {
			fmt.Fprintf(c.Out, "there was an error generating from api config: %v\n", err)
		}

		fmt.Fprintf(c.Out, "api configuration file %s validated successfully and models generated\n", c.Args.Filename)
	}

	return nil
}

func generateFromAPIConfig(cfg *spec.APIConfig, command *ConfigCommand) error {
	var err error

	// todo: generate.Project
	generate.Project(cfg, command.File)
	// todo: generate.Server

	// todo: generate.Routes

	// todo: generate.Docs
	for _, model := range cfg.Models {
		err = generate.Model(model, command.File)
		if err != nil {
			return err
		}
	}
	return err
}

func validateConfigFlags(c *ConfigCommand) error {
	count := 0
	if c.API {
		count++
	}

	if c.CMD {
		count++
	}

	if count > 1 {
		return errors.New("only one of --api or --cmd can be specified")
	}

	if count == 0 {
		return errors.New("one of --api or --cmd must be specified")
	}

	return nil
}

func parseFile[T any](file *os.File) (*T, error) {
	var zero T

	if _, ok := any(zero).(Config); ok {
		return parseYAML[T](file)
	}

	if _, ok := any(zero).(spec.APIConfig); ok {
		return parseYAML[T](file)
	}

	return nil, fmt.Errorf("unsupported type %T", zero)
}

// Alternative implementation using type switch.
// func parseFile2[T any](file *os.File) (*T, error) {
// 	var (
// 		err    error
// 		v      T
// 		result any
// 	)

// 	switch any(v).(type) {
// 	case APIConfig:
// 		result, err = parseYAML[APIConfig](file)
// 	case Config:
// 		result, err = parseYAML[Config](file)
// 	default:
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	typed, ok := result.(*T)
// 	if !ok {
// 		return nil, fmt.Errorf("type mismatch: expected %T", *new(T))
// 	}

// 	return typed, nil
// }

func parseYAML[T any](file *os.File) (*T, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var v T
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func validateAndOpenFile(c *ConfigCommand) (*os.File, error) {
	cleanPath := filepath.Clean(c.Args.Filepath)
	cleanFile := filepath.Clean(c.Args.Filename)
	sanitisedFilePath := helpers.GetAbsoluteSanitiseFilePath(cleanPath, cleanFile)

	err := helpers.CheckSymlinks(sanitisedFilePath)
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

	err = helpers.ValidateFileInfo(fileInfo)
	if err != nil {
		fmt.Fprintf(c.Out, "file info validation error: %v\n", err)
		return nil, err
	}

	fmt.Fprintf(c.Out, "found file: %s\n", fileInfo.Name())

	f, err := os.Open(sanitisedFilePath)
	if err != nil {
		return nil, err
	}

	return f, nil
}
