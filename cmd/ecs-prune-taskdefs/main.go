package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
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

	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(fmt.Errorf("session: new session: %w", err))
	}

	cl := ecs.New(sess)

	for _, s := range cfg.ServiceNames {
		del, err := findObsoleteTaskdefs(context.Background(), cl, s)
		if err != nil {
			log.Fatalf("couldn't find obsolete taskdefs: %s", err)
		}

		for _, d := range del {
			if err := deactivateTaskdef(context.Background(), cl, d); err != nil {
				log.Fatalf("couldn't deactivate taskdef (%s): %s", d, err)
			}
			time.Sleep(500 * time.Millisecond)
		}

		for i := 0; i < len(del); i += 10 {
			// var sl []string
			// if i+10 > len(del) {
			// 	sl = del[i:]
			// } else {
			// 	sl = del[i : i+10]
			// }

			// if err := deleteTaskdefs(context.Background(), cl, sl); err != nil {
			// 	log.Fatalf("couldn't delete taskdefs: %s", err)
			// }

			time.Sleep(500 * time.Millisecond)
		}
	}
}

func findObsoleteTaskdefs(ctx context.Context, cl *ecs.ECS, service string) ([]string, error) {
	defs, err := cl.ListTaskDefinitionsWithContext(ctx, &ecs.ListTaskDefinitionsInput{
		FamilyPrefix: aws.String(service),
		Status:       aws.String(ecs.TaskDefinitionFamilyStatusActive),
		Sort:         aws.String(ecs.SortOrderDesc),
	})
	if err != nil {
		return nil, fmt.Errorf("ecs: list task defs: %w", err)
	}

	var latestArn string
	var delete []string

	for _, defArn := range defs.TaskDefinitionArns {
		if latestArn != "" {
			delete = append(delete, aws.StringValue(defArn))
			continue
		}

		def, err := cl.DescribeTaskDefinitionWithContext(ctx, &ecs.DescribeTaskDefinitionInput{
			Include:        aws.StringSlice([]string{ecs.TaskDefinitionFieldTags}),
			TaskDefinition: defArn,
		})
		if err != nil {
			return nil, fmt.Errorf("ecs: describe task def (%s): %w", defArn, err)
		}

		for _, t := range def.Tags {
			if aws.StringValue(t.Key) == "CreatedBy" && aws.StringValue(t.Value) == "Terraform" {
				latestArn = aws.StringValue(defArn)
				break
			}
		}
	}

	return delete, nil
}

func deleteTaskdefs(ctx context.Context, cl *ecs.ECS, del []string) error {
	if _, err := cl.DeleteTaskDefinitionsWithContext(ctx, &ecs.DeleteTaskDefinitionsInput{
		TaskDefinitions: aws.StringSlice(del),
	}); err != nil {
		return fmt.Errorf("ecs: delete task def (%s): %w", del, err)
	}

	return nil
}

func deactivateTaskdef(ctx context.Context, cl *ecs.ECS, del string) error {
	if _, err := cl.DeregisterTaskDefinitionWithContext(ctx, &ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: aws.String(del),
	}); err != nil {
		return fmt.Errorf("ecs: deregister task def (%s): %w", del, err)
	}
	return nil
}
