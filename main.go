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

var version = "v0.0.3"

func main() {
	// load aws config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// initialize services
	ec2Ins := instance.Instance{AWSConfig: cfg}
	ec2Ssm := ssm.SSM{AWSConfig: cfg}
	ecs := ecs.ECS{AWSConfig: cfg}

	// cli commands
	app := &cli.Command{
		Name:  "awsctl",
		Usage: "Simplify aws cli",
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Get resources",
				Commands: []*cli.Command{
					ec2Ins.GetCLI(),
					ecs.GetClustersCLI(),
					ecs.GetCronsCLI(),
					ecs.GetServicesCLI(),
				},
			},
			{
				Name:  "describe",
				Usage: "Describe resources",
				Commands: []*cli.Command{
					ec2Ins.DescribeCLI(),
					ecs.DescribeServiceCLI(),
					ecs.DescribeClusterCLI(),
					ecs.DescribeCronCLI(),
				},
			},
			{
				Name:  "stop",
				Usage: "Stop resources",
				Commands: []*cli.Command{
					ec2Ins.StopCLI(),
					ecs.StopServiceCLI(),
				},
			},
			{
				Name:  "restart",
				Usage: "Restart resources",
				Commands: []*cli.Command{
					ec2Ins.RestartCLI(),
					ecs.RestartServiceCLI(),
				},
			},
			{
				Name:  "start",
				Usage: "Start resources",
				Commands: []*cli.Command{
					ec2Ins.StartCLI(),
					ecs.StartServiceCLI(),
					ec2Ssm.StartSSMCLI(),
				},
			},
			ec2Ssm.StartSSMCLI(),
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
