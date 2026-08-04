package loadbalancer

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

var loadbalancerAliases = []string{"lb", "lbs", "loadbalancers"}

type LoadBalancer struct {
	AWSConfig aws.Config
}

type LoadBalancerList struct {
	Name       string `json:"name" yaml:"name"`
	Type       string `json:"type" yaml:"type"`
	Scheme     string `json:"scheme" yaml:"scheme"`
	State      string `json:"state" yaml:"state"`
	DNSName    string `json:"dnsName" yaml:"dnsName"`
	VpcId      string `json:"vpcId" yaml:"vpcId"`
	CreateTime string `json:"createTime" yaml:"createTime"`
}

type LoadBalancerDescribe struct {
	Name              string            `json:"name" yaml:"name"`
	ARN               string            `json:"arn" yaml:"arn"`
	DNSName           string            `json:"dnsName" yaml:"dnsName"`
	Type              string            `json:"type" yaml:"type"`
	Scheme            string            `json:"scheme" yaml:"scheme"`
	State             string            `json:"state" yaml:"state"`
	VpcId             string            `json:"vpcId" yaml:"vpcId"`
	AvailabilityZones []string          `json:"availabilityZones" yaml:"availabilityZones"`
	SecurityGroups    []string          `json:"securityGroups" yaml:"securityGroups"`
	IpAddressType     string            `json:"ipAddressType" yaml:"ipAddressType"`
	CustomerOwnedIpv4Pool string        `json:"customerOwnedIpv4Pool" yaml:"customerOwnedIpv4Pool"`
	Tags              map[string]string `json:"tags" yaml:"tags"`
	Listeners         []ListenerInfo    `json:"listeners" yaml:"listeners"`
	TargetGroups      []TargetGroupInfo `json:"targetGroups" yaml:"targetGroups"`
}

type ListenerInfo struct {
	Port     int32  `json:"port" yaml:"port"`
	Protocol string `json:"protocol" yaml:"protocol"`
	ARN      string `json:"arn" yaml:"arn"`
}

type TargetGroupInfo struct {
	Name       string `json:"name" yaml:"name"`
	ARN        string `json:"arn" yaml:"arn"`
	Port       int32  `json:"port" yaml:"port"`
	Protocol   string `json:"protocol" yaml:"protocol"`
	TargetType string `json:"targetType" yaml:"targetType"`
}
