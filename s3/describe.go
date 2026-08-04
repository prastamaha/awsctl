package s3

import (
	"context"
	"fmt"

	"github/prastamaha/awsctl/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func (s *S3) DescribeCLI() *cli.Command {
	return &cli.Command{
		Name:    "bucket",
		Aliases: bucketAliases,
		Usage:   "Describe an S3 bucket",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Get(0) == "" {
				s.DescribeCommandFzf()
			} else {
				s.DescribeCommand(cmd.Args().Get(0))
			}
			return nil
		},
	}
}

func (s *S3) DescribeCommand(bucket string) {
	svc := s3.NewFromConfig(s.AWSConfig)

	output := BucketDescribe{
		Name: bucket,
	}

	tagsResp, err := svc.GetBucketTagging(context.TODO(), &s3.GetBucketTaggingInput{
		Bucket: aws.String(bucket),
	})
	if err == nil && tagsResp.TagSet != nil {
		output.Tags = utils.ConvertS3Tags(tagsResp.TagSet)
	}

	versioningResp, err := svc.GetBucketVersioning(context.TODO(), &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucket),
	})
	if err == nil {
		if versioningResp.Status != "" {
			output.Versioning = string(versioningResp.Status)
		} else {
			output.Versioning = "Disabled"
		}
	}

	encryptionResp, err := svc.GetBucketEncryption(context.TODO(), &s3.GetBucketEncryptionInput{
		Bucket: aws.String(bucket),
	})
	if err == nil && encryptionResp.ServerSideEncryptionConfiguration != nil {
		if len(encryptionResp.ServerSideEncryptionConfiguration.Rules) > 0 {
			rule := encryptionResp.ServerSideEncryptionConfiguration.Rules[0]
			if rule.ApplyServerSideEncryptionByDefault != nil {
				output.Encryption = string(rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
			}
		}
	}

	pabResp, err := svc.GetPublicAccessBlock(context.TODO(), &s3.GetPublicAccessBlockInput{
		Bucket: aws.String(bucket),
	})
	if err == nil && pabResp.PublicAccessBlockConfiguration != nil {
		pab := pabResp.PublicAccessBlockConfiguration
		output.PublicAccessBlock = &PublicAccess{
			BlockPublicAcls:       *pab.BlockPublicAcls,
			IgnorePublicAcls:      *pab.IgnorePublicAcls,
			BlockPublicPolicy:     *pab.BlockPublicPolicy,
			RestrictPublicBuckets: *pab.RestrictPublicBuckets,
		}
	}

	yamlData, err := yaml.Marshal(output)
	if err != nil {
		fmt.Println("Error marshalling to YAML:", err)
		return
	}
	fmt.Println(string(yamlData))
}

func (s *S3) DescribeCommandFzf() {
	buckets := s.GetBuckets()
	if len(buckets) == 0 {
		fmt.Println("No buckets found")
		return
	}

	items := make([]string, len(buckets))
	for i, v := range buckets {
		items[i] = v.Name
	}

	data := utils.FuzzySearch("Select a bucket to describe: ", items)
	for _, i := range data {
		name := buckets[i].Name
		s.DescribeCommand(name)
	}
}
