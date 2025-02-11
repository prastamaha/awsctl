package ecs

import (
	"context"
	"fmt"
	"github/prastamaha/awsctl/utils"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"gopkg.in/yaml.v3"
)

func (e *ECS) DescribeClusterCommand(name string) {
	client := ecs.NewFromConfig(e.AWSConfig)

	resp, err := client.DescribeClusters(context.TODO(), &ecs.DescribeClustersInput{
		Clusters: []string{name},
	})
	if err != nil {
		log.Fatalf("unable to describe clusters, %v", err)
	}

	if len(resp.Clusters) == 0 {
		fmt.Printf("Error from server (NotFound): ecs cluster %s not found\n", name)
		return
	}

	cluster := resp.Clusters[0]
	outputs := ECSClusterDescribe{
		Name:           *cluster.ClusterName,
		Arn:            *cluster.ClusterArn,
		ActiveServices: cluster.ActiveServicesCount,
		RunningTasks:   cluster.RunningTasksCount,
		PendingTasks:   cluster.PendingTasksCount,
		Status:         *cluster.Status,
		Tags:           utils.ConvertECSTags(cluster.Tags),
	}

	yamlData, err := yaml.Marshal(outputs)
	if err != nil {
		fmt.Println("Error marshalling to YAML:", err)
		return
	}
	fmt.Println(string(yamlData))

}
