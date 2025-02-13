package ecs

import (
	"context"
	"fmt"
	"github/prastamaha/awsctl/utils"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"gopkg.in/yaml.v3"
)

func convertDeployments(deployments []types.Deployment) []ECSServicesDeploymentDescribe {
	var result []ECSServicesDeploymentDescribe
	for _, d := range deployments {
		result = append(result, ECSServicesDeploymentDescribe{
			Id:             *d.Id,
			Status:         *d.Status,
			DesiredCount:   d.DesiredCount,
			FailedTasks:    d.FailedTasks,
			PendingCount:   d.PendingCount,
			RunningCount:   d.RunningCount,
			CreatedAt:      d.CreatedAt.String(),
			UpdatedAt:      d.UpdatedAt.String(),
			RolloutState:   string(d.RolloutState),
			TaskDefinition: *d.TaskDefinition,
		})
	}
	return result
}

func convertLoadBalancers(loadBalancers []types.LoadBalancer) []ECSServicesLoadbalancerDescribe {
	var result []ECSServicesLoadbalancerDescribe
	for _, lb := range loadBalancers {
		LoadbalancerName := ""
		if lb.ContainerName == nil {
			LoadbalancerName = *lb.LoadBalancerName
		}
		result = append(result, ECSServicesLoadbalancerDescribe{
			LoadbalancerName: LoadbalancerName,
			TargetGroupArn:   *lb.TargetGroupArn,
			ContainerName:    *lb.ContainerName,
			ContainerPort:    *lb.ContainerPort,
		})
	}
	return result
}

func convertEvents(events []types.ServiceEvent) []ECSServicesEventDescribe {
	var result []ECSServicesEventDescribe
	for _, e := range events {
		result = append(result, ECSServicesEventDescribe{
			CreatedAt: e.CreatedAt.String(),
			Message:   *e.Message,
		})
	}
	return result
}

func (e *ECS) DescribeServiceCommand(service string, cluster string) {
	client := ecs.NewFromConfig(e.AWSConfig)

	respGetCluster, err := client.DescribeClusters(context.TODO(), &ecs.DescribeClustersInput{
		Clusters: []string{cluster},
	})
	if err != nil {
		log.Fatalf("unable to describe clusters, %v", err)
	}

	resp, err := client.DescribeServices(context.TODO(), &ecs.DescribeServicesInput{
		Cluster:  respGetCluster.Clusters[0].ClusterArn,
		Services: []string{service},
	})
	if err != nil {
		log.Fatalf("unable to describe services, %v", err)
	}

	if len(resp.Services) == 0 {
		fmt.Printf("Error from server (NotFound): ecs service %s not found\n", service)
		return
	}

	svc := resp.Services[0]
	outputs := ECSServiceDescribe{
		Name:                 *svc.ServiceName,
		Arn:                  *svc.ServiceArn,
		Events:               convertEvents(svc.Events),
		CreatedAt:            svc.CreatedAt.String(),
		DesiredCount:         svc.DesiredCount,
		RunningCount:         svc.RunningCount,
		PendingCount:         svc.PendingCount,
		Status:               *svc.Status,
		EnableExecuteCommand: svc.EnableExecuteCommand,
		TaskDefinition:       *svc.TaskDefinition,
		Deployments:          convertDeployments(svc.Deployments),
		LoadBalancers:        convertLoadBalancers(svc.LoadBalancers),
		SecurityGroups:       svc.NetworkConfiguration.AwsvpcConfiguration.SecurityGroups,
		Subnets:              svc.NetworkConfiguration.AwsvpcConfiguration.Subnets,
	}

	yamlData, err := yaml.Marshal(outputs)
	if err != nil {
		fmt.Println("Error marshalling to YAML:", err)
		return
	}
	fmt.Println(string(yamlData))
}

func (e *ECS) DescribeServiceCommandFzf() {
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
		if len(ecsServices) == 0 {
			fmt.Println("No ecs service found")
			return
		}

		items := make([]string, len(ecsServices))
		for i, v := range ecsServices {
			items[i] = v.Name
		}

		data := utils.FuzzySearch("Select an ecs service to describe:: ", items)
		for _, i := range data {
			serviceName := ecsServices[i].Name
			e.DescribeServiceCommand(serviceName, clusterName)
		}
	}
}
