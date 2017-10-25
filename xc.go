package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/xcomponent/xc-cli/commands"
	"os"
	. "github.com/daniellavoie/go-shim/osshim"
	. "github.com/daniellavoie/go-shim/httpshim"
	. "github.com/daniellavoie/go-shim/ioshim"
)

func main() {
	app := cli.NewApp()
	app.Name = "XC CLI"
	app.Usage = "XComponent Command Line Interface"
	app.Version = "0.2.0"

	os.TempDir()
	app.Commands = commands.GetCommands(new(OsShim), new(HttpShim), new(IoShim))

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
