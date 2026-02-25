package templatehelper

import (
	fileHelper "binh-swagger/cmd/swagger/commands/adaptor"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"
)

var ErrTemplateNotFound = errors.New("template not found")

var templatePaths = map[string]string{
	ModelTemplateKey:   "model_template.tmpl",
	HandlerTemplateKey: "handler_template.tmpl",
	RouteTemplateKey:   "route_template.tmpl",
	ServerTemplateKey:  "server_template.tmpl",
}

const (
	HandlerTemplateKey = "handler"
	ModelTemplateKey   = "model"
	RouteTemplateKey   = "route"
	ServerTemplateKey  = "server"
)

type TemplateDefines struct {
	RouteDefine       string
	ImportsDefine     string
	ServerMainDefine  string
	ModelStructDefine string
}

var Templates = TemplateDefines{
	RouteDefine:       "registerRoutes",
	ImportsDefine:     "imports",
	ServerMainDefine:  "serverMain",
	ModelStructDefine: "modelStruct",
}

type TemplateModel interface {
	*ImportsTemplateModel |
		*OperationModel |
		*RoutesTemplateModel |
		*ServerTemplateModel |
		*ModelTemplateModel
}

type RoutesTemplateModel struct {
	Routes []RouteModel
}

type ServerTemplateModel struct {
	Paths            []string
	RoutesImportPath string
}

type RouteModel struct {
	PathName    string
	Method      string
	OperationID string
}

type OperationModel struct {
	MethodType  string
	Summary     string
	OperationID string
	Produces    []string
	ReturnType  string // GIN does not use a return type on handlers Leaving this here for reference if I want to add support for a framework that does require a return type in the future. For now it will be used to generate the return code in the handler template.
	Responses   map[int]ResponseModel
}

type ResponseModel struct {
	Description       string
	Type              string
	Ref               string
	SuccessReturnCode string
}

type ImportsTemplateModel struct {
	ModelImportPath   string
	HandlerImportPath string
	RoutesImportPath  string
}

type ModelTemplateModel struct {
	ImportPath []string
	Name       string
	Fields     []FieldsModel
}

type FieldsModel struct {
	Name string
	Type string
	JSON string
}

func LoadModelTemplate(fileHefileHelper fileHelper.FileHelper, templateName string) (*template.Template, error) {
	templateFilename, err := getTemplatePath(templateName)
	if err != nil {
		return nil, err
	}

	currentDir, err := getCurrentFileDir()
	if err != nil {
		return nil, err
	}

	// path of template is coupled to the location of this file.
	// also not properly tested unit test needs to reflect where this path is
	templatePath := fileHefileHelper.GetAbsoluteSanitiseFilePath(filepath.Join(currentDir, "..", "templates"), templateFilename)

	template := template.Must(template.ParseFiles(templatePath))

	return template, nil

	// old code here as an example of reading file content if needed
	// tmpl, err := os.ReadFile(templatePath)
	// if err != nil {
	// 	return nil, err
	// }
	// return template.New(templateName).Parse(string(tmpl))
}

func getTemplatePath(templateName string) (string, error) {
	path := templatePaths[templateName]
	if path == "" {
		return "", fmt.Errorf("%w: %s", ErrTemplateNotFound, templateName)
	}

	return path, nil
}

// getCurrentFileDir returns the directory of the current file needs to be here to get templates based on relative path.
func getCurrentFileDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to get directory")
	}

	return filepath.Dir(file), nil
}

func GetBaseTemplate(buf *bytes.Buffer, fHelper fileHelper.FileHelper, templateKey string) (*template.Template, error) {
	tmpl, err := LoadModelTemplate(fHelper, templateKey)
	if err != nil {
		return nil, err
	}

	if err = tmpl.ExecuteTemplate(buf, templateKey, nil); err != nil {
		return nil, err
	}

	return tmpl, nil
}

func ExecuteTemplate[T TemplateModel](data T, tmpl *template.Template, buf *bytes.Buffer, templateName string) error {
	if err := tmpl.ExecuteTemplate(buf, templateName, data); err != nil {
		return err
	}

	return nil
}

func SetModelTemplateImportPaths(models *ModelTemplateModel) {
	paths := map[string]string{}

	for _, f := range models.Fields {
		importPath := getImportPath(f.Type)
		if importPath != "" {
			paths[importPath] = importPath
		}

		models.ImportPath = buildImportStringSlice(paths)
	}
}

func buildImportStringSlice(paths map[string]string) []string {
	imports := make([]string, 0, len(paths))
	for _, path := range paths {
		imports = append(imports, path)
	}

	return imports
}

func getImportPath(fieldType string) string {
	switch fieldType {
	case "time.Time", "time.Duration":
		return "time"
	case "uuid.UUID":
		return "github.com/google/uuid"
	case "decimal.Decimal":
		return "github.com/shopspring/decimal"
	case "sql.NullString", "sql.DB", "sql.NullInt64":
		return "database/sql"
	case "url.URL":
		return "net/url"
	case "http.Request":
		return "net/http"
	default:
		return ""
	}
}
