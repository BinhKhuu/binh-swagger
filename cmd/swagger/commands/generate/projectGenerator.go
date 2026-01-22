package generate

import (
	fileHelper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/internal/spec"
)

func Project(cfg *spec.APIConfig, fHelper fileHelper.FileHelper) error {
	rootDir, err := fHelper.CreateDirectory(cfg.ProjectRoot)
	if err != nil {
		return err
	}
	_ = rootDir // Use rootDir to generate project structure as needed

	// Generate Server folders
	return nil
}
