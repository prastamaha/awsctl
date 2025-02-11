package ecs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/thediveo/klo"
)

func (e *ECS) GetClusterCommand() {
	client := ecs.NewFromConfig(e.AWSConfig)

	respListClusters, err := client.ListClusters(context.TODO(), &ecs.ListClustersInput{})
	if err != nil {
		log.Fatalf("unable to list clusters, %v", err)
	}

	resp, err := client.DescribeClusters(context.TODO(), &ecs.DescribeClustersInput{
		Clusters: respListClusters.ClusterArns,
	})
	if err != nil {
		log.Fatalf("unable to list clusters, %v", err)
	}

	var outputs []ECSClusterList
	for _, cluster := range resp.Clusters {
		outputs = append(outputs, ECSClusterList{
			Name:          *cluster.ClusterName,
			ActiveServices: cluster.ActiveServicesCount,
			RunningTasks:  cluster.RunningTasksCount,
			PendingTasks:  cluster.PendingTasksCount,
			Status:        *cluster.Status,
		})
	}

	if len(outputs) == 0 {
		fmt.Printf("No resources found in %s region\n", os.Getenv("AWS_REGION"))
		return
	}

	prn, err := klo.PrinterFromFlag("", &klo.Specs{DefaultColumnSpec: "NAME:{.Name},ACTIVE_SERVICES:{.ActiveServices},RUNNING_TASKS:{.RunningTasks},PENDING_TASKS:{.PendingTasks},STATUS:{.Status}"})
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
