package generate

import (
	templateHelper "binh-swagger/cmd/swagger/commands/generate/internal"
	"binh-swagger/cmd/swagger/commands/internal/pkg"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrInvalidPathName = errors.New("invalid path name")

// assuming server command is created by the orchestrator and contains the necessary data for server generation.
func Server(cmd *ServerCommand, cfg Config) error {
	_, err := GetProjectStructure()
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	// load the server template
	tmpl, err := GetBaseTemplate(&buf, cfg.FileHelper(), templateHelper.ServerTemplateKey)
	if err != nil {
		return err
	}
	pathNames := santisePathNames(cmd.PathNames)
	if err = validatePathNames(pathNames); err != nil {
		return err
	}

	tmplModel := serverTemplateModel{
		Paths:            pathNames,
		RoutesImportPath: filepath.Join(projectImportPath, projectStructure["routes"]),
	}
	importTemplateModel := importsTemplateModel{
		RoutesImportPath: tmplModel.RoutesImportPath,
	}

	err = ExecuteTemplate(&importTemplateModel, tmpl, &buf, templateHelper.Templates.ImportsDefine)
	if err != nil {
		return err
	}

	err = ExecuteTemplate(&tmplModel, tmpl, &buf, templateHelper.Templates.ServerMainDefine)
	if err != nil {
		return err
	}
	// todo move server.go to a constant
	outputFile := cfg.FileHelper().GetAbsoluteSanitiseFilePath(projectStructure["server"], "server.go")
	return os.WriteFile(outputFile, buf.Bytes(), pkg.FilePermOwnerReadWrite)
}

func validatePathNames(pathNames []string) error {
	for _, path := range pathNames {
		if path == "" {
			return ErrInvalidPathName
		}
		if !strings.HasPrefix(path, "/") {
			return ErrInvalidPathName
		}
	}
	return nil
}

func santisePathNames(pathNames []string) []string {
	var santised []string
	for _, path := range pathNames {
		if !strings.HasPrefix(path, "/") {
			santised = append(santised, "/"+path)
		} else {
			santised = append(santised, path)
		}
	}
	return santised
}
