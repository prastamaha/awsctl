package instance

import (
	"context"
	"fmt"
	"log"

	"github/prastamaha/awsctl/utils"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"gopkg.in/yaml.v3"
)

func (i *Instance) DescribeCommand(id string) {
	// Initialize Config
	client := ec2.NewFromConfig(i.AWSConfig)

	// Build the request with its input parameters
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
		// search from Tags with key "Name"
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
