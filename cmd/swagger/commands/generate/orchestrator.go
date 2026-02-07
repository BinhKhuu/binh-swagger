package generate

import (
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"path/filepath"
)

type Orchestrator interface {
	FromAPIConfig(cfg *spec.APIConfig, command Config) error
}

type DefaultOrchestrator struct{}

func (o *DefaultOrchestrator) FromAPIConfig(cfg *spec.APIConfig, command Config) error {
	return fromAPIConfig(cfg, command)
}

func fromAPIConfig(cfg *spec.APIConfig, command Config) error {
	// Check prerequisites
	fileHelper := command.FileHelper()
	if err := fileHelper.HasGoModFile(); err != nil {
		return err
	}

	// Generate Project Structure
	if _, err := Project(cfg, fileHelper, ""); err != nil {
		return err
	}

	// Generate models
	if _, err := GetProjectStructure(); err != nil {
		return err
	}

	for _, model := range cfg.Models {
		modelCmd, err := SpecToModelCommand(model)
		if err != nil {
			return err
		}
		// Pass fileHelper wrapper that implements Config interface
		if err := Model(modelCmd, command); err != nil {
			return err
		}
	}

	serverCmd := &ServerCommand{
		RoutesImportPath: filepath.Join(projectImportPath + "/" + projectStructure["routes"]),
	}
	// Generate paths
	for path, pathSpec := range cfg.Paths {
		pathCmd, err := SpecToPathCommand(pathSpec, path)
		if err != nil {
			return err
		}
		if err := Path(pathCmd, command); err != nil {
			return err
		}
		serverCmd.PathNames = append(serverCmd.PathNames, path)
	}

	// Generate server
	if err := Server(serverCmd, command); err != nil {
		return err
	}
	return nil
}
