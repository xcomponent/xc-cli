package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
	"github.com/xcomponent/xc-cli/commands"
)

func main() {
	app := cli.NewApp()
	app.Name = "XC CLI"
	app.Usage = "XComponent Command Line Interface"
	app.Version = "0.2.0"

	app.Commands = commands.GetCommands()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
