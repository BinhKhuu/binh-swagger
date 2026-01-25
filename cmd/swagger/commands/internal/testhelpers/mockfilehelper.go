package testhelpers

import (
	"binh-swagger/cmd/swagger/commands/helpers"
	"io"
	"os"
)

type MockFileHelper struct {
	ValidateFileInfoFn                 func(os.FileInfo) error
	GetAbsoluteSanitiseFilePathFn      func(string, string) string
	CheckSymlinksFn                    func(string) error
	EnsureOutputDirectoryExistsFn      func(string, io.Writer, io.Reader) (bool, error)
	CreateDirectoryFn                  func(string) (string, error)
	ReadAllChildDirectoriesRecursiveFn func(string) ([]string, error)
}

func (m MockFileHelper) ValidateFileInfo(fi os.FileInfo) error {
	return m.ValidateFileInfoFn(fi)
}

func (m MockFileHelper) GetAbsoluteSanitiseFilePath(loc, name string) string {
	return m.GetAbsoluteSanitiseFilePathFn(loc, name)
}

func (m MockFileHelper) CheckSymlinks(path string) error {
	return m.CheckSymlinksFn(path)
}

func (m MockFileHelper) EnsureOutputDirectoryExists(out string, w io.Writer, r io.Reader) (bool, error) {
	return m.EnsureOutputDirectoryExistsFn(out, w, r)
}

func (m MockFileHelper) CreateDirectory(root string) (string, error) {
	return m.CreateDirectoryFn(root)
}

func (m MockFileHelper) ReadAllChildDirectoriesRecursive(_ string) ([]string, error) {
	return []string{}, nil
}

func CreateMockFileHelper() MockFileHelper {
	mock := MockFileHelper{
		ValidateFileInfoFn: func(_ os.FileInfo) error {
			return nil
		},
		GetAbsoluteSanitiseFilePathFn: helpers.GetAbsoluteSanitiseFilePath,
		CheckSymlinksFn: func(_ string) error {
			return nil
		},
		EnsureOutputDirectoryExistsFn: func(_ string, _ io.Writer, _ io.Reader) (bool, error) {
			return true, nil
		},
		CreateDirectoryFn: func(_ string) (string, error) {
			return "", nil
		},
	}

	return mock
}
