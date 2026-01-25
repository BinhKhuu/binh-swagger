package generate

import (
	fileHelper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"os"
)

func handleOperation(method string, cfg *spec.PathSpec, fHelper fileHelper.FileHelper) error {
	if cfg == nil {
		return nil
	}
	projectStructure, err := GetProjectStructure()
	if err != nil {
		return err
	}

	// generate handler
	err = createHandlerFile(projectStructure["handlers"], cfg, fHelper)
	if err != nil {
		return err
	}
	// generate route

	return nil
}

func createHandlerFile(handlersDir string, cfg *spec.PathSpec, fHelper fileHelper.FileHelper) error {
	outputFile := fHelper.GetAbsoluteSanitiseFilePath(handlersDir, cfg.Name+"_handler.go")
	// todo make tmpl file
	// parse spec into tmp file

	// save
	return os.WriteFile(outputFile, []byte("// Handler code here"), filePermOwnerReadWrite)
}

func Path(cfg *spec.PathSpec, fHelper fileHelper.FileHelper) error {
	ops := []struct {
		method string
		op     *spec.Operation
	}{
		{"GET", cfg.Get},
		{"POST", cfg.Post},
		{"PUT", cfg.Put},
		{"DELETE", cfg.Delete},
		{"PATCH", cfg.Patch},
	}

	for _, o := range ops {
		err := handleOperation(o.method, cfg, fHelper)
		if err != nil {
			return err
		}
	}

	return nil
}
