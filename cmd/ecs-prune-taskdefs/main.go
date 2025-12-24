package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ecssvc "github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/kelseyhightower/envconfig"
)

var cfg struct {
	ClusterName  string   `envconfig:"CLUSTER_NAME" required:"true"`
	ServiceNames []string `envconfig:"SERVICE_NAMES" required:"true"`
}

func main() {
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(fmt.Errorf("couldn't process env: %w", err))
	}

	ctx := context.Background()

	awscfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	cl := ecssvc.NewFromConfig(awscfg)

	for _, s := range cfg.ServiceNames {
		obsolete, err := findObsoleteTaskdefs(context.Background(), cl, s)
		if err != nil {
			log.Fatalf("couldn't find obsolete taskdefs: %s", err)
		}

		log.Println("obsolete:")
		for _, d := range obsolete {
			log.Println(" - ", d)
		}

		for _, d := range obsolete {
			log.Println("deactivating", d)
			if err := deactivateTaskdef(context.Background(), cl, d); err != nil {
				log.Fatalf("couldn't deactivate taskdef (%s): %s", d, err)
			}

			time.Sleep(500 * time.Millisecond)
		}

		deactivated, err := findDeactivatedTaskdefs(context.Background(), cl, s)
		if err != nil {
			log.Fatalf("couldn't find deactivated taskdefs: %s", err)
		}

		log.Println("will be deleted:")
		for _, d := range deactivated {
			log.Println(" - ", d)
		}

		for i := 0; i < len(deactivated); i += 10 {
			var sl []string
			if i+10 > len(deactivated) {
				sl = deactivated[i:]
			} else {
				sl = deactivated[i : i+10]
			}

			log.Println("deleting", sl)
			if err := deleteTaskdefs(context.Background(), cl, sl); err != nil {
				log.Fatalf("couldn't delete taskdefs: %s", err)
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}

func findObsoleteTaskdefs(ctx context.Context, cl *ecssvc.Client, service string) ([]string, error) {
	var next *string
	var tbr []string
	var latestArn string

	for {
		defs, err := cl.ListTaskDefinitions(ctx, &ecssvc.ListTaskDefinitionsInput{
			NextToken:    next,
			FamilyPrefix: aws.String(service),
			Status:       ecstypes.TaskDefinitionStatusActive,
			Sort:         ecstypes.SortOrderDesc,
		})
		if err != nil {
			return nil, fmt.Errorf("ecs: list task defs: %w", err)
		}

		for _, arn := range defs.TaskDefinitionArns {
			if latestArn != "" {
				tbr = append(tbr, arn)
				continue
			}

			def, err := cl.DescribeTaskDefinition(ctx, &ecssvc.DescribeTaskDefinitionInput{
				Include:        []ecstypes.TaskDefinitionField{ecstypes.TaskDefinitionFieldTags},
				TaskDefinition: aws.String(arn),
			})
			if err != nil {
				return nil, fmt.Errorf("ecs: describe task def (%s): %w", arn, err)
			}

			for _, t := range def.Tags {
				if aws.ToString(t.Key) == "CreatedBy" && aws.ToString(t.Value) == "Terraform" {
					latestArn = arn
					break
				}
			}
		}

		next = defs.NextToken
		if next == nil {
			break
		}
	}

	return tbr, nil
}

func findDeactivatedTaskdefs(ctx context.Context, cl *ecssvc.Client, service string) ([]string, error) {
	var next *string
	var tbr []string

	for {
		defs, err := cl.ListTaskDefinitions(ctx, &ecssvc.ListTaskDefinitionsInput{
			NextToken:    next,
			FamilyPrefix: aws.String(service),
			Status:       ecstypes.TaskDefinitionStatusInactive,
			Sort:         ecstypes.SortOrderDesc,
		})
		if err != nil {
			return nil, fmt.Errorf("ecs: list task defs: %w", err)
		}

		tbr = append(tbr, defs.TaskDefinitionArns...)

		next = defs.NextToken
		if next == nil {
			break
		}
	}

	return tbr, nil
}

func deleteTaskdefs(ctx context.Context, cl *ecssvc.Client, del []string) error {
	if _, err := cl.DeleteTaskDefinitions(ctx, &ecssvc.DeleteTaskDefinitionsInput{
		TaskDefinitions: del,
	}); err != nil {
		return fmt.Errorf("ecs: delete task def (%s): %w", del, err)
	}

	return nil
}

func deactivateTaskdef(ctx context.Context, cl *ecssvc.Client, del string) error {
	if _, err := cl.DeregisterTaskDefinition(ctx, &ecssvc.DeregisterTaskDefinitionInput{
		TaskDefinition: aws.String(del),
	}); err != nil {
		return fmt.Errorf("ecs: deregister task def (%s): %w", del, err)
	}
	return nil
}
