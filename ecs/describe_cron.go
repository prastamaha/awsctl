package ecs

import (
	"context"
	"fmt"
	"github/prastamaha/awsctl/utils"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"gopkg.in/yaml.v3"
)

func (e *ECS) DescribeCronCommand(cronName string) {
	outputs := e.DescribeECSCron(cronName)

	yamlData, err := yaml.Marshal(outputs)
	if err != nil {
		fmt.Println("Error marshalling to YAML:", err)
		return
	}
	fmt.Println(string(yamlData))
}

func (e *ECS) DescribeCronCommandFzf() {
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

		ecsClusters := e.GetECSCron(clusterName)
		if len(ecsClusters) == 0 {
			fmt.Println("No ecs cron found")
			return
		}

		items := make([]string, len(ecsClusters))
		for i, v := range ecsClusters {
			items[i] = v.Name
		}

		data := utils.FuzzySearch("Select an ecs cron: ", items)
		for _, i := range data {
			cronName := ecsClusters[i].Name
			e.DescribeCronCommand(cronName)
		}
	}
}

func (e *ECS) DescribeECSCron(cronName string) ECSCronDescribe {
	eventbridgeClient := eventbridge.NewFromConfig(e.AWSConfig)

	descRule, err := eventbridgeClient.DescribeRule(context.TODO(), &eventbridge.DescribeRuleInput{
		Name: &cronName,
	})
	if err != nil {
		log.Fatalf("unable to describe rules, %v", err)
	}

	roleArn := ""
	if descRule.RoleArn != nil {
		roleArn = *descRule.RoleArn
	}

	outputs := ECSCronDescribe{
		Name:        *descRule.Name,
		Schedule:    *descRule.ScheduleExpression,
		State:       string(descRule.State),
		Arn:         *descRule.Arn,
		RoleArn:     roleArn,
		Description: *descRule.Description,
	}

	return outputs
}
