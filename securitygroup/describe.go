package securitygroup

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github/prastamaha/awsctl/utils"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func (s *SecurityGroup) DescribeCLI() *cli.Command {
	return &cli.Command{
		Name:    "securitygroup",
		Aliases: securityGroupAliases,
		Usage:   "Describe a security group",
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

func (s *SecurityGroup) DescribeCommand(id string) {
	svc := ec2.NewFromConfig(s.AWSConfig)

	resp, err := svc.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{id},
	})
	if err != nil {
		log.Fatalf("failed to describe security group, %v", err)
	}

	if len(resp.SecurityGroups) == 0 {
		fmt.Printf("Security group %s not found\n", id)
		return
	}

	sg := resp.SecurityGroups[0]

	// Collect all referenced security group IDs and prefix list IDs
	allSgIds := make(map[string]bool)
	allPlIds := make(map[string]bool)

	for _, perm := range sg.IpPermissions {
		if perm.UserIdGroupPairs != nil {
			for _, sgPair := range perm.UserIdGroupPairs {
				if sgPair.GroupId != nil {
					allSgIds[*sgPair.GroupId] = true
				}
			}
		}
		if perm.PrefixListIds != nil {
			for _, pl := range perm.PrefixListIds {
				if pl.PrefixListId != nil {
					allPlIds[*pl.PrefixListId] = true
				}
			}
		}
	}

	for _, perm := range sg.IpPermissionsEgress {
		if perm.UserIdGroupPairs != nil {
			for _, sgPair := range perm.UserIdGroupPairs {
				if sgPair.GroupId != nil {
					allSgIds[*sgPair.GroupId] = true
				}
			}
		}
		if perm.PrefixListIds != nil {
			for _, pl := range perm.PrefixListIds {
				if pl.PrefixListId != nil {
					allPlIds[*pl.PrefixListId] = true
				}
			}
		}
	}

	// Fetch security group names
	sgNameMap := make(map[string]string)
	if len(allSgIds) > 0 {
		sgIds := make([]string, 0, len(allSgIds))
		for id := range allSgIds {
			sgIds = append(sgIds, id)
		}
		
		sgResp, err := svc.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
			GroupIds: sgIds,
		})
		if err == nil {
			for _, refSg := range sgResp.SecurityGroups {
				name := ""
				if refSg.Tags != nil {
					for _, tag := range refSg.Tags {
						if tag.Key != nil && tag.Value != nil && *tag.Key == "Name" {
							name = *tag.Value
							break
						}
					}
				}
				if refSg.GroupId != nil {
					sgNameMap[*refSg.GroupId] = name
				}
			}
		}
	}

	// Fetch prefix list names
	plNameMap := make(map[string]string)
	if len(allPlIds) > 0 {
		plIds := make([]string, 0, len(allPlIds))
		for id := range allPlIds {
			plIds = append(plIds, id)
		}
		
		plResp, err := svc.DescribeManagedPrefixLists(context.TODO(), &ec2.DescribeManagedPrefixListsInput{
			PrefixListIds: plIds,
		})
		if err == nil {
			for _, pl := range plResp.PrefixLists {
				if pl.PrefixListId != nil && pl.PrefixListName != nil {
					plNameMap[*pl.PrefixListId] = *pl.PrefixListName
				}
			}
		}
	}

	// Build output
	name := ""
	tags := make(map[string]string)
	if sg.Tags != nil {
		for _, tag := range sg.Tags {
			if tag.Key != nil && tag.Value != nil {
				if *tag.Key == "Name" {
					name = *tag.Value
				}
				tags[*tag.Key] = *tag.Value
			}
		}
	}

	inboundRules := make([]SecurityRule, 0)
	for _, perm := range sg.IpPermissions {
		rule := SecurityRule{}
		
		if perm.IpProtocol != nil {
			rule.Protocol = *perm.IpProtocol
		}
		if perm.FromPort != nil {
			rule.FromPort = *perm.FromPort
		}
		if perm.ToPort != nil {
			rule.ToPort = *perm.ToPort
		}

		cidrBlocks := make([]string, 0)
		if perm.IpRanges != nil {
			for _, cidr := range perm.IpRanges {
				if cidr.CidrIp != nil {
					cidrBlocks = append(cidrBlocks, *cidr.CidrIp)
				}
			}
		}
		if perm.Ipv6Ranges != nil {
			for _, cidr := range perm.Ipv6Ranges {
				if cidr.CidrIpv6 != nil {
					cidrBlocks = append(cidrBlocks, *cidr.CidrIpv6)
				}
			}
		}
		rule.CidrBlocks = cidrBlocks

		securityGroups := make([]SecurityGroupRef, 0)
		if perm.UserIdGroupPairs != nil {
			for _, sgPair := range perm.UserIdGroupPairs {
				if sgPair.GroupId != nil {
					ref := SecurityGroupRef{
						ID:   *sgPair.GroupId,
						Name: sgNameMap[*sgPair.GroupId],
					}
					securityGroups = append(securityGroups, ref)
				}
			}
		}
		rule.SecurityGroups = securityGroups

		prefixLists := make([]PrefixListRef, 0)
		if perm.PrefixListIds != nil {
			for _, pl := range perm.PrefixListIds {
				if pl.PrefixListId != nil {
					ref := PrefixListRef{
						ID:   *pl.PrefixListId,
						Name: plNameMap[*pl.PrefixListId],
					}
					prefixLists = append(prefixLists, ref)
				}
			}
		}
		rule.PrefixLists = prefixLists

		inboundRules = append(inboundRules, rule)
	}

	outboundRules := make([]SecurityRule, 0)
	for _, perm := range sg.IpPermissionsEgress {
		rule := SecurityRule{}
		
		if perm.IpProtocol != nil {
			rule.Protocol = *perm.IpProtocol
		}
		if perm.FromPort != nil {
			rule.FromPort = *perm.FromPort
		}
		if perm.ToPort != nil {
			rule.ToPort = *perm.ToPort
		}

		cidrBlocks := make([]string, 0)
		if perm.IpRanges != nil {
			for _, cidr := range perm.IpRanges {
				if cidr.CidrIp != nil {
					cidrBlocks = append(cidrBlocks, *cidr.CidrIp)
				}
			}
		}
		if perm.Ipv6Ranges != nil {
			for _, cidr := range perm.Ipv6Ranges {
				if cidr.CidrIpv6 != nil {
					cidrBlocks = append(cidrBlocks, *cidr.CidrIpv6)
				}
			}
		}
		rule.CidrBlocks = cidrBlocks

		securityGroups := make([]SecurityGroupRef, 0)
		if perm.UserIdGroupPairs != nil {
			for _, sgPair := range perm.UserIdGroupPairs {
				if sgPair.GroupId != nil {
					ref := SecurityGroupRef{
						ID:   *sgPair.GroupId,
						Name: sgNameMap[*sgPair.GroupId],
					}
					securityGroups = append(securityGroups, ref)
				}
			}
		}
		rule.SecurityGroups = securityGroups

		prefixLists := make([]PrefixListRef, 0)
		if perm.PrefixListIds != nil {
			for _, pl := range perm.PrefixListIds {
				if pl.PrefixListId != nil {
					ref := PrefixListRef{
						ID:   *pl.PrefixListId,
						Name: plNameMap[*pl.PrefixListId],
					}
					prefixLists = append(prefixLists, ref)
				}
			}
		}
		rule.PrefixLists = prefixLists

		outboundRules = append(outboundRules, rule)
	}

	output := SecurityGroupDescribe{
		ID:            *sg.GroupId,
		Name:          name,
		Description:   *sg.Description,
		VpcId:         *sg.VpcId,
		Tags:          tags,
		InboundRules:  inboundRules,
		OutboundRules: outboundRules,
	}

	yamlData, err := yaml.Marshal(output)
	if err != nil {
		fmt.Printf("Error marshalling to YAML: %v\n", err)
		return
	}
	fmt.Println(string(yamlData))
}

func (s *SecurityGroup) DescribeCommandFzf() {
	allSGs := s.GetSecurityGroups()
	if len(allSGs) == 0 {
		fmt.Println("No security groups found")
		return
	}

	items := make([]string, len(allSGs))
	for i, v := range allSGs {
		items[i] = fmt.Sprintf("%s %s (%s)", v.ID, v.Name, v.VpcId)
	}

	data := utils.FuzzySearch("Select a security group to describe: ", items)
	for _, i := range data {
		id := allSGs[i].ID
		s.DescribeCommand(id)
	}
}
