package main

import (
	"fmt"
	"os"
	"xcomponent.com/xc/commands"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

const (
	Command = "xc"
)

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 {
		fmt.Printf("Please specify an argument.\n")

		return
	}

	commandName := argsWithoutProg[0]

	commands := map[string]commands.Command{
		"push": commands.PushCommand{"xc-cloud-dev"},
	}

	command := commands[commandName]

	if command != nil {
		awsSession, err := NewAwsSession()

		if err == nil {
			err := command.Execute(argsWithoutProg[1:], s3.New(awsSession))

			if err != nil {
				fmt.Printf("Error : %s\n", err)
			}else {
				fmt.Printf("Uploaded component to AWS.\n")
			}

		} else {
			fmt.Printf("Error while connecting to AWS %s\n", err)

			return
		}
	} else {
		fmt.Printf("Invalid argument %s\n", commandName)
	}
}

func NewAwsSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: credentials.NewSharedCredentials("", "xc-cloud"),
	})
}
