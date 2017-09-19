package commands

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/xcomponent/xc-cli/cli"
	"os"
	"strings"
)

type PushCommand struct {
	BucketName string
}

func (pushCommand PushCommand) Execute(args []string, s3 *s3.S3) error {
	pushParams, err := parseParams(args)

	if err != nil {
		return err
	}

	uploader := s3manager.NewUploaderWithClient(s3)

	file, err := os.Open(pushParams.filePath)

	if err != nil {
		return nil
	}

	key := pushParams.componentName + "-" + pushParams.label

	// Perform an upload.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: &pushCommand.BucketName,
		Key:    &key,
		Body:   file,
	})

	return err
}

type PushParams struct {
	componentName string
	label         string
	filePath      string
}

func parseParams(args []string) (*PushParams, error) {
	if len(args) != 2 {
		return nil, errors.New("\"" + cli.Command + " push\" requires two arguments.")
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
