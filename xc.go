package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"xcomponent.com/xc/cli"
	"xcomponent.com/xc/commands"
	"xcomponent.com/xc/service"
)

const (
	Command = "xc"
)

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 {
		log.Fatalf("Please specify an argument.\n%s", buildUsage())
	}

	commandName := argsWithoutProg[0]

	config, err := cli.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	componentService := service.NewComponentService(config.ServerUrl)
	service := commands.NewCommandService(config.ServerUrl)

	command := map[string]commands.Command{
		"push": commands.PushCommand{ComponentService: componentService},
		"run":  commands.RunCommand{Service: service},
		"logs":  commands.LogsCommand{Service: service},
	}[commandName]

	if command != nil {
		err := command.Execute(argsWithoutProg[1:])

		if err != nil {
			fmt.Printf("Error : %s\n", err)
		}
	} else {
		log.Fatalf("Invalid command.\n%s", buildUsage())
	}
}

func buildUsage() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n")
	buffer.WriteString("Usage : " + cli.Command + " COMMAND\n")
	buffer.WriteString("\n")
	buffer.WriteString("A command line interface for XComponent Cloud Platform\n")
	buffer.WriteString("\n")
	buffer.WriteString("Commands:\n")
	buffer.WriteString("	 push      Upload an XC component to XC Cloud Platform\n")
	buffer.WriteString("	 run       Launch a component stack based on docker-compose\n")
	buffer.WriteString("	 logs      Access the logs of a service instance\n")

	return buffer.String()
}
