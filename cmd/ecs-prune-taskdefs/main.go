package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
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

	awscfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(fmt.Errorf("config: load default config: %w", err))
	}

	client := ecs.NewFromConfig(awscfg)

	for _, s := range cfg.ServiceNames {
		del, err := findObsoleteTaskdefs(context.Background(), client, s)
		if err != nil {
			log.Fatalf("couldn't find obsolete taskdefs: %s", err)
		}

		for _, d := range del {
			if err := deactivateTaskdef(context.Background(), client, d); err != nil {
				log.Fatalf("couldn't deactivate taskdef (%s): %s", d, err)
			}
			time.Sleep(500 * time.Millisecond)
		}

		for i := 0; i < len(del); i += 10 {
			var sl []string
			if i+10 > len(del) {
				sl = del[i:]
			} else {
				sl = del[i : i+10]
			}

			if err := deleteTaskdefs(context.Background(), client, sl); err != nil {
				log.Fatalf("couldn't delete taskdefs: %s", err)
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}

func findObsoleteTaskdefs(ctx context.Context, client *ecs.Client, service string) ([]string, error) {
	defs, err := client.ListTaskDefinitions(ctx, &ecs.ListTaskDefinitionsInput{
		FamilyPrefix: aws.String(service),
		Status:       ecstypes.TaskDefinitionStatusActive,
		Sort:         ecstypes.SortOrderDesc,
	})
	if err != nil {
		return nil, fmt.Errorf("ecs: list task defs: %w", err)
	}

	var latestArn string
	var delete []string

	for _, defArn := range defs.TaskDefinitionArns {
		if latestArn != "" {
			delete = append(delete, defArn)
			continue
		}

		def, err := client.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
			Include:        []ecstypes.TaskDefinitionField{ecstypes.TaskDefinitionFieldTags},
			TaskDefinition: aws.String(defArn),
		})
		if err != nil {
			return nil, fmt.Errorf("ecs: describe task def (%s): %w", defArn, err)
		}

		for _, t := range def.Tags {
			if aws.ToString(t.Key) == "CreatedBy" && aws.ToString(t.Value) == "Terraform" {
				latestArn = defArn
				break
			}
		}
	}

	return delete, nil
}

func deleteTaskdefs(ctx context.Context, client *ecs.Client, del []string) error {
	if _, err := client.DeleteTaskDefinitions(ctx, &ecs.DeleteTaskDefinitionsInput{
		TaskDefinitions: del,
	}); err != nil {
		return fmt.Errorf("ecs: delete task def (%s): %w", del, err)
	}

	return nil
}

func deactivateTaskdef(ctx context.Context, client *ecs.Client, del string) error {
	if _, err := client.DeregisterTaskDefinition(ctx, &ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: aws.String(del),
	}); err != nil {
		return fmt.Errorf("ecs: deregister task def (%s): %w", del, err)
	}
	return nil
}
