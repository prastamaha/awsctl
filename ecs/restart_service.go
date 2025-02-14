package ecs

import (
	"context"
	"fmt"
	"github/prastamaha/awsctl/utils"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/urfave/cli/v3"
)

func (e *ECS) RestartServiceCLI() *cli.Command {
	return &cli.Command{
		Name:    "ecs-service",
		Aliases: ecsServiceAliases,
		Usage:   "Restart an ECS service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "cluster",
				Required: false,
				Usage:    "Cluster name",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.String("cluster") == "" && cmd.Args().Get(0) == "" {
				e.RestartServiceCommandFzf()
			} else {
				e.RestartServiceCommand(cmd.String("cluster"), cmd.Args().Get(0))
			}
			return nil
		},
	}
}

func (e *ECS) RestartServiceCommand(cluster string, serviceName string) {
	client := ecs.NewFromConfig(e.AWSConfig)

	_, err := client.UpdateService(context.TODO(), &ecs.UpdateServiceInput{
		Cluster:            &cluster,
		Service:            &serviceName,
		ForceNewDeployment: true,
	})
	if err != nil {
		log.Fatalf("failed to restart ecs service, %v", err)
	}

	fmt.Printf("Restarting instance %s...\n", serviceName)
}

func (e *ECS) RestartServiceCommandFzf() {
	ecsClusters := e.GetAllECSCluster()
	if len(ecsClusters) == 0 {
		fmt.Println("No ecs clusters found")
		return
	}

	items := make([]string, len(ecsClusters))
	for i, v := range ecsClusters {
		items[i] = v.Name
	}

	data := utils.FuzzySearch("Select an ecs cluster: ", items)
	for _, i := range data {
		clusterName := ecsClusters[i].Name

		ecsServices := e.GetAllServices(clusterName)
		if len(ecsClusters) == 0 {
			fmt.Println("No ecs services found")
			return
		}

		items := make([]string, len(ecsServices))
		for i, v := range ecsServices {
			items[i] = v.Name
		}

		data := utils.FuzzySearch("Select an ecs service to restart: ", items)
		for _, i := range data {
			serviceName := ecsServices[i].Name
			e.RestartServiceCommand(clusterName, serviceName)
		}
	}
}
