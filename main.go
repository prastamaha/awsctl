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

var version = "v0.0.1"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	ins := instance.Instance{
		AWSConfig: cfg,
	}

	ec2Ssm := ssm.SSM{
		AWSConfig: cfg,
	}

	ecs := ecs.ECS{
		AWSConfig: cfg,
	}

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
						Aliases: []string{"in", "ins", "instances", "ec2"},
						Usage:   "Get EC2 Instances",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							ins.GetCommand()
							return nil
						},
					},
					{
						Name:    "ecs-cluster",
						Aliases: []string{"ecsc", "ecscluster", "ecsclusters"},
						Usage:   "Get ECS Clusters",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							ecs.GetClusterCommand()
							return nil
						},
					},
					{
						Name:    "ecs-service",
						Aliases: []string{"ecss", "ecsservice", "ecsservices"},
						Usage:   "Get ECS Services",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "cluster",
								Required: true,
								Usage:    "Cluster name",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							ecs.GetServicesCommand(cmd.String("cluster"))
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
						Aliases:   []string{"in", "ins", "instances", "ec2"},
						Usage:     "Describe EC2 Instances",
						ArgsUsage: "[instance id]",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							ins.DescribeCommand(cmd.Args().Get(0))
							return nil
						},
					},
					{
						Name:      "ecs-cluster",
						Aliases:   []string{"ecsc", "ecscluster", "ecsclusters"},
						Usage:     "Describe ECS Clusters",
						ArgsUsage: "[cluster name]",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							ecs.DescribeClusterCommand(cmd.Args().Get(0))
							return nil
						},
					},
					{
						Name:    "ecs-service",
						Aliases: []string{"ecss", "ecsservice", "ecsservices"},
						Usage:   "Describe ECS Services",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "cluster",
								Required: true,
								Usage:    "Cluster name",
							},
						},
						ArgsUsage: "[service name]",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							ecs.DescribeServiceCommand(cmd.Args().Get(0), cmd.String("cluster"))
							return nil
						},
					},
				},
			},
			{
				Name:  "ssm",
				Usage: "Start SSM",
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
