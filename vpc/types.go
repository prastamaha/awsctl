package vpc

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

var vpcAliases = []string{"vpcs"}

type VPC struct {
	AWSConfig aws.Config
}

type VPCList struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	CidrBlock   string `json:"cidrBlock" yaml:"cidrBlock"`
	State       string `json:"state" yaml:"state"`
	IsDefault   bool   `json:"isDefault" yaml:"isDefault"`
}

type VPCDescribe struct {
	ID                  string            `json:"id" yaml:"id"`
	Name                string            `json:"name" yaml:"name"`
	CidrBlock           string            `json:"cidrBlock" yaml:"cidrBlock"`
	State               string            `json:"state" yaml:"state"`
	IsDefault           bool              `json:"isDefault" yaml:"isDefault"`
	InstanceTenancy     string            `json:"instanceTenancy" yaml:"instanceTenancy"`
	DhcpOptionsId       string            `json:"dhcpOptionsId" yaml:"dhcpOptionsId"`
	Tags                map[string]string `json:"tags" yaml:"tags"`
	CidrBlockAssociations []CidrBlockAssociation `json:"cidrBlockAssociations" yaml:"cidrBlockAssociations"`
}

type CidrBlockAssociation struct {
	CidrBlock string `json:"cidrBlock" yaml:"cidrBlock"`
	State     string `json:"state" yaml:"state"`
}
