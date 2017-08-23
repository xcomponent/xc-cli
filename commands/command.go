package commands

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

type Command interface {
	Execute(args []string, awsSession *s3.S3) error
}
