package instance

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/thediveo/klo"
)

func (i *Instance) GetCommand() {
	outputs := i.GetInstance()
	if len(outputs) == 0 {
		fmt.Printf("No resources found in %s region\n", os.Getenv("AWS_REGION"))
		return
	}

	prn, err := klo.PrinterFromFlag("", &klo.Specs{DefaultColumnSpec: "NAME:{.Name},ID:{.ID},STATE:{.State},TYPE:{.Type},PRIVATE_IP:{.PrivateIP},PUBLIC_IP:{.PublicIP},LAUNCH_TIME:{.LaunchTime}"})
	if err != nil {
		panic(err)
	}

	table, err := klo.NewSortingPrinter("{.Name}", prn)
	if err != nil {
		panic(err)
	}
	table.Fprint(os.Stdout, outputs)
}

func (in *Instance) GetInstance(filter ...ec2types.Filter) []EC2InstanceList {
	svc := ec2.NewFromConfig(in.AWSConfig)

	var resp *ec2.DescribeInstancesOutput
	var err error
	if len(filter) > 0 {
		resp, err = svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
			Filters: filter,
		})
		if err != nil {
			log.Fatalf("failed to list tables, %v", err)
		}
	} else {
		resp, err = svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
		if err != nil {
			log.Fatalf("failed to list tables, %v", err)
		}
	}

	var outputs []EC2InstanceList
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			privateIP := ""
			if instance.PrivateIpAddress != nil {
				privateIP = *instance.PrivateIpAddress
			}

			publicIP := ""
			if instance.PublicIpAddress != nil {
				publicIP = *instance.PublicIpAddress
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

			outputs = append(outputs, EC2InstanceList{
				Name:       name,
				ID:         *instance.InstanceId,
				State:      string(instance.State.Name),
				Type:       string(instance.InstanceType),
				PrivateIP:  privateIP,
				PublicIP:   publicIP,
				LaunchTime: instance.LaunchTime.String(),
			})
		}
	}

	return outputs
}
