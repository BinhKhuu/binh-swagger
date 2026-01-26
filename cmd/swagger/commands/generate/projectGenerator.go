package generate

import (
	fileHelper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"errors"
	"os"
)

var projectStructure map[string]string

func Project(cfg *spec.APIConfig, fHelper fileHelper.FileHelper) (string, error) {
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

func GetProjectStructure() (map[string]string, error) {
	if projectStructure == nil {
		return nil, errors.New("project structure not initialized")
	}
	return projectStructure, nil
}

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
