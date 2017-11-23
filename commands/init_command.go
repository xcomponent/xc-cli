package commands

import (
	. "os"
	"fmt"
	"errors"
	"path/filepath"
	"strings"
	"github.com/xcomponent/xc-cli/services"
	"io"
	"archive/zip"
	"github.com/xcomponent/xc-cli/util"
)

type InitCommand struct {
	http services.HttpService
	io   services.IoService
	os   services.OsService
	zip  services.ZipService
}

func NewInitCommand(http services.HttpService, io services.IoService, os services.OsService,
	zip services.ZipService) *InitCommand {
	return &InitCommand{http, io, os, zip}
}

func (initCommand *InitCommand) Init(
	projectFolder string, githubOrg, templateName string, projectName string) error {

	err := initCommand.prepareProjectFolder(projectFolder)

	if err != nil {
		return err
	}

	fileName, err := initCommand.downloadTemplate(projectFolder, githubOrg, templateName)
	defer initCommand.clean(fileName)

	if err != nil {
		return err
	}

	err = initCommand.unzip(fileName, projectFolder)
	if err != nil {
		return err
	}

	if projectName != "" {
		return util.NewFileUtils(initCommand.os, initCommand.io).RenameAndReplaceFiles(
			projectFolder, "NewProject", projectName)
	}

	return nil
}

func (initCommand *InitCommand) prepareProjectFolder(projectFolder string) error {
	fileInfo, err := initCommand.os.Stat(projectFolder)
	if err != nil {
		if initCommand.os.IsNotExist(err) {
			return initCommand.os.MkdirAll(projectFolder, 0700)
		} else {
			return err
		}
	} else {
		if !fileInfo.IsDir() {
			return errors.New(fmt.Sprintf("%s is not a directory", projectFolder))
		}

		empty, err := initCommand.isEmpty(projectFolder)
		if err != nil {
			return err
		} else if !empty {
			return errors.New(fmt.Sprintf("%s is not empty", projectFolder))
		}
	}

	return nil
}

func (initCommand *InitCommand) isEmpty(name string) (bool, error) {
	f, err := initCommand.os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func (initCommand *InitCommand) downloadTemplate(projectFolder string, githubOrg string, templateName string) (string, error) {
	fmt.Println("Downloading template")

	fileName := fmt.Sprintf("%s/%s.zip", projectFolder, templateName)

	out, err := initCommand.os.Create(fileName)
	defer out.Close()
	if err != nil {
		return fileName, err
	}

	resp, err := initCommand.http.Get(fmt.Sprintf("https://github.com/%s/%s/archive/master.zip", githubOrg, templateName))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fileName, err
	}

	if resp.StatusCode == 404 {
		return fileName, errors.New(fmt.Sprintf("Project template %s/%s could not be found", githubOrg, templateName))
	}

	_, err = initCommand.io.Copy(out, resp.Body)
	if err != nil {
		return fileName, err
	}

	return fileName, nil
}

func (initCommand *InitCommand) unzip(archive string, target string) error {
	fmt.Println(fmt.Sprintf("Extracting template from %s", archive))

	reader, err := initCommand.zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := initCommand.os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		err = initCommand.unzipFile(file, target)
		if err != nil {
			return err
		}
	}

	return nil
}

func (initCommand *InitCommand) unzipFile(file *zip.File, target string) error {
	fileName := file.Name[strings.Index(file.Name, "/")+1:len(file.Name)]

	path := filepath.Join(target, fileName)
	if file.FileInfo().IsDir() {
		initCommand.os.MkdirAll(path, file.Mode())
		return nil
	}

	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	targetFile, err := initCommand.os.OpenFile(path, O_WRONLY|O_CREATE|O_TRUNC, file.Mode())
	if targetFile != nil {
		defer targetFile.Close()
	}
	if err != nil {
		return err
	}

	if _, err := initCommand.io.Copy(targetFile, fileReader); err != nil {
		return err
	}

	return nil
}

func (initCommand *InitCommand) clean(fileName string) error {
	fmt.Println(fmt.Sprintf("Cleaning %s", fileName))

	return initCommand.os.Remove(fileName)
}
