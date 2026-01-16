package generate

import (
	"binh-swagger/cmd/swagger/commands"
	"bytes"
	"html/template"
	"os"
	"path/filepath"
)

func GenerateModel(config commands.ModelSpec) error {
	var buf bytes.Buffer
	tmpl, err := loadModelTemplate()
	if err != nil {
		return err
	}

	err = tmpl.Execute(&buf, config)
	if err != nil {
		return err
	}
	outputDir := filepath.Join("..", "testdata/output")
	outputFile := filepath.Join(outputDir, "models.go")
	return os.WriteFile(outputFile, buf.Bytes(), 0o644)
}

func loadModelTemplate() (*template.Template, error) {
	tmpl, err := os.ReadFile("./templates/model_template.tmpl")
	if err != nil {
		return nil, err
	}

	return template.New("./templates").Parse(string(tmpl))
}
