package ecs

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/thediveo/klo"
)

func (e *ECS) GetServicesCommand(cluster string) {
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

	if len(outputs) == 0 {
		fmt.Printf("No resources found in %s region\n", os.Getenv("AWS_REGION"))
		return
	}

	prn, err := klo.PrinterFromFlag("", &klo.Specs{DefaultColumnSpec: "NAME:{.Name},REPLICAS:{.Replicas},TASK_DEFINITION:{.TaskDefinition},STATUS:{.Status},TYPE:{.Type},CREATED_AT:{.CreatedAt}"})
	if err != nil {
		panic(err)
	}

	// Use a table sorter and tell it to sort by the Name field of our column objects.
	table, err := klo.NewSortingPrinter("{.Name}", prn)
	if err != nil {
		panic(err)
	}
	table.Fprint(os.Stdout, outputs)
}
