package filehelper

import (
	"binh-swagger/cmd/swagger/commands/helpers"
	"io"
	"os"
)

type FileHelper interface {
	ValidateFileInfo(fileInfo os.FileInfo) error
	GetAbsoluteSanitiseFilePath(filelocation string, filename string) string
	CheckSymlinks(absFilePath string) error
	EnsureOutputDirectoryExists(outputPath string, output io.Writer, input io.Reader) (bool, error)
}

type DefaultFileHelper struct{}

func (f *DefaultFileHelper) ValidateFileInfo(fileInfo os.FileInfo) error {
	return helpers.ValidateFileInfo(fileInfo)
}

func (f *DefaultFileHelper) GetAbsoluteSanitiseFilePath(filelocation string, filename string) string {
	return helpers.GetAbsoluteSanitiseFilePath(filelocation, filename)
}

func (f *DefaultFileHelper) CheckSymlinks(absFilePath string) error {
	return helpers.CheckSymlinks(absFilePath)
}

func (f *DefaultFileHelper) EnsureOutputDirectoryExists(outputPath string, output io.Writer, input io.Reader) (bool, error) {
	return helpers.EnsureOutputDirectoryExists(outputPath, output, input)
}
