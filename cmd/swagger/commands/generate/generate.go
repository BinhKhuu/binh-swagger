package generate

import (
	filehelper "binh-swagger/cmd/swagger/commands/adaptor"
	"binh-swagger/cmd/swagger/commands/internal/spec"
	"errors"
	"fmt"
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
	Responses   map[int]ResponseCommand
}

// combining schema spec struct with responsecommand spec
type ResponseCommand struct {
	Description string
	Type        string
	Ref         string
}

type ServerCommand struct {
	PathNames []string
}

// func SpecToPathCommand(cfg spec.PathSpec, pathName string) (*PathCommand, error) {
// 	paths, err := GetProjectStructure()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &PathCommand{
// 		ModelImportPath:   filepath.Join(projectImportPath, paths["models"]),
// 		HandlerImportPath: filepath.Join(projectImportPath, paths["handlers"]),
// 		Name:              strings.Replace(pathName, "/", "", 1),
// 		Get:               cfg.Get,
// 		Post:              cfg.Post,
// 		Put:               cfg.Put,
// 		Patch:             cfg.Patch,
// 		Delete:            cfg.Delete,
// 	}, nil
// }

// todo this will replace SpecToPathCommand
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
	return opCmd, nil
}

/*
Look up $ref in the reponses then look up the definitions path and get the model
IF model is there generate the response
IF model is not there then return error
The success response is what is used to generate the return in the body of the handler function the other repsones should be commended out in the handler body.
*/
func SpecToOperationResponses(cfg spec.Operation, mCdm map[string]*ModelCommand) (map[int]ResponseCommand, error) {
	responses := make(map[int]ResponseCommand)

	// for each response get the Description, Schema.Type, Schema.Ref
	for responseKey, res := range cfg.Responses {
		// if reference is found in ModelCommand then we can generate the response otherwise we return an error
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
			resCmd.Type = res.Schema.Type
			resCmd.Ref = res.Schema.Ref
			resCmd.Description = res.Description
		}
		responses[responseKey] = *resCmd
	}

	// maybe store the name of the model so handler can generate the return statement with the correct model name

	return responses, nil
}

func getModelSchemaKey(res spec.ResponseSpec) string {
	// ref := path.Base(res.Schema.Ref)
	if res.Schema == nil {
		return ""
	}

	ref := res.Schema.Ref
	return ref
}

// SpecToModelCommand todo test this.
func SpecToModelCommand(cfg spec.ModelSpec) (*ModelCommand, error) {
	return &ModelCommand{
		Name:   cfg.Name,
		Fields: cfg.Fields,
	}, nil
}
