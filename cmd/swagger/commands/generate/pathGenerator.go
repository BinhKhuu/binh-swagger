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

type templateModel interface {
	*importsTemplateModel | *spec.Operation | *routesTemplateModel
}

// todo rename this
type routesTemplateModel struct {
	Routes []routeModel
}

// todo rename this
type routeModel struct {
	PathName    string
	Method      string
	OperationId string
}

// todo rename this
type importsTemplateModel struct {
	ModelImportPath   string
	HandlerImportPath string
}

func Path(cmd *PathCommand, config Config) error {
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

// todo test this.
func createRoutes(cmd *PathCommand, config Config, importsData importsTemplateModel, ops []operation) error {
	routeFileName := cmd.Name + routesFileSuffix
	fHelper := config.FileHelper()
	var buf bytes.Buffer
	tmpl, err := getBaseTemplate(&buf, fHelper, templateHelper.RouteTemplateKey)
	if err != nil {
		return err
	}

	if err = executeTemplate(&importsData, tmpl, &buf, templateHelper.ImportsDefine); err != nil {
		return err
	}

	// todo check if path starts with /
	routesData := createRoutesData("/"+cmd.Name, ops)
	routesTemplateModel := routesTemplateModel{
		Routes: routesData,
	}
	if err = executeTemplate(&routesTemplateModel, tmpl, &buf, templateHelper.RouteDefine); err != nil {
		return err
	}

	outputFile := fHelper.GetAbsoluteSanitiseFilePath(projectStructure["routes"], routeFileName)
	return os.WriteFile(outputFile, buf.Bytes(), pkg.FilePermOwnerReadWrite)
}

func createHandlers(cmd *PathCommand, config Config, importsData importsTemplateModel, ops []operation) error {
	fHelper := config.FileHelper()
	handlerFilename := cmd.Name + handlerFileSuffix
	if _, err := GetProjectStructure(); err != nil {
		return err
	}

	var buf bytes.Buffer
	tmpl, err := getBaseTemplate(&buf, fHelper, templateHelper.HandlerTemplateKey)
	if err != nil {
		return err
	}

	if err = executeTemplate(&importsData, tmpl, &buf, templateHelper.ImportsDefine); err != nil {
		return err
	}
	for _, o := range ops {
		if o.op == nil {
			continue
		}
		err := executeTemplate(o.op, tmpl, &buf, o.method)
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
			OperationId: o.op.OperationID,
		}

		routes = append(routes, route)
	}
	return routes
}

// getBaseTemplate todo unit test this
func getBaseTemplate(buf *bytes.Buffer, fHelper fileHelper.FileHelper, templateKey string) (*template.Template, error) {
	tmpl, err := templateHelper.LoadModelTemplate(fHelper, templateKey)
	if err != nil {
		return nil, err
	}

	if err = tmpl.ExecuteTemplate(buf, templateKey, nil); err != nil {
		return nil, err
	}

	return tmpl, nil
}

func executeTemplate[T templateModel](data T, tmpl *template.Template, buf *bytes.Buffer, templateName string) error {
	if err := tmpl.ExecuteTemplate(buf, templateName, data); err != nil {
		return err
	}
	return nil
}
