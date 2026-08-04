package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github/prastamaha/awsctl/utils"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func (v *VPC) DescribeCLI() *cli.Command {
	return &cli.Command{
		Name:    "vpc",
		Aliases: vpcAliases,
		Usage:   "Describe a VPC",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Get(0) == "" {
				v.DescribeCommandFzf()
			} else {
				v.DescribeCommand(cmd.Args().Get(0))
			}
			return nil
		},
	}
}

func (v *VPC) DescribeCommand(id string) {
	svc := ec2.NewFromConfig(v.AWSConfig)

	resp, err := svc.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{
		VpcIds: []string{id},
	})
	if err != nil {
		log.Fatalf("failed to describe VPC, %v", err)
	}

	if len(resp.Vpcs) == 0 {
		fmt.Printf("VPC %s not found\n", id)
		return
	}

	vpc := resp.Vpcs[0]

	name := ""
	tags := make(map[string]string)
	if vpc.Tags != nil {
		for _, tag := range vpc.Tags {
			if tag.Key != nil && tag.Value != nil {
				if *tag.Key == "Name" {
					name = *tag.Value
				}
				tags[*tag.Key] = *tag.Value
			}
		}
	}

	var cidrAssociations []CidrBlockAssociation
	if vpc.CidrBlockAssociationSet != nil {
		for _, assoc := range vpc.CidrBlockAssociationSet {
			cidrAssoc := CidrBlockAssociation{}
			if assoc.CidrBlock != nil {
				cidrAssoc.CidrBlock = *assoc.CidrBlock
			}
			if assoc.CidrBlockState != nil {
				cidrAssoc.State = string(assoc.CidrBlockState.State)
			}
			cidrAssociations = append(cidrAssociations, cidrAssoc)
		}
	}

	output := VPCDescribe{
		ID:                    *vpc.VpcId,
		Name:                  name,
		CidrBlock:             *vpc.CidrBlock,
		State:                 string(vpc.State),
		IsDefault:             *vpc.IsDefault,
		InstanceTenancy:       string(vpc.InstanceTenancy),
		DhcpOptionsId:         *vpc.DhcpOptionsId,
		Tags:                  tags,
		CidrBlockAssociations: cidrAssociations,
	}

	yamlData, err := yaml.Marshal(output)
	if err != nil {
		log.Fatalf("failed to marshal VPC data, %v", err)
	}
	fmt.Println(string(yamlData))
}

func (v *VPC) DescribeCommandFzf() {
	allVpcs := v.GetVPCs()
	if len(allVpcs) == 0 {
		fmt.Println("No VPCs found")
		return
	}

	items := make([]string, len(allVpcs))
	for i, vpc := range allVpcs {
		items[i] = fmt.Sprintf("%s %s (%s)", vpc.ID, vpc.Name, vpc.CidrBlock)
	}

	data := utils.FuzzySearch("Select a VPC to describe: ", items)
	for _, i := range data {
		id := allVpcs[i].ID
		v.DescribeCommand(id)
	}
}
