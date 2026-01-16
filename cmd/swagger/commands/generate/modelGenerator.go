package generate

import (
	"binh-swagger/cmd/swagger/commands"
	"bytes"
	"html/template"
	"os"
	"path/filepath"
)

const (
	// filePermOwnerReadWrite is the file permission for owner read/write only (rw-------).
	filePermOwnerReadWrite = 0o600
)

func Model(config commands.ModelSpec) error {
	var buf bytes.Buffer
	tmpl, err := loadModelTemplate()
	if err != nil {
		return err
	}

	err = tmpl.Execute(&buf, config)
	if err != nil {
		return err
	}
	outputFile := filepath.Join(config.OutputPath, config.OutputFile)
	return os.WriteFile(outputFile, buf.Bytes(), filePermOwnerReadWrite)
}

func loadModelTemplate() (*template.Template, error) {
	tmpl, err := os.ReadFile("./templates/model_template.tmpl")
	if err != nil {
		return nil, err
	}

	return template.New("./templates").Parse(string(tmpl))
}
