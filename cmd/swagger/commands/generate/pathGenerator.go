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

// todo move this to the helper and share between the generators
type templateModel interface {
	*importsTemplateModel | *spec.Operation | *routesTemplateModel | *serverTemplateModel
}

type routesTemplateModel struct {
	Routes []routeModel
}

// todo refine the fields because this matches with ServerCommand
type serverTemplateModel struct {
	Paths        []string
	RoutesImportPath string
}

type routeModel struct {
	PathName    string
	Method      string
	OperationID string
}

type importsTemplateModel struct {
	ModelImportPath   string
	HandlerImportPath string
	RoutesImportPath  string
}

func Path(cmd *PathCommand, config Config) error {
	if _, err := GetProjectStructure(); err != nil {
		return err
	}

	// should ops contain the PathName e.g. /users this will help with route generation
	ops := []operation{
		{"GET", cmd.Get},
		{"POST", cmd.Post},
		{"PUT", cmd.Put},
		{"DELETE", cmd.Delete},
		{"PATCH", cmd.Patch},
	}
	importsData := importsTemplateModel{
		ModelImportPath:   cmd.ModelImportPath,
		HandlerImportPath: cmd.HandlerImportPath,
	}
	if err := createHandlers(cmd, config, importsData, ops); err != nil {
		return err
	}
	// todo could be optimized by looking through ops one and generating both handlers and routes in one pass.
	if err := createRoutes(cmd, config, importsData, ops); err != nil {
		return err
	}

	return nil
}

func createRoutes(cmd *PathCommand, config Config, importsData importsTemplateModel, ops []operation) error {
	routeFileName := cmd.Name + routesFileSuffix
	fHelper := config.FileHelper()
	var buf bytes.Buffer
	tmpl, err := GetBaseTemplate(&buf, fHelper, templateHelper.RouteTemplateKey)
	if err != nil {
		return err
	}

	if err = ExecuteTemplate(&importsData, tmpl, &buf, templateHelper.Templates.ImportsDefine); err != nil {
		return err
	}

	// todo check if path starts with /
	rd := createRoutesData("/"+cmd.Name, ops)
	rm := routesTemplateModel{
		Routes: rd,
	}
	if err = ExecuteTemplate(&rm, tmpl, &buf, templateHelper.Templates.RouteDefine); err != nil {
		return err
	}

	outputFile := fHelper.GetAbsoluteSanitiseFilePath(projectStructure["routes"], routeFileName)
	return os.WriteFile(outputFile, buf.Bytes(), pkg.FilePermOwnerReadWrite)
}

func createHandlers(cmd *PathCommand, config Config, importsData importsTemplateModel, ops []operation) error {
	fHelper := config.FileHelper()
	handlerFilename := cmd.Name + handlerFileSuffix

	var buf bytes.Buffer
	tmpl, err := GetBaseTemplate(&buf, fHelper, templateHelper.HandlerTemplateKey)
	if err != nil {
		return err
	}

	if err = ExecuteTemplate(&importsData, tmpl, &buf, templateHelper.Templates.ImportsDefine); err != nil {
		return err
	}
	for _, o := range ops {
		if o.op == nil {
			continue
		}
		err := ExecuteTemplate(o.op, tmpl, &buf, o.method)
		if err != nil {
			return err
		}
	}
	outputFile := fHelper.GetAbsoluteSanitiseFilePath(projectStructure["handlers"], handlerFilename)
	return os.WriteFile(outputFile, buf.Bytes(), pkg.FilePermOwnerReadWrite)
}

func createRoutesData(pathName string, ops []operation) []routeModel {
	routes := make([]routeModel, 0, len(ops))
	for _, o := range ops {
		if o.op == nil {
			continue
		}
		route := routeModel{
			Method:      o.method,
			PathName:    pathName,
			OperationID: o.op.OperationID,
		}

		routes = append(routes, route)
	}
	return routes
}

// todo move this to a helper for the other generators to use. and change templateKey to the type which holds the tempalte keys
func GetBaseTemplate(buf *bytes.Buffer, fHelper fileHelper.FileHelper, templateKey string) (*template.Template, error) {
	tmpl, err := templateHelper.LoadModelTemplate(fHelper, templateKey)
	if err != nil {
		return nil, err
	}

	if err = tmpl.ExecuteTemplate(buf, templateKey, nil); err != nil {
		return nil, err
	}

	return tmpl, nil
}

// todo move this to a helper for the other generators to use. and change templateKey to the type which holds the tempalte keys
func ExecuteTemplate[T templateModel](data T, tmpl *template.Template, buf *bytes.Buffer, templateName string) error {
	if err := tmpl.ExecuteTemplate(buf, templateName, data); err != nil {
		return err
	}
	return nil
}
