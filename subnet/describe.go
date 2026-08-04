package subnet

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github/prastamaha/awsctl/utils"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func (s *Subnet) DescribeCLI() *cli.Command {
	return &cli.Command{
		Name:    "subnet",
		Aliases: subnetAliases,
		Usage:   "Describe a subnet",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Get(0) == "" {
				s.DescribeCommandFzf()
			} else {
				s.DescribeCommand(cmd.Args().Get(0))
			}
			return nil
		},
	}
}

func (s *Subnet) DescribeCommand(id string) {
	svc := ec2.NewFromConfig(s.AWSConfig)

	resp, err := svc.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{
		SubnetIds: []string{id},
	})
	if err != nil {
		log.Fatalf("failed to describe subnet, %v", err)
	}

	if len(resp.Subnets) == 0 {
		fmt.Printf("Subnet %s not found\n", id)
		return
	}

	subnet := resp.Subnets[0]

	name := ""
	tags := make(map[string]string)
	if subnet.Tags != nil {
		for _, tag := range subnet.Tags {
			if tag.Key != nil && tag.Value != nil {
				if *tag.Key == "Name" {
					name = *tag.Value
				}
				tags[*tag.Key] = *tag.Value
			}
		}
	}

	output := SubnetDescribe{
		ID:                      *subnet.SubnetId,
		Name:                    name,
		VpcId:                   *subnet.VpcId,
		CidrBlock:               *subnet.CidrBlock,
		AvailabilityZone:        *subnet.AvailabilityZone,
		AvailableIpAddressCount: *subnet.AvailableIpAddressCount,
		State:                   string(subnet.State),
		AssignPublicIp:          *subnet.AssignIpv6AddressOnCreation,
		MapPublicIpOnLaunch:     *subnet.MapPublicIpOnLaunch,
		DefaultForAz:            *subnet.DefaultForAz,
		Tags:                    tags,
	}

	yamlData, err := yaml.Marshal(output)
	if err != nil {
		log.Fatalf("failed to marshal subnet data, %v", err)
	}
	fmt.Println(string(yamlData))
}

func (s *Subnet) DescribeCommandFzf() {
	allSubnets := s.GetSubnets()
	if len(allSubnets) == 0 {
		fmt.Println("No subnets found")
		return
	}

	items := make([]string, len(allSubnets))
	for i, subnet := range allSubnets {
		items[i] = fmt.Sprintf("%s %s (%s)", subnet.ID, subnet.Name, subnet.CidrBlock)
	}

	data := utils.FuzzySearch("Select a subnet to describe: ", items)
	for _, i := range data {
		id := allSubnets[i].ID
		s.DescribeCommand(id)
	}
}
