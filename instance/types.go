package instance

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

var instanceAliases = []string{"in", "ins", "instances", "ec2"}

type Instance struct {
	AWSConfig aws.Config
}

type EC2InstanceList struct {
	Name       string `json:"name" yaml:"name"`
	ID         string `json:"id" yaml:"id"`
	State      string `json:"state" yaml:"state"`
	Type       string `json:"type" yaml:"type"`
	PrivateIP  string `json:"privateIp" yaml:"privateIp"`
	PublicIP   string `json:"publicIp" yaml:"publicIp"`
	LaunchTime string `json:"launchTime" yaml:"launchTime"`
}

type EC2InstanceDescribe struct {
	Name           string            `json:"name" yaml:"name"`
	IAMRoleArn     string            `json:"iamRoleArn" yaml:"iamRoleArn"`
	ID             string            `json:"id" yaml:"id"`
	State          string            `json:"state" yaml:"state"`
	Type           string            `json:"type" yaml:"type"`
	ImageID        string            `json:"imageId" yaml:"imageId"`
	SubnetID       string            `json:"subnetId" yaml:"subnetId"`
	VPCID          string            `json:"vpcId" yaml:"vpcId"`
	PrivateIP      string            `json:"privateIp" yaml:"privateIp"`
	PublicIP       string            `json:"publicIp" yaml:"publicIp"`
	DNSPublic      string            `json:"dnsPublic" yaml:"dnsPublic"`
	DnsPrivate     string            `json:"dnsPrivate" yaml:"dnsPrivate"`
	SecurityGroups map[string]string `json:"securityGroups" yaml:"securityGroups"`
	Tags           map[string]string `json:"tags" yaml:"tags"`
	LaunchTime     string            `json:"launchTime" yaml:"launchTime"`
}
