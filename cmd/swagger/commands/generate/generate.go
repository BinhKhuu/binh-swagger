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
	ReturnType  string
	Responses   map[int]ResponseCommand
}

type ResponseCommand struct {
	Description       string
	Type              string
	Ref               string
	SuccessReturnCode string
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
			if res.Schema.Type == "array" {
				resCmd.Ref = res.Schema.Items.Ref.Ref
			} else {
				resCmd.Ref = res.Schema.Ref.Ref
			}
			resCmd.Type = res.Schema.Type
			resCmd.Description = res.Description

			setReturnCode(res, mCdm, ref, resCmd)
		}
		responses[responseKey] = *resCmd
	}

	// maybe store the name of the model so handler can generate the return statement with the correct model name

	return responses, nil
}

// todo test this
// coudl combine mCdm and ref into a single parameter since ref is only used to get the model name from mCdm
// change so it loops through all the possible return types and build one return string contraining all the return paths
func setReturnCode(res spec.ResponseSpec, mCdm map[string]*ModelCommand, ref string, resCmd *ResponseCommand) {
	var code string
	switch res.Schema.Type {
	case "array":
		code = fmt.Sprintf("var result = []models.%s{}", mCdm[ref].Name)
	case "object":
		code = fmt.Sprintf("var result = models.%s{}", mCdm[ref].Name)
	default:
		code = "// todo add comment on return type, currently only supports array and object types"
	}

	// todo create string where all the return types are defined, then have the template reference that string instead of hardcoding the return code here. This will make it easier to maintain and update the return code in the future.
	resCmd.SuccessReturnCode = code
}

func getModelSchemaKey(res spec.ResponseSpec) string {
	if res.Schema == nil {
		return ""
	}

	switch res.Schema.Type {
	case "array":
		return path.Base(res.Schema.Items.Ref.Ref)
	case "object":
		return path.Base(res.Schema.Ref.Ref)
	default:
		return ""
	}
}

// SpecToModelCommand todo test this.
func SpecToModelCommand(cfg spec.ModelSpec) (*ModelCommand, error) {
	return &ModelCommand{
		Name:   cfg.Name,
		Fields: cfg.Fields,
	}, nil
}

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
	case "array":
		key := getModelSchemaKey(*chosen)
		if key == "" {
			return ""
		}
		return "[]models." + key

	case "object", "":
		key := getModelSchemaKey(*chosen)
		if key == "" {
			return ""
		}
		return "models." + key

	default:
		return ""
	}
}
