package main

import (
	"context"
	"fmt"
	"github/prastamaha/awsctl/ecs"
	"github/prastamaha/awsctl/instance"
	"github/prastamaha/awsctl/ssm"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/urfave/cli/v3"
)

var version = "v0.0.2"

func main() {
	// load aws config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// initialize services
	ins := instance.Instance{AWSConfig: cfg}
	ec2Ssm := ssm.SSM{AWSConfig: cfg}
	ecs := ecs.ECS{AWSConfig: cfg}

	// aliases
	instanceAliases := []string{"in", "ins", "instances", "ec2"}
	ssmAliases := []string{"start-session", "ssh"}
	ecsClusterAliases := []string{"ecsc", "ecscluster", "ecsclusters"}
	ecsServiceAliases := []string{"ecss", "ecsservice", "ecsservices", "ecssvc"}

	// cli commands
	app := &cli.Command{
		Name:  "awsctl",
		Usage: "Simplify aws cli",
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Get resources",
				Commands: []*cli.Command{
					{
						Name:    "instance",
						Aliases: instanceAliases,
						Usage:   "Get EC2 Instances",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							ins.GetCommand()
							return nil
						},
					},
					{
						Name:    "ecs-cluster",
						Aliases: ecsClusterAliases,
						Usage:   "Get ECS Clusters",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							ecs.GetClusterCommand()
							return nil
						},
					},
					{
						Name:    "ecs-service",
						Aliases: ecsServiceAliases,
						Usage:   "Get ECS Services",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "cluster",
								Required: false,
								Usage:    "Cluster name",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							if cmd.String("cluster") == "" {
								ecs.GetServicesCommandFzf()
							} else {
								ecs.GetServicesCommand(cmd.String("cluster"))
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "describe",
				Usage: "Describe resources",
				Commands: []*cli.Command{
					{
						Name:      "instance",
						Aliases:   instanceAliases,
						Usage:     "Describe EC2 Instances",
						ArgsUsage: "[instance id]",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							if cmd.Args().Get(0) == "" {
								ins.DescribeCommandFzf()
							} else {
								ins.DescribeCommand(cmd.Args().Get(0))
							}
							return nil
						},
					},
					{
						Name:      "ecs-cluster",
						Aliases:   ecsClusterAliases,
						Usage:     "Describe ECS Clusters",
						ArgsUsage: "[cluster name]",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							if cmd.Args().Get(0) == "" {
								ecs.DescribeClusterCommandFzf()
							} else {
								ecs.DescribeClusterCommand(cmd.Args().Get(0))
							}
							return nil
						},
					},
					{
						Name:    "ecs-service",
						Aliases: ecsServiceAliases,
						Usage:   "Describe ECS Services",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "cluster",
								Required: false,
								Usage:    "Cluster name",
							},
						},
						ArgsUsage: "[service name]",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							if cmd.Args().Get(0) == "" && cmd.String("cluster") == "" {
								ecs.DescribeServiceCommandFzf()
							} else {
								ecs.DescribeServiceCommand(cmd.Args().Get(0), cmd.String("cluster"))
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "stop",
				Usage: "Stop resources",
				Commands: []*cli.Command{
					{
						Name:      "instance",
						Aliases:   instanceAliases,
						Usage:     "Stop EC2 Instance",
						ArgsUsage: "[instance id]",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							if cmd.Args().Get(0) == "" {
								ins.StopCommandFzf()
							} else {
								ins.StopCommand(cmd.Args().Get(0))
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "restart",
				Usage: "Restart resources",
				Commands: []*cli.Command{
					{
						Name:      "instance",
						Aliases:   instanceAliases,
						Usage:     "Restart EC2 Instance",
						ArgsUsage: "[instance id]",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							if cmd.Args().Get(0) == "" {
								ins.RestartCommandFzf()
							} else {
								ins.RestartCommand(cmd.Args().Get(0))
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "start",
				Usage: "Start resources",
				Commands: []*cli.Command{
					{
						Name:      "instance",
						Aliases:   instanceAliases,
						Usage:     "Start EC2 Instance",
						ArgsUsage: "[instance id]",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							if cmd.Args().Get(0) == "" {
								ins.StartCommandFzf()
							} else {
								ins.StartCommand(cmd.Args().Get(0))
							}
							return nil
						},
					},
				},
			},
			{
				Name:    "ssm",
				Usage:   "Start SSM",
				Aliases: ssmAliases,
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Get(0) == "" {
						ec2Ssm.StartSessionFzf()
					} else {
						ec2Ssm.StartSessionTarget(cmd.Args().Get(0))
					}
					return nil
				},
			},
			{
				Name:  "version",
				Usage: "Current Version",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println(version)
					return nil
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
