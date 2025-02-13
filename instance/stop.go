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

func (in *Instance) StopCommand(id string) {
	client := ec2.NewFromConfig(in.AWSConfig)

	_, err := client.StopInstances(context.TODO(), &ec2.StopInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		log.Fatalf("failed to start instance, %v", err)
	}

	fmt.Printf("Stopping instance %s...\n", id)
}

func (in *Instance) StopCommandFzf() {
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

	data := utils.FuzzySearch("Select an instance to stop: ", items)
	for _, i := range data {
		id := allInstances[i].ID
		in.StopCommand(id)
	}
}
