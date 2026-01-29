package generate

import (
	templateHelper "binh-swagger/cmd/swagger/commands/generate/internal"
	"bytes"
	"os"
	"strings"
)

const (
	filePermOwnerReadWrite = 0o600
)

func Model(cmd *ModelCommand, generateConfig Config) error {
	var buf bytes.Buffer
	fHelper := generateConfig.FileHelper()
	tmpl, err := templateHelper.LoadModelTemplate(fHelper, templateHelper.ModelTemplateKey)
	if err != nil {
		return err
	}

	err = tmpl.Execute(&buf, cmd)
	if err != nil {
		return err
	}

	// todo think about coupling to GetProjectStructure here
	_, err = GetProjectStructure()
	if err != nil {
		return err
	}
	outputPath := projectStructure["models"]
	outFile := strings.ToLower(cmd.Name) + ".go"

	outputFile := fHelper.GetAbsoluteSanitiseFilePath(outputPath, outFile)
	shouldReturn, err := fHelper.EnsureOutputDirectoryExists(outputPath, os.Stdout, os.Stdin)
	if shouldReturn {
		return err
	}

	return os.WriteFile(outputFile, buf.Bytes(), filePermOwnerReadWrite)
}
