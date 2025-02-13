package ssm

import (
	"fmt"
	"github/prastamaha/awsctl/instance"
	"log"

	"github.com/koki-develop/go-fzf"
	"github.com/mmmorris1975/ssm-session-client/ssmclient"
)

func (s *SSM) StartSessionFzf() {
	cfg := s.AWSConfig

	f, err := fzf.New(fzf.WithPrompt("Select an instance to start session: "))
	if err != nil {
		log.Fatal(err)
	}

	ins := instance.Instance{
		AWSConfig: cfg,
	}

	// convert instanceList to []string of Name and ID
	instanceList := ins.GetRunningInstance()
	items := make([]string, len(instanceList))
	for i, v := range instanceList {
		items[i] = fmt.Sprintf("%s %s", v.ID, v.Name)
	}

	idxs, err := f.Find(items, func(i int) string { return items[i] })
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range idxs {
		target := instanceList[i].ID
		err := ssmclient.ShellSession(cfg, target)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *SSM) StartSessionTarget(target string) {
	cfg := s.AWSConfig

	err := ssmclient.ShellSession(cfg, target)
	if err != nil {
		log.Fatal(err)
	}
}
