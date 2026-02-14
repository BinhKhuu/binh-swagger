package generate

import (
	templateHelper "binh-swagger/cmd/swagger/commands/generate/internal/templateHelper"
	"binh-swagger/cmd/swagger/commands/internal/pkg"
	"bytes"
	"os"
	"strings"
)

func Model(cmd *ModelCommand, generateConfig Config) error {
	if _, err := GetProjectStructure(); err != nil {
		return err
	}

	var buf bytes.Buffer
	fHelper := generateConfig.FileHelper()
	tmpl, err := templateHelper.GetBaseTemplate(&buf, generateConfig.FileHelper(), templateHelper.ModelTemplateKey)
	if err != nil {
		return err
	}

	tModel := toModelTemplate(cmd)
	if err = templateHelper.ExecuteTemplate(&tModel, tmpl, &buf, templateHelper.Templates.ImportsDefine); err != nil {
		return err
	}

	if err = templateHelper.ExecuteTemplate(&tModel, tmpl, &buf, templateHelper.Templates.ModelStructDefine); err != nil {
		return err
	}

	err = tmpl.Execute(&buf, cmd)
	if err != nil {
		return err
	}

	outputPath := projectStructure["models"]
	outFile := strings.ToLower(cmd.Name) + ".go"

	outputFile := fHelper.GetAbsoluteSanitiseFilePath(outputPath, outFile)
	return os.WriteFile(outputFile, buf.Bytes(), pkg.FilePermOwnerReadWrite)
}

func toModelTemplate(cmd *ModelCommand) templateHelper.ModelTemplateModel {
	fields := []templateHelper.FieldsModel{}
	for _, f := range cmd.Fields {
		if f.Name == "" || f.Type == "" {
			continue
		}
		fm := templateHelper.FieldsModel{
			Name: f.Name,
			Type: f.Type,
			JSON: f.JSON,
		}
		fields = append(fields, fm)
	}
	models := templateHelper.ModelTemplateModel{
		Name:   cmd.Name,
		Fields: fields,
	}

	templateHelper.SetModelTemplateImportPaths(&models)
	return models
}
