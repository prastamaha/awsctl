package subnet

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/thediveo/klo"
	"github.com/urfave/cli/v3"
)

func (s *Subnet) GetCLI() *cli.Command {
	return &cli.Command{
		Name:    "subnet",
		Aliases: subnetAliases,
		Usage:   "Get subnets",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			s.GetCommand()
			return nil
		},
	}
}

func (s *Subnet) GetCommand() {
	outputs := s.GetSubnets()
	if len(outputs) == 0 {
		fmt.Printf("No resources found in %s region\n", os.Getenv("AWS_REGION"))
		return
	}

	prn, err := klo.PrinterFromFlag("", &klo.Specs{DefaultColumnSpec: "ID:{.ID},NAME:{.Name},VPC_ID:{.VpcId},CIDR_BLOCK:{.CidrBlock},AVAILABILITY_ZONE:{.AvailabilityZone},AVAILABLE_IP:{.AvailableIpAddressCount}"})
	if err != nil {
		panic(err)
	}

	table, err := klo.NewSortingPrinter("{.Name}", prn)
	if err != nil {
		panic(err)
	}
	table.Fprint(os.Stdout, outputs)
}

func (s *Subnet) GetSubnets() []SubnetList {
	svc := ec2.NewFromConfig(s.AWSConfig)

	resp, err := svc.DescribeSubnets(context.TODO(), &ec2.DescribeSubnetsInput{})
	if err != nil {
		log.Fatalf("failed to list subnets, %v", err)
	}

	var outputs []SubnetList
	for _, subnet := range resp.Subnets {
		name := ""
		if subnet.Tags != nil {
			for _, tag := range subnet.Tags {
				if tag.Key != nil && tag.Value != nil && *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}
		}

		output := SubnetList{}
		if subnet.SubnetId != nil {
			output.ID = *subnet.SubnetId
		}
		output.Name = name
		if subnet.VpcId != nil {
			output.VpcId = *subnet.VpcId
		}
		if subnet.CidrBlock != nil {
			output.CidrBlock = *subnet.CidrBlock
		}
		if subnet.AvailabilityZone != nil {
			output.AvailabilityZone = *subnet.AvailabilityZone
		}
		if subnet.AvailableIpAddressCount != nil {
			output.AvailableIpAddressCount = *subnet.AvailableIpAddressCount
		}

		outputs = append(outputs, output)
	}

	return outputs
}
