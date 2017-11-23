package commands

import (
	"github.com/urfave/cli"
	"strings"
	"runtime"
	"github.com/xcomponent/xc-cli/services"
)

const (
	installConfigUrl = "https://raw.githubusercontent.com/xcomponent/xc-cli/install-config-v1/install-config.json"
	projectName      = ""

	githubOrg = "xcomponent-templates"
)

func GetCommands(os services.OsService, http services.HttpService, io services.IoService, zip services.ZipService,
	exec services.ExecService) ([]cli.Command) {
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
				var installCommand = InstallCommand{os: os, io: io, http: http, exec: exec}

				return installCommand.Install(
					c.String("install-config-url"),
					runtime.GOOS,
					runtime.GOARCH)
			},
		},
		{
			Name:      "init",
			Usage:     "Initialize a new XComponent project",
			ArgsUsage: "[ORGANIZATION:][TEMPLATE-NAME]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "project-name",
					Usage: "Project name.",
					Value: projectName},
			},
			Action: func(c *cli.Context) error {
				if len(c.Args()) > 1 {
					cli.ShowCommandHelp(c, "init")

					os.Exit(1)
					return nil
				}

				var githubOrg = githubOrg
				var templateName = "default"
				if len(c.Args()) == 1 {
					templateArg := c.Args().Get(0)
					i := strings.Index(templateArg, ":")
					if i != -1 {
						githubOrg = templateArg[:i]
						templateName = templateArg[i+1:]
					} else {
						templateName = templateArg
					}
				}

				workDir, err := os.Getwd()
				if err != nil {
					return err
				}

				var initCommand = InitCommand{http, io, os, zip}

				return initCommand.Init(workDir, githubOrg, templateName, c.String("project-name"))
			},
		},
	}
}
