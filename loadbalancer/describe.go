package loadbalancer

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github/prastamaha/awsctl/utils"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func (l *LoadBalancer) DescribeCLI() *cli.Command {
	return &cli.Command{
		Name:    "loadbalancer",
		Aliases: loadbalancerAliases,
		Usage:   "Describe a load balancer",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Get(0) == "" {
				l.DescribeCommandFzf()
			} else {
				l.DescribeCommand(cmd.Args().Get(0))
			}
			return nil
		},
	}
}

func (l *LoadBalancer) DescribeCommand(name string) {
	svc := elasticloadbalancingv2.NewFromConfig(l.AWSConfig)

	resp, err := svc.DescribeLoadBalancers(context.TODO(), &elasticloadbalancingv2.DescribeLoadBalancersInput{
		Names: []string{name},
	})
	if err != nil {
		log.Fatalf("failed to describe load balancer, %v", err)
	}

	if len(resp.LoadBalancers) == 0 {
		fmt.Printf("Error from server (NotFound): load balancer %s not found\n", name)
		return
	}

	lb := resp.LoadBalancers[0]

	state := ""
	if lb.State != nil && lb.State.Code != "" {
		state = string(lb.State.Code)
	}

	availabilityZones := make([]string, 0)
	for _, az := range lb.AvailabilityZones {
		if az.ZoneName != nil {
			availabilityZones = append(availabilityZones, *az.ZoneName)
		}
	}

	customerOwnedIpv4Pool := ""
	if lb.CustomerOwnedIpv4Pool != nil {
		customerOwnedIpv4Pool = *lb.CustomerOwnedIpv4Pool
	}

	output := LoadBalancerDescribe{
		Name:              *lb.LoadBalancerName,
		ARN:               *lb.LoadBalancerArn,
		DNSName:           *lb.DNSName,
		Type:              string(lb.Type),
		Scheme:            string(lb.Scheme),
		State:             state,
		VpcId:             *lb.VpcId,
		AvailabilityZones: availabilityZones,
		SecurityGroups:    lb.SecurityGroups,
		IpAddressType:     string(lb.IpAddressType),
		CustomerOwnedIpv4Pool: customerOwnedIpv4Pool,
		Tags:              l.getTags(*lb.LoadBalancerArn),
		Listeners:         l.getListeners(*lb.LoadBalancerArn),
		TargetGroups:      l.getTargetGroups(*lb.LoadBalancerArn),
	}

	yamlData, err := yaml.Marshal(output)
	if err != nil {
		fmt.Println("Error marshalling to YAML:", err)
		return
	}
	fmt.Println(string(yamlData))
}

func (l *LoadBalancer) getTags(arn string) map[string]string {
	svc := elasticloadbalancingv2.NewFromConfig(l.AWSConfig)

	resp, err := svc.DescribeTags(context.TODO(), &elasticloadbalancingv2.DescribeTagsInput{
		ResourceArns: []string{arn},
	})
	if err != nil {
		return nil
	}

	tags := make(map[string]string)
	for _, desc := range resp.TagDescriptions {
		for _, tag := range desc.Tags {
			tags[*tag.Key] = *tag.Value
		}
	}

	return tags
}

func (l *LoadBalancer) getListeners(arn string) []ListenerInfo {
	svc := elasticloadbalancingv2.NewFromConfig(l.AWSConfig)

	resp, err := svc.DescribeListeners(context.TODO(), &elasticloadbalancingv2.DescribeListenersInput{
		LoadBalancerArn: &arn,
	})
	if err != nil {
		return nil
	}

	listeners := make([]ListenerInfo, 0)
	for _, listener := range resp.Listeners {
		listeners = append(listeners, ListenerInfo{
			Port:     *listener.Port,
			Protocol: string(listener.Protocol),
			ARN:      *listener.ListenerArn,
		})
	}

	return listeners
}

func (l *LoadBalancer) getTargetGroups(lbArn string) []TargetGroupInfo {
	svc := elasticloadbalancingv2.NewFromConfig(l.AWSConfig)

	resp, err := svc.DescribeTargetGroups(context.TODO(), &elasticloadbalancingv2.DescribeTargetGroupsInput{
		LoadBalancerArn: &lbArn,
	})
	if err != nil {
		return nil
	}

	targetGroups := make([]TargetGroupInfo, 0)
	for _, tg := range resp.TargetGroups {
		targetGroups = append(targetGroups, TargetGroupInfo{
			Name:       *tg.TargetGroupName,
			ARN:        *tg.TargetGroupArn,
			Port:       *tg.Port,
			Protocol:   string(tg.Protocol),
			TargetType: string(tg.TargetType),
		})
	}

	return targetGroups
}

func (l *LoadBalancer) DescribeCommandFzf() {
	allLBs := l.GetLoadBalancers()
	if len(allLBs) == 0 {
		fmt.Println("No load balancers found")
		return
	}

	items := make([]string, len(allLBs))
	for i, v := range allLBs {
		items[i] = fmt.Sprintf("%s (%s)", v.Name, v.Type)
	}

	data := utils.FuzzySearch("Select a load balancer to describe: ", items)
	for _, i := range data {
		name := allLBs[i].Name
		l.DescribeCommand(name)
	}
}
