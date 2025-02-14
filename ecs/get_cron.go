package ecs

import (
	"context"
	"fmt"
	"github/prastamaha/awsctl/utils"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/thediveo/klo"
	"github.com/urfave/cli/v3"
)

func (e *ECS) GetCronsCLI() *cli.Command {
	return &cli.Command{
		Name:    "ecs-cron",
		Aliases: ecsCronliases,
		Usage:   "Get ECS crons",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "cluster",
				Required: false,
				Usage:    "Cluster name",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.String("cluster") == "" {
				e.GetCronCommandFzf()
			} else {
				e.GetCronCommand(cmd.String("cluster"))
			}
			return nil
		},
	}
}

func (e *ECS) GetCronCommand(cluster string) {
	outputs := e.GetECSCron(cluster)
	if len(outputs) == 0 {
		fmt.Printf("No resources found in %s region\n", os.Getenv("AWS_REGION"))
		return
	}

	prn, err := klo.PrinterFromFlag("", &klo.Specs{DefaultColumnSpec: "NAME:{.Name},SCHEDULE:{.Schedule},STATE:{.State}"})
	if err != nil {
		panic(err)
	}

	table, err := klo.NewSortingPrinter("{.Name}", prn)
	if err != nil {
		panic(err)
	}
	table.Fprint(os.Stdout, outputs)
}

func (e *ECS) GetCronCommandFzf() {
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
		e.GetCronCommand(name)
	}
}

func (e *ECS) GetECSCron(cluster string) []ECSCronList {

	ecsClient := ecs.NewFromConfig(e.AWSConfig)

	ecsResponse, err := ecsClient.DescribeClusters(context.TODO(), &ecs.DescribeClustersInput{
		Clusters: []string{cluster},
	})
	if err != nil {
		log.Fatalf("unable to list ecs cluster, %v", err)
	}

	eventbridgeClient := eventbridge.NewFromConfig(e.AWSConfig)

	respRules, err := eventbridgeClient.ListRuleNamesByTarget(context.TODO(), &eventbridge.ListRuleNamesByTargetInput{
		TargetArn: ecsResponse.Clusters[0].ClusterArn,
	})
	if err != nil {
		log.Fatalf("unable to list rules, %v", err)
	}

	var outputs []ECSCronList
	var nextToken *string

	for {
		for _, rule := range respRules.RuleNames {
			descRule, err := eventbridgeClient.DescribeRule(context.TODO(), &eventbridge.DescribeRuleInput{
				Name: &rule,
			})
			if err != nil {
				log.Fatalf("unable to describe rules, %v", err)
			}
			outputs = append(outputs, ECSCronList{
				Name:     *descRule.Name,
				Schedule: *descRule.ScheduleExpression,
				State:    string(descRule.State),
			})
		}

		if respRules.NextToken == nil {
			break
		}
		nextToken = respRules.NextToken
		respRules, err = eventbridgeClient.ListRuleNamesByTarget(context.TODO(), &eventbridge.ListRuleNamesByTargetInput{
			TargetArn: ecsResponse.Clusters[0].ClusterArn,
			NextToken: nextToken,
		})
		if err != nil {
			log.Fatalf("unable to paginate rules, %v", err)
		}
	}

	return outputs
}
