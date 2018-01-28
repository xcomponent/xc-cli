package commands

import (
	"github.com/urfave/cli"
	"github.com/xcomponent/xc-cli/services"
	"strings"
	"github.com/xcomponent/xc-cli/util"
)

const (
	installConfigUrl = "https://raw.githubusercontent.com/xcomponent/xc-cli/install-config-v1/install-config.json"
	projectName      = ""

	projectTemplatesGithubOrg = "xcomponent-templates"
	addTemplatesGithubOrg     = "xcomponent-add-templates"
)

func GetCommands(workDir string, os services.OsService, http services.HttpService, io services.IoService, zip services.ZipService,
	exec services.ExecService) []cli.Command {
	return []cli.Command{
		GetCliCommand(os, http, io, exec),
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

				var githubOrg = projectTemplatesGithubOrg
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

				command := NewInitCommand(http, io, os, zip)

				return command.Init(workDir, githubOrg, templateName, c.String("project-name"))
			},
		},
		{
			Name:      "add",
			Usage:     "Add new element to the project",
			ArgsUsage: "ELEMENT-NAME",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "element-type",
					Usage: "Element type.",
					Value: "component"},
			},
			Action: func(c *cli.Context) error {
				if len(c.Args()) != 1 {
					cli.ShowCommandHelp(c, "add")

					os.Exit(1)
					return nil
				}

				command := NewAddCommand(os, http, io, zip)

				return command.Execute(workDir, c.String("element-type"), c.Args().Get(0),
					addTemplatesGithubOrg)
			},
		},
		{
			Name:      "rename-project",
			Usage:     "Rename the current project name",
			ArgsUsage: " NEW-PROJECT-NAME",
			Action: func(c *cli.Context) error {
				if len(c.Args()) != 1 {
					cli.ShowCommandHelp(c, "rename-project")

					os.Exit(1)
					return nil
				}

				command := NewRenameProjectCommand(io, os, util.NewFileUtils(os, io))

				return command.Execute(
					workDir, c.Args().Get(0));
			},
		},
	}
}
