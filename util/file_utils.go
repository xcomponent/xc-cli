package util

import (
	"github.com/xcomponent/xc-cli/services"
	"strings"
	"fmt"
	"os"
	"path/filepath"
	"errors"
	"io"
)

type FileUtils interface {
	ListFiles(folder string, predicates []StringPredicate) ([]string, error)
	RenameAndReplaceFiles(folder string, old string, new string) error
	ReplaceInFileName(filePath string, old string, new string) (string, error)
	SearchAndReplace(fileName string, old string, new string) error
}

func NewFileUtils(os services.OsService, io services.IoService) *FileUtilsImpl {
	return &FileUtilsImpl{os, io}
}

type FileUtilsImpl struct {
	os services.OsService
	io services.IoService
}

func CopyFile(source string, desination string, perm os.FileMode) error {
	from, err := os.Open(source)
	if from != nil {
		defer from.Close()
	}
	if err != nil {
		return err
	}

	to, err := os.OpenFile(desination, os.O_RDWR|os.O_CREATE, perm)
	if to != nil{
		defer to.Close()
	}
	if err != nil {
		return err
	}

	_, err = io.Copy(to, from)

	return err
}

func (fileUtils *FileUtilsImpl) ListFiles(folder string, predicates []StringPredicate) ([]string, error) {
	fileList := make([]string, 0)
	err := filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
		for _, predicate := range predicates {
			if predicate(path) {
				fileList = append(fileList, path)
				break
			}
		}

		return err
	})
	if err != nil {
		return nil, err
	}

	return fileList, nil
}

func (fileUtils *FileUtilsImpl) SearchAndReplace(fileName string, old string, new string) error {
	bytes, err := fileUtils.io.ReadFile(fileName)
	if err != nil {
		return err
	}

	newContents := strings.Replace(string(bytes), old, new, -1)

	err = fileUtils.io.WriteFile(fileName, []byte(newContents), 0)
	if err != nil {
		return err
	}

	return nil
}

func (fileUtils *FileUtilsImpl) RenameAndReplaceFiles(folder string, old string, new string) error {
	renameVisitor := renameVisitor{fileUtils, folder, old, new, nil}
	replaceVisitor := replaceVisitor{fileUtils, old, new}

	filepath.Walk(folder, renameVisitor.visit)

	if renameVisitor.err != nil {
		return renameVisitor.err
	}

	err := filepath.Walk(folder, replaceVisitor.visit)

	if err != nil {
		return err
	}

	return nil
}

func (fileUtils *FileUtilsImpl) ReplaceInFileName(filePath string, old string, new string) (string, error) {
	pathSeperator := fmt.Sprintf("%c", fileUtils.os.GetPathSeperator())

	paths := strings.Split(filePath, pathSeperator)

	filename := paths[len(paths)-1]

	if !strings.Contains(filename, old) {
		return filePath, nil
	}

	newFilename := fmt.Sprintf(
		"%s%s", filePath[0:len(filePath)-len(filename)],
		strings.Replace(filename, old, new, -1))

	err := fileUtils.os.Rename(filePath, newFilename)
	if err != nil {
		return "", err
	}

	return newFilename, nil
}

type renameVisitor struct {
	fileUtils *FileUtilsImpl
	rootPath  string
	old       string
	new       string
	err       error
}

func (visitor *renameVisitor) visit(path string, fi os.FileInfo, err error) error {
	newPath, err := visitor.fileUtils.ReplaceInFileName(path, visitor.old, visitor.new)
	if err != nil {
		return visitor.returnError(err)
	}

	if newPath != path {
		newVisitor := renameVisitor{visitor.fileUtils, visitor.rootPath, visitor.old, visitor.new, nil}
		filepath.Walk(visitor.rootPath, newVisitor.visit)

		if newVisitor.err != nil {
			return visitor.returnError(newVisitor.err)
		}

		return errors.New("file renamed")
	}

	return nil
}

func (visitor *renameVisitor) returnError(err error) error {
	visitor.err = err

	return err
}

type replaceVisitor struct {
	fileUtils *FileUtilsImpl
	old       string
	new       string
}

func (visitor *replaceVisitor) visit(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !!fi.IsDir() {
		return nil
	}

	return visitor.fileUtils.SearchAndReplace(path, visitor.old, visitor.new)
}
