package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/xcomponent/xc-cli/commands"
	"os"
	. "github.com/daniellavoie/go-shim/httpshim"
	. "github.com/daniellavoie/go-shim/zipshim"
	"github.com/xcomponent/xc-cli/services"
	"github.com/daniellavoie/go-shim/execshim"
)

func main() {
	app := cli.NewApp()
	app.Name = "XC CLI"
	app.Usage = "XComponent Command Line Interface"
	app.Version = "0.5.0"

	workDir, err := os.Getwd()
	if err == err {
		app.Commands = commands.GetCommands(
			workDir,
			services.NewOsService(),
			services.NewHttpService(new(HttpShim)),
			services.NewIoService(),
			services.NewZipService(new(ZipShim)),
			services.NewExecService(&execshim.ExecShim{}))

		err = app.Run(os.Args)
	}

	if err != nil {
		fmt.Printf("Error : %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
