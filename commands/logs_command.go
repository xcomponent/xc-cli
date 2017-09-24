package commands

import (
	"github.com/daniellavoie/xc-cli/service"
)

type LogsCommandRequest struct {
	StackName   string `json:"stackName"`
	ServiceName string `json:"serviceName"`
	Instance    string `json:"instance"`
}

func RetreiveLogs(request LogsCommandRequest, service service.CloudPlatformService) error {
	_, err := service.SendCommand("logs", request)

	return err
}
