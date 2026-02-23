package generate

import (
	filehelper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

const (
	schemaTypeArray  = "array"
	schemaTypeObject = "object"
)

var ErrSchemaNotFound = errors.New("schema not found for response")

type Config interface {
	FileHelper() filehelper.FileHelper
}

type ModelCommand struct {
	Name   string           `yaml:"name"`
	Fields []spec.FieldSpec `yaml:"fields"`
}

type PathCommand struct {
	ModelImportPath   string
	HandlerImportPath string
	Name              string
	Get               *OperationCommand
	Post              *OperationCommand
	Put               *OperationCommand
	Patch             *OperationCommand
	Delete            *OperationCommand
}

type OperationCommand struct {
	Summary     string
	OperationID string
	Produces    []string
	ReturnType  string // GIN does not use a return type on handlers Leaving this here for reference if I want to add support for a framework that does require a return type in the future. For now it will be used to generate the return code in the handler template.
	Responses   map[int]ResponseCommand
	ReturnCode  string
}

type ResponseCommand struct {
	Description       string
	Type              string
	Ref               string
	SuccessReturnCode string // this is the code that will be used in the handler template to generate the return statement for successful responses. It is generated based on the response schema and the model command associated with that schema. It is currently only generated for array and object types, but can be extended to support other types in the future if needed.
}

type ServerCommand struct {
	PathNames []string
}

func SpecToPathCommand(cfg spec.PathSpec, pathName string, mCmd map[string]*ModelCommand) (*PathCommand, error) {
	paths, err := GetProjectStructure()
	if err != nil {
		return nil, err
	}

	var getCmd, postCmd, putCmd, patchCmd, delCmd *OperationCommand

	if cfg.Get != nil {
		getCmd, err = SpecToOperation(*cfg.Get, mCmd)
		if err != nil {
			return nil, err
		}
	}

	if cfg.Post != nil {
		postCmd, err = SpecToOperation(*cfg.Post, mCmd)
		if err != nil {
			return nil, err
		}
	}

	if cfg.Put != nil {
		putCmd, err = SpecToOperation(*cfg.Put, mCmd)
		if err != nil {
			return nil, err
		}
	}

	if cfg.Patch != nil {
		patchCmd, err = SpecToOperation(*cfg.Patch, mCmd)
		if err != nil {
			return nil, err
		}
	}

	if cfg.Delete != nil {
		delCmd, err = SpecToOperation(*cfg.Delete, mCmd)
		if err != nil {
			return nil, err
		}
	}

	return &PathCommand{
		ModelImportPath:   filepath.Join(projectImportPath, paths["models"]),
		HandlerImportPath: filepath.Join(projectImportPath, paths["handlers"]),
		Name:              strings.Replace(pathName, "/", "", 1),
		Get:               getCmd,
		Post:              postCmd,
		Put:               putCmd,
		Patch:             patchCmd,
		Delete:            delCmd,
	}, nil
}

func SpecToOperation(cfg spec.Operation, mCmd map[string]*ModelCommand) (*OperationCommand, error) {
	opCmd := &OperationCommand{
		Summary:     cfg.Summary,
		OperationID: cfg.OperationID,
		Produces:    cfg.Produces,
	}
	resp, err := SpecToOperationResponses(cfg, mCmd)
	if err != nil {
		return nil, err
	}
	opCmd.Responses = resp
	opCmd.ReturnType = deriveReturnType(cfg)
	opCmd.ReturnCode = combineReturnCodes(resp)
	return opCmd, nil
}

func SpecToOperationResponses(cfg spec.Operation, mCdm map[string]*ModelCommand) (map[int]ResponseCommand, error) {
	responses := make(map[int]ResponseCommand)

	for responseKey, res := range cfg.Responses {
		resCmd := &ResponseCommand{
			Type:        "",
			Ref:         "",
			Description: "",
		}
		ref := getModelSchemaKey(res)
		resCmd.Description = res.Description
		if ref != "" {
			if mCdm[ref] == nil {
				return nil, fmt.Errorf("%w: schema %s not found", ErrSchemaNotFound, ref)
			}
			if res.Schema.Type == schemaTypeArray {
				resCmd.Ref = res.Schema.Items.Ref.Ref
			} else {
				resCmd.Ref = res.Schema.Ref.Ref
			}
			resCmd.Type = res.Schema.Type
			resCmd.Description = res.Description

			setReturnCode(res, mCdm[ref], resCmd, responseKey)
		}
		responses[responseKey] = *resCmd
	}
	return responses, nil
}

func setReturnCode(res spec.ResponseSpec, mCdm *ModelCommand, resCmd *ResponseCommand, responseKey int) {
	var code string
	switch res.Schema.Type {
	case schemaTypeArray:
		code = fmt.Sprintf("c.JSON(%d,[]models.%s{})", responseKey, mCdm.Name)
	case schemaTypeObject:
		code = fmt.Sprintf("c.JSON(%d, models.%s{})", responseKey, mCdm.Name)
	default:
		code = "// todo add comment on return type, currently only supports array and object types"
	}

	resCmd.SuccessReturnCode = code
}

func combineReturnCodes(responses map[int]ResponseCommand) string {
	var codes []string
	for _, res := range responses {
		if res.SuccessReturnCode != "" {
			codes = append(codes, res.SuccessReturnCode)
		}
	}
	return strings.Join(codes, "\n")
}

func getModelSchemaKey(res spec.ResponseSpec) string {
	if res.Schema == nil {
		return ""
	}

	switch res.Schema.Type {
	case schemaTypeArray:
		return path.Base(res.Schema.Items.Ref.Ref)
	case schemaTypeObject:
		return path.Base(res.Schema.Ref.Ref)
	default:
		return ""
	}
}

func SpecToModelCommand(cfg spec.ModelSpec) (*ModelCommand, error) {
	return &ModelCommand{
		Name:   cfg.Name,
		Fields: cfg.Fields,
	}, nil
}

// deriveReturnType GIN does not have a return type Im leaving this here for reference if I want to add support for a framework that does require a return type in the future. For now it will be used to generate the return code in the handler template.
func deriveReturnType(op spec.Operation) string {
	if op.Responses == nil {
		return ""
	}

	// preferred success codes in order
	successCodes := []int{200, 201, 202, 204}

	var chosen *spec.ResponseSpec

	for _, code := range successCodes {
		if r, ok := op.Responses[code]; ok {
			chosen = &r
			break
		}
	}

	if chosen == nil {
		for code, r := range op.Responses {
			if code >= 200 && code < 300 {
				rr := r
				chosen = &rr
				break
			}
		}
	}

	if chosen == nil || chosen.Schema == nil {
		return ""
	}

	s := chosen.Schema

	switch s.Type {
	case schemaTypeArray:
		key := getModelSchemaKey(*chosen)
		if key == "" {
			return ""
		}
		return "[]models." + key

	case schemaTypeObject, "":
		key := getModelSchemaKey(*chosen)
		if key == "" {
			return ""
		}
		return "models." + key

	default:
		return ""
	}
}
