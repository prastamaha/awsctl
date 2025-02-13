package instance

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/koki-develop/go-fzf"
)

func (i *Instance) StartCommand(id string) {
	client := ec2.NewFromConfig(i.AWSConfig)

	_, err := client.StartInstances(context.TODO(), &ec2.StartInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		log.Fatalf("failed to start instance, %v", err)
	}

	fmt.Printf("Starting instance %s...\n", id)
}

func (i *Instance) StartCommandFzf() {
	client := ec2.NewFromConfig(i.AWSConfig)

	f, err := fzf.New(fzf.WithPrompt("Select an instance to start: "))
	if err != nil {
		log.Fatal(err)
	}

	stoppedInstance := i.GetStoppedInstance()
	if len(stoppedInstance) == 0 {
		fmt.Println("No stopped instances found")
		return
	}

	items := make([]string, len(stoppedInstance))
	for i, v := range stoppedInstance {
		items[i] = fmt.Sprintf("%s %s", v.ID, v.Name)
	}

	idxs, err := f.Find(items, func(i int) string { return items[i] })
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range idxs {
		id := stoppedInstance[i].ID
		_, err := client.StartInstances(context.TODO(), &ec2.StartInstancesInput{
			InstanceIds: []string{id},
		})
		if err != nil {
			log.Fatalf("failed to start instance, %v", err)
		}
		fmt.Printf("Starting instance %s...\n", id)
	}
}
