package securitygroup

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/thediveo/klo"
	"github.com/urfave/cli/v3"
)

func (s *SecurityGroup) GetCLI() *cli.Command {
	return &cli.Command{
		Name:    "securitygroup",
		Aliases: securityGroupAliases,
		Usage:   "Get security groups",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "vpc",
				Aliases: []string{"v"},
				Usage:   "Filter by VPC ID",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			vpcId := cmd.String("vpc")
			if vpcId == "" {
				s.GetCommand()
			} else {
				s.GetCommandVpc(vpcId)
			}
			return nil
		},
	}
}

func (s *SecurityGroup) GetCommand() {
	outputs := s.GetSecurityGroups()
	if len(outputs) == 0 {
		fmt.Printf("No resources found in %s region\n", os.Getenv("AWS_REGION"))
		return
	}

	prn, err := klo.PrinterFromFlag("", &klo.Specs{DefaultColumnSpec: "ID:{.ID},NAME:{.Name},VPC_ID:{.VpcId},DESCRIPTION:{.Description}"})
	if err != nil {
		panic(err)
	}

	table, err := klo.NewSortingPrinter("{.Name}", prn)
	if err != nil {
		panic(err)
	}
	table.Fprint(os.Stdout, outputs)
}

func (s *SecurityGroup) GetCommandVpc(vpcId string) {
	outputs := s.GetSecurityGroupsByVpc(vpcId)
	if len(outputs) == 0 {
		fmt.Printf("No security groups found in VPC %s\n", vpcId)
		return
	}

	prn, err := klo.PrinterFromFlag("", &klo.Specs{DefaultColumnSpec: "ID:{.ID},NAME:{.Name},VPC_ID:{.VpcId},DESCRIPTION:{.Description}"})
	if err != nil {
		panic(err)
	}

	table, err := klo.NewSortingPrinter("{.Name}", prn)
	if err != nil {
		panic(err)
	}
	table.Fprint(os.Stdout, outputs)
}

func (s *SecurityGroup) GetSecurityGroups() []SecurityGroupList {
	svc := ec2.NewFromConfig(s.AWSConfig)

	resp, err := svc.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		log.Fatalf("failed to list security groups, %v", err)
	}

	var outputs []SecurityGroupList
	for _, sg := range resp.SecurityGroups {
		name := ""
		if sg.Tags != nil {
			for _, tag := range sg.Tags {
				if tag.Key != nil && tag.Value != nil && *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}
		}

		output := SecurityGroupList{}
		if sg.GroupId != nil {
			output.ID = *sg.GroupId
		}
		output.Name = name
		if sg.Description != nil {
			output.Description = *sg.Description
		}
		if sg.VpcId != nil {
			output.VpcId = *sg.VpcId
		}

		outputs = append(outputs, output)
	}

	return outputs
}

func (s *SecurityGroup) GetSecurityGroupsByVpc(vpcId string) []SecurityGroupList {
	svc := ec2.NewFromConfig(s.AWSConfig)

	resp, err := svc.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcId},
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to list security groups, %v", err)
	}

	var outputs []SecurityGroupList
	for _, sg := range resp.SecurityGroups {
		name := ""
		if sg.Tags != nil {
			for _, tag := range sg.Tags {
				if tag.Key != nil && tag.Value != nil && *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}
		}

		output := SecurityGroupList{}
		if sg.GroupId != nil {
			output.ID = *sg.GroupId
		}
		output.Name = name
		if sg.Description != nil {
			output.Description = *sg.Description
		}
		if sg.VpcId != nil {
			output.VpcId = *sg.VpcId
		}

		outputs = append(outputs, output)
	}

	return outputs
}
