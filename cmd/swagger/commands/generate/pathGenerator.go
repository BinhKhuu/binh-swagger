package generate

import (
	templateHelper "binh-swagger/cmd/swagger/commands/generate/internal/templateHelper"
	"binh-swagger/cmd/swagger/commands/internal/pkg"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"bytes"
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

func Path(cmd *PathCommand, config Config) error {
	if _, err := GetProjectStructure(); err != nil {
		return err
	}

	// should ops contain the PathName e.g. /users this will help with route generation
	ops := createOperationModels(cmd)
	importsData := templateHelper.ImportsTemplateModel{
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

func createOperationModels(cmd *PathCommand) []templateHelper.OperationModel {
	var ops []templateHelper.OperationModel

	if cmd.Get != nil {
		ops = append(ops, toOperationsModel(*cmd.Get, "GET"))
	}
	if cmd.Post != nil {
		ops = append(ops, toOperationsModel(*cmd.Post, "POST"))
	}
	if cmd.Put != nil {
		ops = append(ops, toOperationsModel(*cmd.Put, "PUT"))
	}
	if cmd.Delete != nil {
		ops = append(ops, toOperationsModel(*cmd.Delete, "DELETE"))
	}
	if cmd.Patch != nil {
		ops = append(ops, toOperationsModel(*cmd.Patch, "PATCH"))
	}

	return ops
}

func toOperationsModel(op spec.Operation, methodType string) templateHelper.OperationModel {
	responseModel := make(map[int]templateHelper.ResponseModel, len(op.Responses))
	for code, response := range op.Responses {
		responseModel[code] = templateHelper.ResponseModel{
			Description: response.Description,
			Type:        response.Schema.Type,
			Ref:         response.Schema.Ref,
		}
	}

	return templateHelper.OperationModel{
		MethodType:  methodType,
		OperationID: op.OperationID,
		Summary:     op.Summary,
		Produces:    op.Produces,
		Responses:   responseModel,
	}
}

func createRoutes(cmd *PathCommand, config Config, importsData templateHelper.ImportsTemplateModel, ops []templateHelper.OperationModel) error {
	routeFileName := cmd.Name + routesFileSuffix
	fHelper := config.FileHelper()
	var buf bytes.Buffer
	tmpl, err := templateHelper.GetBaseTemplate(&buf, fHelper, templateHelper.RouteTemplateKey)
	if err != nil {
		return err
	}

	if err = templateHelper.ExecuteTemplate(&importsData, tmpl, &buf, templateHelper.Templates.ImportsDefine); err != nil {
		return err
	}

	// todo check if path starts with /
	rd := createRoutesData("/"+cmd.Name, ops)
	rm := templateHelper.RoutesTemplateModel{
		Routes: rd,
	}
	if err = templateHelper.ExecuteTemplate(&rm, tmpl, &buf, templateHelper.Templates.RouteDefine); err != nil {
		return err
	}

	outputFile := fHelper.GetAbsoluteSanitiseFilePath(projectStructure["routes"], routeFileName)
	return os.WriteFile(outputFile, buf.Bytes(), pkg.FilePermOwnerReadWrite)
}

func createHandlers(cmd *PathCommand, config Config, importsData templateHelper.ImportsTemplateModel, ops []templateHelper.OperationModel) error {
	fHelper := config.FileHelper()
	handlerFilename := cmd.Name + handlerFileSuffix

	var buf bytes.Buffer
	tmpl, err := templateHelper.GetBaseTemplate(&buf, fHelper, templateHelper.HandlerTemplateKey)
	if err != nil {
		return err
	}

	if err = templateHelper.ExecuteTemplate(&importsData, tmpl, &buf, templateHelper.Templates.ImportsDefine); err != nil {
		return err
	}
	for _, o := range ops {
		if o.OperationID == "" {
			continue
		}
		err := templateHelper.ExecuteTemplate(&o, tmpl, &buf, o.MethodType)
		if err != nil {
			return err
		}
	}
	outputFile := fHelper.GetAbsoluteSanitiseFilePath(projectStructure["handlers"], handlerFilename)
	return os.WriteFile(outputFile, buf.Bytes(), pkg.FilePermOwnerReadWrite)
}

func createRoutesData(pathName string, ops []templateHelper.OperationModel) []templateHelper.RouteModel {
	routes := make([]templateHelper.RouteModel, 0, len(ops))
	for _, o := range ops {
		if o.OperationID == "" {
			continue
		}
		route := templateHelper.RouteModel{
			Method:      o.MethodType,
			PathName:    pathName,
			OperationID: o.OperationID,
		}

		routes = append(routes, route)
	}
	return routes
}
