package commands

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"github.com/xcomponent/xc-cli/service"
)

type PushCommand struct {
	ComponentService *service.ComponentServiceOp
}

func (pushCommand PushCommand) Execute(args []string) error {
	pushParams, err := loadPushParams(args)

	if err != nil {
		return err
	}

	_, err = pushCommand.ComponentService.Upload(pushParams.componentName, pushParams.label, pushParams.filePath)

	return err
}

type PushParams struct {
	componentName string
	label         string
	filePath      string
}

func loadPushParams(args []string) (*PushParams, error) {
	argsLength := len(args)
	if argsLength != 2 {
		return nil, errors.New(fmt.Sprintf("Invalid parameters specified.\n%s", printPushUsage()))
	}

	componentNameAndLabel := strings.Split(args[0], ":")

	if len(componentNameAndLabel) > 2 {
		return nil, errors.New("Component name parameter does not match the pattern \"component-name:label\".")
	}

	componentName := componentNameAndLabel[0]
	var label = "latest"
	if len(componentNameAndLabel) == 2 {
		label = componentNameAndLabel[1]
	}

	filePath := args[1]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.New(filePath + " could not be found.")
	}

	return &PushParams{componentName, label, filePath}, nil
}

func printPushUsage() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n")
	buffer.WriteString("Usage: xc push COMPONENT-NAME[:LABEL] COMPONENT-PATH\n")
	buffer.WriteString("\n")

	return buffer.String()
}
