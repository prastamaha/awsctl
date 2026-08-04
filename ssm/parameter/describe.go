package parameter

import (
	"context"
	"fmt"
	"log"

	"github/prastamaha/awsctl/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func (p *Parameter) DescribeCLI() *cli.Command {
	return &cli.Command{
		Name:    "parameter",
		Aliases: parameterAliases,
		Usage:   "Describe an SSM parameter",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Get(0) == "" {
				p.DescribeCommandFzf()
			} else {
				p.DescribeCommand(cmd.Args().Get(0))
			}
			return nil
		},
	}
}

func (p *Parameter) DescribeCommand(name string) {
	svc := ssm.NewFromConfig(p.AWSConfig)

	resp, err := svc.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("failed to get parameter, %v", err)
	}

	if resp.Parameter == nil {
		fmt.Printf("Error from server (NotFound): parameter %s not found\n", name)
		return
	}

	param := resp.Parameter
	lastModified := ""
	arn := ""

	if param.LastModifiedDate != nil {
		lastModified = param.LastModifiedDate.String()
	}
	if param.ARN != nil {
		arn = *param.ARN
	}

	paramMeta := p.GetParameterMeta(name)

	tagsResp, err := svc.ListTagsForResource(context.TODO(), &ssm.ListTagsForResourceInput{
		ResourceType: ssmtypes.ResourceTypeForTaggingParameter,
		ResourceId:   &name,
	})

	tags := make(map[string]string)
	if err == nil && tagsResp.TagList != nil {
		tags = utils.ConvertSSMTags(tagsResp.TagList)
	}

	output := ParameterDescribe{
		Name:         *param.Name,
		Type:         string(param.Type),
		Value:        *param.Value,
		Version:      param.Version,
		LastModified: lastModified,
		ARN:          arn,
		Description:  paramMeta.Description,
		Tags:         tags,
	}

	yamlData, err := yaml.Marshal(output)
	if err != nil {
		fmt.Println("Error marshalling to YAML:", err)
		return
	}
	fmt.Println(string(yamlData))
}

func (p *Parameter) DescribeCommandFzf() {
	allParams := p.GetParameters("/")
	if len(allParams) == 0 {
		fmt.Println("No parameters found")
		return
	}

	items := make([]string, len(allParams))
	for i, v := range allParams {
		items[i] = fmt.Sprintf("%s (%s)", v.Name, v.Type)
	}

	data := utils.FuzzySearch("Select a parameter to describe: ", items)
	for _, i := range data {
		name := allParams[i].Name
		p.DescribeCommand(name)
	}
}

func (p *Parameter) GetParameterMeta(name string) ParameterList {
	svc := ssm.NewFromConfig(p.AWSConfig)

	resp, err := svc.DescribeParameters(context.TODO(), &ssm.DescribeParametersInput{
		Filters: []ssmtypes.ParametersFilter{
			{
				Key:    ssmtypes.ParametersFilterKeyName,
				Values: []string{name},
			},
		},
	})
	if err != nil {
		return ParameterList{}
	}

	for _, param := range resp.Parameters {
		if *param.Name == name {
			description := ""
			if param.Description != nil {
				description = *param.Description
			}
			lastModified := ""
			if param.LastModifiedDate != nil {
				lastModified = param.LastModifiedDate.String()
			}
			return ParameterList{
				Name:         *param.Name,
				Type:         string(param.Type),
				Description:  description,
				LastModified: lastModified,
				Version:      param.Version,
			}
		}
	}

	return ParameterList{}
}
