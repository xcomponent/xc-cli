package commands

import (
	"os"
	"strings"
	"errors"
	"fmt"
	"github.com/xcomponent/xc-cli/util"
	"path"
	"io/ioutil"
	"path/filepath"
	"github.com/xcomponent/xc-cli/services"
)

const projectConfigFileSuffix = "_Model.xcml"

type RenameProjectCommand struct {
	io        services.IoService
	os        services.OsService
	fileUtils util.FileUtils
}

type RenameProjectProvider interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
	Rename(oldpath string, newpath string) error
	Stat(name string) (os.FileInfo, error)
}

type RenameProjectProviderImpl struct {
}

func (provider *RenameProjectProviderImpl) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (provider *RenameProjectProviderImpl) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (provider *RenameProjectProviderImpl) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func NewRenameProjectCommand(io services.IoService, os services.OsService, fileUtils util.FileUtils) *RenameProjectCommand {
	return &RenameProjectCommand{io, os, fileUtils}
}

func (command *RenameProjectCommand) Execute(workDir string, newProjectName string) error {
	projectName, err := command.extractProjectNameFromWorkDir(workDir)
	if err != nil {
		return err
	}

	files, err := command.fileUtils.ListFiles(workDir, []util.StringPredicate{func(value string) bool {
		return true
	}})
	for _, file := range files {
		command.fileUtils.SearchAndReplace(file, projectName, newProjectName)
	}

	// List all files that needs to be renamed.
	files, err = command.listFilesToRename(workDir, projectName)

	return command.renameFiles(files, projectName, newProjectName)

	/*
edit strongProject_Model.xcml and set LinkingSchema/@name to strongProject
rename folder Configuration.dummyProject with Configuration.strongProject
rename each file "dummyProject_Deployment_Configuration.xml" by "strongProject_Deployment_Configuration.xml"
	 */

	return nil
}

func (command *RenameProjectCommand) extractProjectNameFromWorkDir(workDir string) (projectName string, err error) {
	files, err := command.io.ReadDir(workDir)
	if err != nil {
		return projectName, err
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), projectConfigFileSuffix) {
			return ExtractProjectName(f.Name()), nil
		}
	}

	return projectName, errors.New(fmt.Sprintf("Could not find a project configuration file in %s", workDir))
}

/*
func (command *RenameProjectCommand) renameProjectConfigFile(workDir string, oldProjectName string, newProjectName string) (newFileName string, err error) {
	path.Join(workDir, strings.Join([]string{workDir, oldProjectName, projectConfigFileSuffix}, ""))
	oldProjectConfigFileName := strings.Join([]string{workDir, oldProjectName, projectConfigFileSuffix}, "")
	newFileName = strings.Join([]string{newProjectName, projectConfigFileSuffix}, "")

	err = command.os.Rename(oldProjectConfigFileName, newFileName)
	if err != nil {
		return "", err
	}

	return newFileName, nil
}*/

func ExtractProjectName(projectConfigFileName string) string {
	return projectConfigFileName[0: len(projectConfigFileName)-len(projectConfigFileSuffix)]
}

/*
func (command *RenameProjectCommand) editFiles(workDir string, oldProjectName string, newProjectName string) error {
	return command.fileUtils.SearchAndReplace(workDir, oldProjectName, newProjectName)
}
*/

func (command *RenameProjectCommand) listFilesToRename(workDir string, oldProjectName string) ([]string, error) {
	return command.fileUtils.ListFiles(workDir, []util.StringPredicate{
		func(value string) bool {
			return strings.Contains(filepath.Base(value), oldProjectName)
		},
	})
}

func (command *RenameProjectCommand) renameFiles(files []string, oldProjectName string, newProjectName string) error {
	fileToProcess := files
	completed := false

	for !completed {
		for index, file := range fileToProcess {
			newFilePath := path.Join(filepath.Dir(file), strings.Replace(filepath.Base(file), oldProjectName, newProjectName, -1))
			err := command.os.Rename(file, newFilePath)

			if err != nil {
				return err
			}

			fileInfo, err := command.os.Stat(newFilePath)
			if err != nil {
				return err
			}

			if fileInfo.Mode().IsDir() {
				fileToProcess = fileToProcess[index+1:]
				for indexToRename, leftOverFile := range fileToProcess {
					if strings.HasPrefix(leftOverFile, file) {
						fileToProcess[indexToRename] = path.Join(newFilePath, leftOverFile[len(file)+1:])
					}
				}

				break
			}

			completed = true
		}
	}

	return nil
}
