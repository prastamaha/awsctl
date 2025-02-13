package instance

import (
	"context"
	"fmt"
	"github/prastamaha/awsctl/utils"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
)

func (in *Instance) RestartCommand(id string) {
	client := ec2.NewFromConfig(in.AWSConfig)

	_, err := client.RebootInstances(context.TODO(), &ec2.RebootInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		log.Fatalf("failed to reboot instance, %v", err)
	}

	fmt.Printf("Restarting instance %s...\n", id)
}

func (in *Instance) RestartCommandFzf() {
	allInstances := in.GetInstance(
		ec2types.Filter{
			Name:   aws.String("instance-state-name"),
			Values: []string{"running"},
		},
	)
	if len(allInstances) == 0 {
		fmt.Println("No instances found")
		return
	}

	items := make([]string, len(allInstances))
	for i, v := range allInstances {
		items[i] = fmt.Sprintf("%s %s", v.ID, v.Name)
	}

	data := utils.FuzzySearch("Select an instance to restart: ", items)
	for _, i := range data {
		id := allInstances[i].ID
		in.RestartCommand(id)
	}
}
