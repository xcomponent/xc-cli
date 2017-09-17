package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
)

type RunCommandRequest struct {
	StackName       string `json:"stackName"`
	StackDefinition string `json:"stackDefinition"`
}

type RunCommand struct {
	Service *ServiceOp
}

func (runCommand RunCommand) Execute(args []string) error {
	runCommandRequest, err := loadRunParams(args)

	if err != nil {
		return err
	}

	_, err = runCommand.Service.ExecuteCommand("run", runCommandRequest)

	return err
}

func loadRunParams(args []string) (*RunCommandRequest, error) {
	if len(args) != 2 {
		return nil, errors.New(fmt.Sprintf("Invalid parameters specified.\n%s", printRunUsage()))
	}

	stackDefinition, err := ioutil.ReadFile(args[1])
	if err != nil {
		return nil, err
	}

	return &RunCommandRequest{StackName: args[0], StackDefinition: string(stackDefinition)}, nil
}

func printRunUsage() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n")
	buffer.WriteString("Usage: xc run STACK-NAME DOCKER-COMPOSE-YML-PATH\n")
	buffer.WriteString("\n")

	return buffer.String()
}
