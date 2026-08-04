package vpc

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/thediveo/klo"
	"github.com/urfave/cli/v3"
)

func (v *VPC) GetCLI() *cli.Command {
	return &cli.Command{
		Name:    "vpc",
		Aliases: vpcAliases,
		Usage:   "Get VPCs",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			v.GetCommand()
			return nil
		},
	}
}

func (v *VPC) GetCommand() {
	outputs := v.GetVPCs()
	if len(outputs) == 0 {
		fmt.Printf("No resources found in %s region\n", os.Getenv("AWS_REGION"))
		return
	}

	prn, err := klo.PrinterFromFlag("", &klo.Specs{DefaultColumnSpec: "ID:{.ID},NAME:{.Name},CIDR_BLOCK:{.CidrBlock},STATE:{.State},IS_DEFAULT:{.IsDefault}"})
	if err != nil {
		panic(err)
	}

	table, err := klo.NewSortingPrinter("{.Name}", prn)
	if err != nil {
		panic(err)
	}
	table.Fprint(os.Stdout, outputs)
}

func (v *VPC) GetVPCs() []VPCList {
	svc := ec2.NewFromConfig(v.AWSConfig)

	resp, err := svc.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{})
	if err != nil {
		log.Fatalf("failed to list VPCs, %v", err)
	}

	var outputs []VPCList
	for _, vpc := range resp.Vpcs {
		name := ""
		if vpc.Tags != nil {
			for _, tag := range vpc.Tags {
				if tag.Key != nil && tag.Value != nil && *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}
		}

		output := VPCList{}
		if vpc.VpcId != nil {
			output.ID = *vpc.VpcId
		}
		output.Name = name
		if vpc.CidrBlock != nil {
			output.CidrBlock = *vpc.CidrBlock
		}
		output.State = string(vpc.State)
		if vpc.IsDefault != nil {
			output.IsDefault = *vpc.IsDefault
		}

		outputs = append(outputs, output)
	}

	return outputs
}
