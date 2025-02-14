package ecs

import (
	"context"
	"fmt"
	"github/prastamaha/awsctl/utils"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/applicationautoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/urfave/cli/v3"
)

func (e *ECS) StartServiceCLI() *cli.Command {
	return &cli.Command{
		Name:    "ecs-service",
		Aliases: ecsServiceAliases,
		Usage:   "Start an ECS service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "cluster",
				Required: false,
				Usage:    "Cluster name",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.String("cluster") == "" && cmd.Args().Get(0) == "" {
				e.StartServiceCommandFzf()
			} else {
				e.StartServiceCommand(cmd.String("cluster"), cmd.Args().Get(0))
			}
			return nil
		},
	}
}

func (e *ECS) StartServiceCommand(cluster string, serviceName string) {
	clientAutoScale := applicationautoscaling.NewFromConfig(e.AWSConfig)

	result, err := clientAutoScale.DescribeScalableTargets(context.TODO(), &applicationautoscaling.DescribeScalableTargetsInput{
		ServiceNamespace: "ecs",
		ResourceIds:      []string{fmt.Sprintf("service/%s/%s", cluster, serviceName)},
	})
	if err != nil {
		log.Fatalf("failed to start ecs service, %v", err)
	}

	if len(result.ScalableTargets) == 0 {
		log.Fatalf("failed to start ecs service, no scalable target found")
	}

	minCapacity := *result.ScalableTargets[0].MinCapacity

	clientEcs := ecs.NewFromConfig(e.AWSConfig)
	_, err = clientEcs.UpdateService(context.TODO(), &ecs.UpdateServiceInput{
		Cluster:      &cluster,
		Service:      &serviceName,
		DesiredCount: &minCapacity,
	})
	if err != nil {
		log.Fatalf("failed to start ecs service, %v", err)
	}

	fmt.Printf("Starting instance %s...\n", serviceName)
}

func (e *ECS) StartServiceCommandFzf() {
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
			e.StartServiceCommand(clusterName, serviceName)
		}
	}
}
