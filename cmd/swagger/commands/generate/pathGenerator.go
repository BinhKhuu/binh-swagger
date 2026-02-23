package generate

import (
	templateHelper "binh-swagger/cmd/swagger/commands/generate/internal/templateHelper"
	"binh-swagger/cmd/swagger/commands/internal/pkg"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"bytes"
	"errors"
	"net/http"
	"os"
	"unicode"
)

const (
	handlerFileSuffix = "_handler.go"
	routesFileSuffix  = "_routes.go"
)

var ErrOperationIDNotDefined = errors.New("operationId is not defined for this operation")

// todo replace this
type operation struct {
	method string
	op     *spec.Operation
}

func Path(cmd *PathCommand, config Config) error {
	if _, err := GetProjectStructure(); err != nil {
		return err
	}

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
		if op, err := toOperationsModel(*cmd.Get, http.MethodGet); err == nil {
			ops = append(ops, op)
		}
	}
	if cmd.Post != nil {
		if op, err := toOperationsModel(*cmd.Post, http.MethodPost); err == nil {
			ops = append(ops, op)
		}
	}
	if cmd.Put != nil {
		if op, err := toOperationsModel(*cmd.Put, http.MethodPut); err == nil {
			ops = append(ops, op)
		}
	}
	if cmd.Delete != nil {
		if op, err := toOperationsModel(*cmd.Delete, http.MethodDelete); err == nil {
			ops = append(ops, op)
		}
	}
	if cmd.Patch != nil {
		if op, err := toOperationsModel(*cmd.Patch, http.MethodPatch); err == nil {
			ops = append(ops, op)
		}
	}

	return ops
}

func toOperationsModel(op OperationCommand, methodType string) (templateHelper.OperationModel, error) {
	responseModel := make(map[int]templateHelper.ResponseModel, len(op.Responses))
	for code, response := range op.Responses {
		responseModel[code] = templateHelper.ResponseModel{
			Description:       response.Description,
			Type:              response.Type,
			Ref:               response.Ref,
			SuccessReturnCode: response.SuccessReturnCode,
		}
	}
	r := []rune(op.OperationID)
	if len(r) == 0 {
		return templateHelper.OperationModel{}, ErrOperationIDNotDefined
	}

	r[0] = unicode.ToUpper(r[0])
	return templateHelper.OperationModel{
		MethodType:  methodType,
		OperationID: string(r),
		Summary:     op.Summary,
		Produces:    op.Produces,
		Responses:   responseModel,
		ReturnType:  op.ReturnType,
	}, nil
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
