package commands

import (
	"bytes"
	"errors"
	"fmt"
	"encoding/json"
)

type LogsCommandRequest struct {
	StackName   string `json:"stackName"`
	ServiceName string `json:"serviceName"`
	Instance    string `json:"instance"`
}

type LogsCommand struct {
	Service *ServiceOp
}

func (logsCommand LogsCommand) Execute(args []string) error {
	logsCommandRequest, err := loadLogsParams(args)

	if err != nil {
		return err
	}

	res, err := logsCommand.Service.ExecuteCommand("logs", logsCommandRequest)
	if err == nil {
		var output Output
		json.NewDecoder(res.Body).Decode(&output)
		fmt.Println(output.Content)
	}

	return err
}

func loadLogsParams(args []string) (*LogsCommandRequest, error) {
	if len(args) != 2 && len(args) != 3 {
		return nil, errors.New(fmt.Sprintf("Invalid parameters specified.\n%s", printLogsUsage()))
	}

	instance := "1"
	if len(args) == 3 {
		instance = args[2]
	}

	return &LogsCommandRequest{StackName: args[0], ServiceName: args[1], Instance: instance}, nil
}

func printLogsUsage() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n")
	buffer.WriteString("Usage: xc logs STACK-NAME SERVICE-NAME [INSTANCE-NUMBER]\n")
	buffer.WriteString("\n")

	return buffer.String()
}
