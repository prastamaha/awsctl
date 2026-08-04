package s3

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/thediveo/klo"
	"github.com/urfave/cli/v3"
)

func (s *S3) GetCLI() *cli.Command {
	return &cli.Command{
		Name:    "bucket",
		Aliases: bucketAliases,
		Usage:   "Get S3 buckets",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			s.GetCommand()
			return nil
		},
	}
}

func (s *S3) GetCommand() {
	outputs := s.GetBuckets()
	if len(outputs) == 0 {
		fmt.Printf("No resources found in %s region\n", os.Getenv("AWS_REGION"))
		return
	}

	prn, err := klo.PrinterFromFlag("", &klo.Specs{DefaultColumnSpec: "NAME:{.Name},CREATION_DATE:{.CreationDate}"})
	if err != nil {
		panic(err)
	}

	table, err := klo.NewSortingPrinter("{.Name}", prn)
	if err != nil {
		panic(err)
	}
	table.Fprint(os.Stdout, outputs)
}

func (s *S3) GetBuckets() []BucketList {
	svc := s3.NewFromConfig(s.AWSConfig)

	resp, err := svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("failed to list buckets, %v", err)
	}

	var outputs []BucketList
	for _, bucket := range resp.Buckets {
		outputs = append(outputs, BucketList{
			Name:         *bucket.Name,
			CreationDate: bucket.CreationDate.String(),
		})
	}

	return outputs
}
