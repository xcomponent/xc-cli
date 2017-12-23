package commands

import (
	"fmt"
	"github.com/xcomponent/xc-cli/util"
	"github.com/xcomponent/xc-cli/services"
	"errors"
	"path"
)

func NewAddCommand(os services.OsService, http services.HttpService, io services.IoService,
	zip services.ZipService) *AddCommand {

	return &AddCommand{os, http, io, zip}
}

type AddCommand struct {
	os   services.OsService
	http services.HttpService
	io   services.IoService
	zip  services.ZipService
}

func (command *AddCommand) Execute(projectFolder string, elementType string, elementName string, addTemplatesGithubOrg string) error {
	var placeholder = ""
	if elementType == "component" {
		placeholder = "NewComponent"
	} else {
		return errors.New(fmt.Sprintf("%s is not a supported element to add to the project", elementType))
	}

	err := util.NewGitHubUtils(command.os, command.http, command.io, command.zip).DownloadTemplate(
		projectFolder, addTemplatesGithubOrg, elementType, elementName)
	if err != nil {
		return err
	}

	var elementPath = path.Join(projectFolder, placeholder)

	if elementName != "" {
		elementPath = path.Join(projectFolder, elementName)

		err = command.os.Rename(path.Join(projectFolder, placeholder), elementPath)
		if err != nil {
			return err
		}

		err = util.NewFileUtils(command.os, command.io).RenameAndReplaceFiles(
			elementPath, placeholder, elementName)

		if err != nil {
			return err
		}
	}

	fmt.Printf("Successfully added %s %s to %s.\n", elementType, elementName, elementPath)

	return nil

}
