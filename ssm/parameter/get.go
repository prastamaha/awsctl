package parameter

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/thediveo/klo"
	"github.com/urfave/cli/v3"
)

func (p *Parameter) GetCLI() *cli.Command {
	return &cli.Command{
		Name:    "parameter",
		Aliases: parameterAliases,
		Usage:   "Get SSM parameters",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Parameter path (e.g., /my/path)",
			},
			&cli.StringFlag{
				Name:    "export",
				Aliases: []string{"e"},
				Usage:   "Export parameters to file (key=value format)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			path := cmd.String("path")
			export := cmd.String("export")
			p.GetCommand(path, export)
			return nil
		},
	}
}

func (p *Parameter) GetCommand(path string, export string) {
	if export != "" {
		p.ExportParameters(path, export)
		return
	}

	outputs := p.GetParameters(path)
	if len(outputs) == 0 {
		fmt.Printf("No resources found in %s region\n", os.Getenv("AWS_REGION"))
		return
	}

	prn, err := klo.PrinterFromFlag("", &klo.Specs{DefaultColumnSpec: "NAME:{.Name},TYPE:{.Type},LAST_MODIFIED:{.LastModified},VERSION:{.Version}"})
	if err != nil {
		panic(err)
	}

	table, err := klo.NewSortingPrinter("{.Name}", prn)
	if err != nil {
		panic(err)
	}
	table.Fprint(os.Stdout, outputs)
}

func (p *Parameter) GetParameters(path string) []ParameterList {
	svc := ssm.NewFromConfig(p.AWSConfig)

	if path == "" {
		path = "/"
	}

	var outputs []ParameterList
	var nextToken *string

	for {
		resp, err := svc.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
			Path:      aws.String(path),
			Recursive: aws.Bool(true),
			NextToken: nextToken,
		})
		if err != nil {
			log.Fatalf("failed to list parameters, %v", err)
		}

		for _, param := range resp.Parameters {
			lastModified := ""
			if param.LastModifiedDate != nil {
				lastModified = param.LastModifiedDate.String()
			}

			outputs = append(outputs, ParameterList{
				Name:         *param.Name,
				Type:         string(param.Type),
				LastModified: lastModified,
				Version:      param.Version,
			})
		}

		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	return outputs
}

func (p *Parameter) ExportParameters(path string, filename string) {
	svc := ssm.NewFromConfig(p.AWSConfig)

	if path == "" {
		path = "/"
	}

	var params []ssmtypes.Parameter
	var nextToken *string

	for {
		resp, err := svc.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
			Path:           aws.String(path),
			Recursive:      aws.Bool(true),
			WithDecryption: aws.Bool(true),
			NextToken:      nextToken,
		})
		if err != nil {
			log.Fatalf("failed to get parameters, %v", err)
		}

		params = append(params, resp.Parameters...)

		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	if len(params) == 0 {
		fmt.Printf("No parameters found in path %s\n", path)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed to create file, %v", err)
	}
	defer file.Close()

	for _, param := range params {
		line := fmt.Sprintf("%s=%s\n", *param.Name, *param.Value)
		if _, err := file.WriteString(line); err != nil {
			log.Fatalf("failed to write to file, %v", err)
		}
	}

	fmt.Printf("Exported %d parameters to %s\n", len(params), filename)
}
