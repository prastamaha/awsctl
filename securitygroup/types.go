package securitygroup

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

var securityGroupAliases = []string{"sg", "sgs", "securitygroups"}

type SecurityGroup struct {
	AWSConfig aws.Config
}

type SecurityGroupList struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	VpcId       string `json:"vpcId" yaml:"vpcId"`
}

type SecurityGroupDescribe struct {
	ID            string            `json:"id" yaml:"id"`
	Name          string            `json:"name" yaml:"name"`
	Description   string            `json:"description" yaml:"description"`
	VpcId         string            `json:"vpcId" yaml:"vpcId"`
	Tags          map[string]string `json:"tags" yaml:"tags"`
	InboundRules  []SecurityRule    `json:"inboundRules" yaml:"inboundRules"`
	OutboundRules []SecurityRule    `json:"outboundRules" yaml:"outboundRules"`
}

type SecurityRule struct {
	Protocol       string              `json:"protocol" yaml:"protocol"`
	FromPort       int32               `json:"fromPort" yaml:"fromPort"`
	ToPort         int32               `json:"toPort" yaml:"toPort"`
	CidrBlocks     []string            `json:"cidrBlocks" yaml:"cidrBlocks"`
	SecurityGroups []SecurityGroupRef  `json:"securityGroups" yaml:"securityGroups"`
	PrefixLists    []PrefixListRef     `json:"prefixLists" yaml:"prefixLists"`
}

type SecurityGroupRef struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

type PrefixListRef struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}
