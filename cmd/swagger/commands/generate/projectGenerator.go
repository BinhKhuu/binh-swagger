package generate

import (
	fileHelper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"os"
)

func Project(cfg *spec.APIConfig, fHelper fileHelper.FileHelper) error {
	rootDir, err := fHelper.CreateDirectory(cfg.ProjectRoot)
	if err != nil {
		return err
	}

	var createdDirs []string

	directories := getProjectStructure(rootDir)

	for _, dir := range directories {
		if _, err := fHelper.CreateDirectory(dir); err != nil {
			// Cleanup created directories
			for _, createdDir := range createdDirs {
				_ = os.RemoveAll(createdDir)
			}
			return err
		}
		createdDirs = append(createdDirs, dir)
	}

	return nil
}

func getProjectStructure(rootDir string) []string {
	return []string{
		rootDir + "/cmd/server",
		rootDir + "/internal/handlers",
		rootDir + "/internal/routes",
		rootDir + "/services",
		rootDir + "/repository",
		rootDir + "/models",
		rootDir + "/middleware",
		rootDir + "/config",
	}
}
