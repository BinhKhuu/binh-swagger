package generate

import (
	fileHelper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"errors"
	"os"
)

var (
	projectStructure          map[string]string
	projectImportPath         string
	ProjectNotInitilizedError = errors.New("project structure not initialized")
)

func Project(cfg *spec.APIConfig, fHelper fileHelper.FileHelper) (string, error) {
	importPath, err := fHelper.GetGoModImportPath()
	if err != nil {
		return "", err
	}

	projectImportPath = importPath
	rootDir, err := fHelper.CreateDirectory(cfg.ProjectRoot)
	if err != nil {
		return "", err
	}

	err = SetProjectStructure(rootDir)
	if err != nil {
		return "", err
	}

	createdDirs := []string{}

	directories, err := GetProjectStructure()
	if err != nil {
		return "", err
	}

	for _, dir := range directories {
		if _, err := fHelper.CreateDirectory(dir); err != nil {
			// Cleanup created directories
			for _, createdDir := range createdDirs {
				_ = os.RemoveAll(createdDir)
			}
			return "", err
		}
		createdDirs = append(createdDirs, dir)
	}

	return rootDir, nil
}

// GetProjectStructure todo think about the coupling of this, model, handler and route generateor needs this to be set before it can execute and the current setter only allows one set action.
// GetProjectStructure add error variable export and change new error to that variable.
func GetProjectStructure() (map[string]string, error) {
	if projectStructure == nil {
		return nil, ProjectNotInitilizedError
	}
	return projectStructure, nil
}

// SetProjectStructure todo think if we need to make this setable more than once to avoid coupling in model, handler and route generators.
func SetProjectStructure(rootDir string) error {
	if projectStructure != nil {
		return errors.New("project structure already initialized")
	}
	projectStructure = map[string]string{
		"server":     rootDir + "/cmd/server",
		"handlers":   rootDir + "/internal/handlers",
		"routes":     rootDir + "/internal/routes",
		"services":   rootDir + "/services",
		"repository": rootDir + "/repository",
		"models":     rootDir + "/models",
		"middleware": rootDir + "/middleware",
		"config":     rootDir + "/config",
	}
	return nil
}

func GetImportPath() (string, error) {
	if projectImportPath == "" {
		return "", errors.New("project import path not initialized")
	}
	return projectImportPath, nil
}
