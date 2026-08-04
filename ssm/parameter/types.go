package parameter

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

var parameterAliases = []string{"param", "params", "parameters"}

type Parameter struct {
	AWSConfig aws.Config
}

type ParameterList struct {
	Name        string `json:"name" yaml:"name"`
	Type        string `json:"type" yaml:"type"`
	Description string `json:"description" yaml:"description"`
	LastModified string `json:"lastModified" yaml:"lastModified"`
	Version     int64  `json:"version" yaml:"version"`
}

type ParameterDescribe struct {
	Name           string            `json:"name" yaml:"name"`
	Type           string            `json:"type" yaml:"type"`
	Value          string            `json:"value" yaml:"value"`
	Version        int64             `json:"version" yaml:"version"`
	LastModified   string            `json:"lastModified" yaml:"lastModified"`
	ARN            string            `json:"arn" yaml:"arn"`
	Description    string            `json:"description" yaml:"description"`
	Tags           map[string]string `json:"tags" yaml:"tags"`
}
