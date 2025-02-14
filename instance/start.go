package instance

import (
	"context"
	"fmt"
	"github/prastamaha/awsctl/utils"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/urfave/cli/v3"
)

func (in *Instance) StartCLI() *cli.Command {
	return &cli.Command{
		Name:    "instance",
		Aliases: instanceAliases,
		Usage:   "Start an ec2 instance",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Get(0) == "" {
				in.StartCommandFzf()
			} else {
				in.StartCommand(cmd.Args().Get(0))
			}
			return nil
		},
	}
}

func (in *Instance) StartCommand(id string) {
	client := ec2.NewFromConfig(in.AWSConfig)

	_, err := client.StartInstances(context.TODO(), &ec2.StartInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		log.Fatalf("failed to start instance, %v", err)
	}

	fmt.Printf("Starting instance %s...\n", id)
}

func (in *Instance) StartCommandFzf() {
	allInstances := in.GetInstance(
		ec2types.Filter{
			Name:   aws.String("instance-state-name"),
			Values: []string{"stopped"},
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

	data := utils.FuzzySearch("Select an instance to start: ", items)
	for _, i := range data {
		id := allInstances[i].ID
		in.StartCommand(id)
	}
}
