package commands

import (
	"errors"
	"os"
	"strings"
	"github.com/daniellavoie/xc-cli/service"
)

type PushRequest struct {
	componentName string
	label         string
	filePath      string
}

func BuildPushRequest(componentNameAndLabel string, filePath string) (*PushRequest, error) {
	componentNameAndLabelArray := strings.Split(componentNameAndLabel, ":")
	if len(componentNameAndLabel) > 2 {
		return nil, errors.New("Component name parameter does not match the pattern \"component-name:label\".")
	}

	componentName := componentNameAndLabelArray[0]
	var label = "latest"
	if len(componentNameAndLabel) == 2 {
		label = componentNameAndLabelArray[1]
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.New(filePath + " could not be found.")
	}

	return &PushRequest{componentName, label, filePath}, nil
}

func PushComponent(request *PushRequest, service service.CloudPlatformService) error {
	_, err := service.UploadComponent(request.componentName, request.label, request.filePath)

	return err
}
