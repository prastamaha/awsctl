package ssm

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

var ssmAliases = []string{"session", "ssh"}

type SSM struct {
	AWSConfig aws.Config
}
