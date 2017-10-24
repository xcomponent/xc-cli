package commands

import (
	"io/ioutil"
	"github.com/xcomponent/xc-cli/service"
)

type RunStackRequest struct {
	StackName       string `json:"stackName"`
	StackDefinition string `json:"stackDefinition"`
}

func BuildRunStackRequest(stackName string, stackDefinitionPath string) (*RunStackRequest, error) {
	stackDefinition, err := ioutil.ReadFile(stackDefinitionPath)
	if err != nil {
		return nil, err
	}

	return &RunStackRequest{StackName: stackName, StackDefinition: string(stackDefinition)}, nil
}

func RunStack(request *RunStackRequest, service service.CloudPlatformService) error {
	_, err := service.SendCommand("run", request)

	return err
}
