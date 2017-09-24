package commands

import (
	"github.com/urfave/cli"
	"github.com/daniellavoie/xc-cli/service"
)

const protocol = "https"

func GetCommands(installConfigUrl string, cloudPlatformService service.CloudPlatformService) ([]cli.Command) {
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
		/*{
			Name:      "push",
			Usage:     "Upload an XC component to XC Cloud Platform",
			ArgsUsage: "COMPONENT-NAME[:LABEL] COMPONENT-PATH",
			Action: func(c *cli.Context) error {
				if len(c.Args()) != 2 {
					cli.ShowCommandHelp(c, "run")
					os.Exit(1)
				}

				componentNameAndLabel := strings.Split(c.Args().Get(0), ":")
				if len(componentNameAndLabel) > 2 {
					cli.ShowCommandHelp(c, "run")
					os.Exit(1)
				}

				request, err := BuildPushRequest(c.Args().Get(0), c.Args().Get(2))
				if err != nil {
					return err
				}

				return PushComponent(request, cloudPlatformService)
			},
		},
		{
			Name:      "run",
			Usage:     "Launch a component stack based on docker-compose",
			ArgsUsage: "STACK-NAME DOCKER-COMPOSE-YML-PATH",
			Action: func(c *cli.Context) error {
				if len(c.Args()) != 2 {
					cli.ShowCommandHelp(c, "run")
					os.Exit(1)
				}

				stackName := c.Args().Get(0)
				stackDefinitionPath := c.Args().Get(2)

				request, err := BuildRunStackRequest(stackName, stackDefinitionPath)
				if err != nil {
					return err
				}

				return RunStack(request, cloudPlatformService)
			},
		},
		{
			Name:      "logs",
			Usage:     "Access the logs of a service instance",
			ArgsUsage: "STACK-NAME SERVICE-NAME [INSTANCE-NUMBER]",
			Action: func(c *cli.Context) error {
				argsLen := len(c.Args())
				if argsLen >= 2 && argsLen <= 3 {
					cli.ShowCommandHelp(c, "logs")
					os.Exit(1)
				}

				request := LogsCommandRequest{StackName: c.Args().First(), ServiceName: c.Args().Get(1), Instance: "1"}
				if argsLen == 3 {
					request.Instance = c.Args().Get(2)
				}

				return RetreiveLogs(request, cloudPlatformService)
			},
		},*/
	}
}
