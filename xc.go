package main

import (
	"fmt"
	"os"
	xccli "github.com/xcomponent/xc-cli/cli"
	"github.com/urfave/cli"
	"github.com/xcomponent/xc-cli/commands"
	"log"
	"github.com/xcomponent/xc-cli/service"
)

const (
	Command          = "xc"
	InstallConfigUrl = "https://raw.githubusercontent.com/daniellavoie/xc-cli/install-config-v1/install-config.json"
)

func main() {
	config, err := xccli.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := cli.NewApp()
	app.Name = "XC CLI"
	app.Usage = "XComponent Command Line Interface"
	app.Version = "0.1.0"

	app.Commands = commands.GetCommands(InstallConfigUrl, service.NewCloudPlatformService(config.ServerUrl))

	err = app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
