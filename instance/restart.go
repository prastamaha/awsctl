package instance

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/koki-develop/go-fzf"
)

func (i *Instance) RestartCommand(id string) {
	client := ec2.NewFromConfig(i.AWSConfig)

	_, err := client.RebootInstances(context.TODO(), &ec2.RebootInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		log.Fatalf("failed to reboot instance, %v", err)
	}

	fmt.Printf("Restarting instance %s...\n", id)
}

func (i *Instance) RestartCommandFzf() {
	client := ec2.NewFromConfig(i.AWSConfig)

	f, err := fzf.New(fzf.WithPrompt("Select an instance to restart: "))
	if err != nil {
		log.Fatal(err)
	}
	runningInstance := i.GetRunningInstance()
	if len(runningInstance) == 0 {
		fmt.Println("No running instances found")
		return
	}

	items := make([]string, len(runningInstance))
	for i, v := range runningInstance {
		items[i] = fmt.Sprintf("%s %s", v.ID, v.Name)
	}

	idxs, err := f.Find(items, func(i int) string { return items[i] })
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range idxs {
		id := runningInstance[i].ID
		_, err := client.RebootInstances(context.TODO(), &ec2.RebootInstancesInput{
			InstanceIds: []string{id},
		})
		if err != nil {
			log.Fatalf("failed to reboot instance, %v", err)
		}
		fmt.Printf("Restarting instance %s...\n", id)
	}
}
