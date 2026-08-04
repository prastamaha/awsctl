package loadbalancer

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/thediveo/klo"
	"github.com/urfave/cli/v3"
)

func (l *LoadBalancer) GetCLI() *cli.Command {
	return &cli.Command{
		Name:    "loadbalancer",
		Aliases: loadbalancerAliases,
		Usage:   "Get load balancers",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			l.GetCommand()
			return nil
		},
	}
}

func (l *LoadBalancer) GetCommand() {
	outputs := l.GetLoadBalancers()
	if len(outputs) == 0 {
		fmt.Printf("No resources found in %s region\n", os.Getenv("AWS_REGION"))
		return
	}

	prn, err := klo.PrinterFromFlag("", &klo.Specs{DefaultColumnSpec: "NAME:{.Name},TYPE:{.Type},SCHEME:{.Scheme},STATE:{.State},DNS_NAME:{.DNSName},VPC_ID:{.VpcId},CREATE_TIME:{.CreateTime}"})
	if err != nil {
		panic(err)
	}

	table, err := klo.NewSortingPrinter("{.Name}", prn)
	if err != nil {
		panic(err)
	}
	table.Fprint(os.Stdout, outputs)
}

func (l *LoadBalancer) GetLoadBalancers() []LoadBalancerList {
	svc := elasticloadbalancingv2.NewFromConfig(l.AWSConfig)

	resp, err := svc.DescribeLoadBalancers(context.TODO(), &elasticloadbalancingv2.DescribeLoadBalancersInput{})
	if err != nil {
		log.Fatalf("failed to list load balancers, %v", err)
	}

	var outputs []LoadBalancerList
	for _, lb := range resp.LoadBalancers {
		state := ""
		if lb.State != nil && lb.State.Code != "" {
			state = string(lb.State.Code)
		}

		createTime := ""
		if lb.CreatedTime != nil {
			createTime = lb.CreatedTime.String()
		}

		outputs = append(outputs, LoadBalancerList{
			Name:       *lb.LoadBalancerName,
			Type:       string(lb.Type),
			Scheme:     string(lb.Scheme),
			State:      state,
			DNSName:    *lb.DNSName,
			VpcId:      *lb.VpcId,
			CreateTime: createTime,
		})
	}

	return outputs
}
