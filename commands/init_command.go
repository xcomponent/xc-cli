package commands

import (
	"fmt"
	"errors"
	"github.com/xcomponent/xc-cli/services"
	"io"
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

	err = util.NewGitHubUtils(initCommand.os, initCommand.http, initCommand.io, initCommand.zip).DownloadTemplate(
		projectFolder, githubOrg, templateName, projectName)
	if err != nil {
		return err
	}

	if projectName != "" {
		err = util.NewFileUtils(initCommand.os, initCommand.io).RenameAndReplaceFiles(
			projectFolder, "NewProject", projectName)
		if err != nil {
			return err
		}
	}else {
		projectName = "NewProject"
	}

	fmt.Printf("Successfully initialized project %s in %s.\n", projectName,projectFolder)

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
