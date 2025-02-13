package instance

import (
	"context"
	"fmt"
	"log"

	"github/prastamaha/awsctl/utils"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"gopkg.in/yaml.v3"
)

func (in *Instance) DescribeCommand(id string) {
	client := ec2.NewFromConfig(in.AWSConfig)

	resp, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		log.Fatalf("failed to list tables, %v", err)
	}

	if len(resp.Reservations) == 0 {
		fmt.Printf("Error from server (NotFound): ec2 instance %s not found\n", id)
		return
	}

	instance := resp.Reservations[0].Instances[0]

	publicIP := ""
	if instance.PublicIpAddress != nil {
		publicIP = *instance.PublicIpAddress
	}

	dnsPublic := ""
	if instance.PublicDnsName == nil {
		dnsPublic = *instance.PublicDnsName
	}

	name := ""
	if len(instance.Tags) > 0 {
		for _, tag := range instance.Tags {
			if *tag.Key == "Name" {
				name = *tag.Value
				break
			}
		}
	}

	outputs := EC2InstanceDescribe{
		Name:       name,
		IAMRoleArn: *instance.IamInstanceProfile.Arn,
		ID:         *instance.InstanceId,
		State:      string(instance.State.Name),
		Type:       string(instance.InstanceType),
		SubnetID:   *instance.SubnetId,
		VPCID:      *instance.VpcId,
		PrivateIP:  *instance.PrivateIpAddress,
		PublicIP:   publicIP,
		DNSPublic:  dnsPublic,
		DnsPrivate: *instance.PrivateDnsName,
		LaunchTime: instance.LaunchTime.String(),
		Tags:       utils.ConvertTags(instance.Tags),
		ImageID:    *instance.ImageId,
	}

	yamlData, err := yaml.Marshal(outputs)
	if err != nil {
		fmt.Println("Error marshalling to YAML:", err)
		return
	}
	fmt.Println(string(yamlData))
}

func (in *Instance) DescribeCommandFzf() {
	allInstances := in.GetInstance()
	if len(allInstances) == 0 {
		fmt.Println("No instances found")
		return
	}

	items := make([]string, len(allInstances))
	for i, v := range allInstances {
		items[i] = fmt.Sprintf("%s %s", v.ID, v.Name)
	}

	data := utils.FuzzySearch("Select an instance to describe: ", items)
	for _, i := range data {
		id := allInstances[i].ID
		in.DescribeCommand(id)
	}
}
