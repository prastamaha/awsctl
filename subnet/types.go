package subnet

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

var subnetAliases = []string{"subnets"}

type Subnet struct {
	AWSConfig aws.Config
}

type SubnetList struct {
	ID                    string `json:"id" yaml:"id"`
	Name                  string `json:"name" yaml:"name"`
	VpcId                 string `json:"vpcId" yaml:"vpcId"`
	CidrBlock             string `json:"cidrBlock" yaml:"cidrBlock"`
	AvailabilityZone      string `json:"availabilityZone" yaml:"availabilityZone"`
	AvailableIpAddressCount int32  `json:"availableIpAddressCount" yaml:"availableIpAddressCount"`
}

type SubnetDescribe struct {
	ID                      string            `json:"id" yaml:"id"`
	Name                    string            `json:"name" yaml:"name"`
	VpcId                   string            `json:"vpcId" yaml:"vpcId"`
	CidrBlock               string            `json:"cidrBlock" yaml:"cidrBlock"`
	AvailabilityZone        string            `json:"availabilityZone" yaml:"availabilityZone"`
	AvailableIpAddressCount int32             `json:"availableIpAddressCount" yaml:"availableIpAddressCount"`
	State                   string            `json:"state" yaml:"state"`
	AssignPublicIp          bool              `json:"assignPublicIp" yaml:"assignPublicIp"`
	MapPublicIpOnLaunch     bool              `json:"mapPublicIpOnLaunch" yaml:"mapPublicIpOnLaunch"`
	DefaultForAz            bool              `json:"defaultForAz" yaml:"defaultForAz"`
	Tags                    map[string]string `json:"tags" yaml:"tags"`
}
