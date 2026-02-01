package generate

import (
	fileHelper "binh-swagger/cmd/swagger/commands/adaptor"
	templateHelper "binh-swagger/cmd/swagger/commands/generate/internal"
	"binh-swagger/cmd/swagger/commands/internal/pkg"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"bytes"
	"html/template"
	"os"
)

const (
	handlerFileSuffix = "_handler.go"
	routesFileSuffix  = "_routes.go"
)

type operation struct {
	method string
	op     *spec.Operation
}

type importsStruct struct {
	ModelImportPath   string
	HandlerImportPath string
}

func Path(cmd *PathCommand, config Config) error {
	fHelper := config.FileHelper()
	ops := []operation{
		{"GET", cmd.Get},
		{"POST", cmd.Post},
		{"PUT", cmd.Put},
		{"DELETE", cmd.Delete},
		{"PATCH", cmd.Patch},
	}
	importsData := importsStruct{
		ModelImportPath:   cmd.ModelImportPath,
		HandlerImportPath: cmd.HandlerImportPath,
	}
	if err := createHandlers(fHelper, importsData, ops, cmd.Name+handlerFileSuffix); err != nil {
		return err
	}
	// todo could be optimized by looking through ops one and generating both handlers and routes in one pass.
	if err := createRoutes(fHelper, importsData, ops, cmd.Name+routesFileSuffix); err != nil {
		return err
	}

	return nil
}

// todo test this.
func createRoutes(fHelper fileHelper.FileHelper, importsData importsStruct, ops []operation, routeFileName string) error {
	var buf bytes.Buffer
	tmpl, err := getRoutesTemplate(&buf, fHelper)
	if err != nil {
		return err
	}

	if err = executeImportsTemplate(importsData, tmpl, &buf); err != nil {
		return err
	}

	if err = executeRoutesTemplate(importsData, tmpl, &buf); err != nil {
		return err
	}

	outputFile := fHelper.GetAbsoluteSanitiseFilePath(projectStructure["routes"], routeFileName)
	return os.WriteFile(outputFile, buf.Bytes(), pkg.FilePermOwnerReadWrite)
}

func createHandlers(fHelper fileHelper.FileHelper, importsData importsStruct, ops []operation, handlerFilename string) error {
	if _, err := GetProjectStructure(); err != nil {
		return err
	}

	var buf bytes.Buffer
	tmpl, err := getHandlerTemplate(&buf, fHelper)
	if err != nil {
		return err
	}

	if err = executeImportsTemplate(importsData, tmpl, &buf); err != nil {
		return err
	}
	for _, o := range ops {
		err := executeHandlerTemplate(tmpl, &buf, o)
		if err != nil {
			return err
		}
	}
	outputFile := fHelper.GetAbsoluteSanitiseFilePath(projectStructure["handlers"], handlerFilename)
	return os.WriteFile(outputFile, buf.Bytes(), pkg.FilePermOwnerReadWrite)
}

// getHandlerTemplate todo could be imporoved by using generics and combining getRoutesTemplate.
func getHandlerTemplate(buf *bytes.Buffer, fHelper fileHelper.FileHelper) (*template.Template, error) {
	tmpl, err := templateHelper.LoadModelTemplate(fHelper, templateHelper.HandlerTemplateKey)
	if err != nil {
		return nil, err
	}

	if err = tmpl.ExecuteTemplate(buf, "handler", nil); err != nil {
		return nil, err
	}

	return tmpl, nil
}

func getRoutesTemplate(buf *bytes.Buffer, fHelper fileHelper.FileHelper) (*template.Template, error) {
	tmpl, err := templateHelper.LoadModelTemplate(fHelper, templateHelper.RouteTemplateKey)
	if err != nil {
		return nil, err
	}

	if err = tmpl.ExecuteTemplate(buf, "route", nil); err != nil {
		return nil, err
	}

	return tmpl, nil
}

// executeImportsTemplate todo could be imporoved by using generics and combining the other execute functions.
func executeRoutesTemplate(importsData importsStruct, tmpl *template.Template, buf *bytes.Buffer) error {
	if err := tmpl.ExecuteTemplate(buf, "registerRoutes", importsData); err != nil {
		return err
	}
	return nil
}

// this is generic enough to be used for both routes and handlers. becuase of the {{ define imports }} in both templates.
func executeImportsTemplate(importsData importsStruct, tmpl *template.Template, buf *bytes.Buffer) error {
	if err := tmpl.ExecuteTemplate(buf, "imports", importsData); err != nil {
		return err
	}
	return nil
}

func executeHandlerTemplate(tmpl *template.Template, buf *bytes.Buffer, cfg operation) error {
	if cfg.op == nil {
		return nil
	}

	if err := tmpl.ExecuteTemplate(buf, cfg.method, cfg.op); err != nil {
		return err
	}
	return nil
}
