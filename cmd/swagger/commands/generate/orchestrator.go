package generate

import (
	"binh-swagger/cmd/swagger/commands/internal/spec"
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

	if _, err := Project(fileHelper, ""); err != nil {
		return err
	}

	if _, err := GetProjectStructure(); err != nil {
		return err
	}

	// Validate the API config by converting the spec to commands.
	// Model Command, Path Command
	// Server Command is generated after all paths are generated to get the list of paths for the server template.
	// When the Commands have been generate THEN generate the project
	// because project path is required for validateAPI spec its called after initializing the project structure and before generating the models and paths.
	if err := ValidateAPISpec(cfg); err != nil {
		return err
	}

	modelCmds := make(map[string]*ModelCommand)
	for key, model := range cfg.Models {
		modelCmd, err := SpecToModelCommand(model)
		if err != nil {
			return err
		}
		if err := Model(modelCmd, command); err != nil {
			return err
		}
		modelCmds[key] = modelCmd
	}

	serverCmd := &ServerCommand{}
	for path, pathSpec := range cfg.Paths {
		pathCmd, err := SpecToPathCommand(pathSpec, path, modelCmds)
		if err != nil {
			return err
		}
		if err := Path(pathCmd, command); err != nil {
			return err
		}
		serverCmd.PathNames = append(serverCmd.PathNames, path)
	}

	if err := Server(serverCmd, command); err != nil {
		return err
	}
	return nil
}

func ValidateAPISpec(cfg *spec.APIConfig) error {
	modelCommands := make(map[string]*ModelCommand)
	pathCommands := make(map[string]*PathCommand)

	for key, model := range cfg.Models {
		modelCmd, err := SpecToModelCommand(model)
		if err != nil {
			return err
		}
		modelCommands[key] = modelCmd
	}

	for key, path := range cfg.Paths {
		pathCmd, err := SpecToPathCommand(path, key, modelCommands)
		if err != nil {
			return err
		}
		pathCommands[key] = pathCmd
	}
	return nil
}
