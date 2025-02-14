package ecs

import (
	"context"
	"fmt"
	"github/prastamaha/awsctl/utils"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/thediveo/klo"
	"github.com/urfave/cli/v3"
)

func (e *ECS) GetServicesCLI() *cli.Command {
	return &cli.Command{
		Name:    "ecs-service",
		Aliases: ecsServiceAliases,
		Usage:   "Get ECS services",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "cluster",
				Required: false,
				Usage:    "Cluster name",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.String("cluster") == "" {
				e.GetServicesCommandFzf()
			} else {
				e.GetServicesCommand(cmd.String("cluster"))
			}
			return nil
		},
	}
}

func (e *ECS) GetServicesCommand(cluster string) {
	outputs := e.GetAllServices(cluster)
	if len(outputs) == 0 {
		fmt.Printf("No resources found in %s region\n", os.Getenv("AWS_REGION"))
		return
	}

	prn, err := klo.PrinterFromFlag("", &klo.Specs{DefaultColumnSpec: "NAME:{.Name},REPLICAS:{.Replicas},TASK_DEFINITION:{.TaskDefinition},STATUS:{.Status},TYPE:{.Type},CREATED_AT:{.CreatedAt}"})
	if err != nil {
		panic(err)
	}

	table, err := klo.NewSortingPrinter("{.Name}", prn)
	if err != nil {
		panic(err)
	}
	table.Fprint(os.Stdout, outputs)
}

func (e *ECS) GetServicesCommandFzf() {
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
		name := ecsClusters[i].Name
		e.GetServicesCommand(name)
	}
}

func (e *ECS) GetAllServices(cluster string) []ECSServicesList {
	client := ecs.NewFromConfig(e.AWSConfig)

	respGetCluster, err := client.DescribeClusters(context.TODO(), &ecs.DescribeClustersInput{
		Clusters: []string{cluster},
	})
	if err != nil {
		log.Fatalf("unable to describe clusters, %v", err)
	}

	var outputs []ECSServicesList
	var nextToken *string

	for {
		respServiceList, err := client.ListServices(context.TODO(), &ecs.ListServicesInput{
			Cluster:   respGetCluster.Clusters[0].ClusterArn,
			NextToken: nextToken,
		})
		if err != nil {
			log.Fatalf("unable to list services, %v", err)
		}

		if len(respServiceList.ServiceArns) == 0 {
			break
		}

		resp, err := client.DescribeServices(context.TODO(), &ecs.DescribeServicesInput{
			Cluster:  respGetCluster.Clusters[0].ClusterArn,
			Services: respServiceList.ServiceArns,
		})
		if err != nil {
			log.Fatalf("unable to describe services, %v", err)
		}

		for _, service := range resp.Services {
			taskDefArnParts := strings.Split(*service.TaskDefinition, ":")
			taskDef := taskDefArnParts[len(taskDefArnParts)-1]

			outputs = append(outputs, ECSServicesList{
				Name:           *service.ServiceName,
				Replicas:       fmt.Sprintf("%d/%d", service.RunningCount, service.DesiredCount),
				Status:         *service.Status,
				Type:           string(service.LaunchType),
				TaskDefinition: taskDef,
				CreatedAt:      service.CreatedAt.String(),
			})
		}

		if respServiceList.NextToken == nil {
			break
		}
		nextToken = respServiceList.NextToken
	}

	return outputs
}
