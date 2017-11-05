package commands

import (
	"github.com/urfave/cli"
	"strings"
	"os"
	"github.com/daniellavoie/go-shim"
)

const (
	installConfigUrl = "https://raw.githubusercontent.com/xcomponent/xc-cli/install-config-v1/install-config.json"

	githubOrg = "xcomponent-templates"

)

func GetCommands(osshim goshim.Os, httpshim goshim.Http, io goshim.Io) ([]cli.Command) {
	return []cli.Command{
		{
			Name:  "install",
			Usage: "Install XComponent",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "keep-temp-files", Usage: "Temporary files will not be cleaned."},
				cli.StringFlag{
					Name:  "install-config-url",
					Usage: "Url for the dependencies manifest.",
					Value: installConfigUrl},
			},
			Action: func(c *cli.Context) error {
				return Install(
					c.String("install-config-url"),
					false,
					c.Bool("keep-temp-files"))
			},
		},
		{
			Name:  "init",
			Usage: "Initialize a new XComponent project",
			ArgsUsage: "[ORGANIZATION:][TEMPLATE-NAME]",
			Action: func(c *cli.Context) error {
				if len(c.Args()) > 1 {
					cli.ShowCommandHelp(c, "init")
					os.Exit(1)
				}

				var githubOrg = githubOrg
				var templateName = "default"
				if len(c.Args()) == 1 {
					templateArg := c.Args().Get(0)
					i := strings.Index(templateArg, ":")
					if i != -1 {
						githubOrg = templateArg[:i]
						templateName = templateArg[i+1:]
					}else {
						templateName = templateArg
					}
				}
				workDir, err := os.Getwd()
				if err != nil {
					return err
				}

				return Init(workDir, githubOrg, templateName, osshim, httpshim, io)
			},
		},
	}
}
