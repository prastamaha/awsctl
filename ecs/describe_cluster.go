package ecs

import (
	"context"
	"fmt"
	"github/prastamaha/awsctl/utils"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func (e *ECS) DescribeClusterCLI() *cli.Command {
	return &cli.Command{
		Name:    "ecs-cluster",
		Aliases: ecsClusterAliases,
		Usage:   "Decribe an ECS cluster",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Get(0) == "" {
				e.DescribeClusterCommandFzf()
			} else {
				e.DescribeClusterCommand(cmd.Args().Get(0))
			}
			return nil
		},
	}
}

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

func (e *ECS) DescribeClusterCommandFzf() {
	ecsClusters := e.GetAllECSCluster()
	if len(ecsClusters) == 0 {
		fmt.Println("No ecs clusters found")
		return
	}

	items := make([]string, len(ecsClusters))
	for i, v := range ecsClusters {
		items[i] = v.Name
	}

	data := utils.FuzzySearch("Select an ecs cluster to describe: ", items)
	for _, i := range data {
		name := ecsClusters[i].Name
		e.DescribeClusterCommand(name)
	}
}
