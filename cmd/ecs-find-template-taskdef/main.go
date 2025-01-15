package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ecssvc "github.com/aws/aws-sdk-go/service/ecs"
)

func main() {
	var expectedTag string
	var imagePlaceholder string

	flag.StringVar(&expectedTag, "tag", "CreatedBy=Terraform", "the flag to look for that indicates a template taskdef")
	flag.StringVar(&imagePlaceholder, "image", "<IMAGE1_NAME>", "the placeholder to use for the image in the container definition")
	flag.Parse()

	family := flag.Arg(0)

	tag, val, ok := strings.Cut(expectedTag, "=")
	if !ok {
		log.Fatal(fmt.Errorf("expected tag should be of format tag=value"))
	}

	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(fmt.Errorf("session: new session: %w", err))
	}

	ecs := ecssvc.New(sess)

	taskDef, err := findTaskDefinition(context.Background(), ecs, family, tag, val)
	if err != nil {
		log.Fatal(fmt.Errorf("find task definition: %w", err))
	}

	def := buildDefinition(taskDef, imagePlaceholder)

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(def)
}

func buildDefinition(taskDef *ecssvc.TaskDefinition, placeholder string) TaskDefinition {
	def := TaskDefinition{
		TaskDefinitionARN: aws.StringValue(taskDef.TaskDefinitionArn),
		ExecutionRoleARN:  aws.StringValue(taskDef.ExecutionRoleArn),
		TaskRoleARN:       aws.StringValue(taskDef.TaskRoleArn),
		Compatibilities:   aws.StringValueSlice(taskDef.Compatibilities),
		NetworkMode:       aws.StringValue(taskDef.NetworkMode),
		CPU:               aws.StringValue(taskDef.Cpu),
		Memory:            aws.StringValue(taskDef.Memory),
		Family:            aws.StringValue(taskDef.Family),
	}

	for _, c := range taskDef.ContainerDefinitions {
		cdef := ContainerDefinition{
			Name:      aws.StringValue(c.Name),
			Image:     placeholder,
			Essential: aws.BoolValue(c.Essential),
		}

		for _, p := range c.PortMappings {
			cdef.PortMappings = append(cdef.PortMappings, PortMapping{
				HostPort:      uint64(aws.Int64Value(p.HostPort)),
				Protocol:      aws.StringValue(p.Protocol),
				ContainerPort: uint64(aws.Int64Value(p.ContainerPort)),
			})
		}

		def.ContainerDefinitions = append(def.ContainerDefinitions, cdef)
	}

	return def
}

type TaskDefinition struct {
	TaskDefinitionARN    string                `json:"taskDefinitionArn"`
	ExecutionRoleARN     string                `json:"executionRoleArn"`
	TaskRoleARN          string                `json:"taskRoleArn"`
	ContainerDefinitions []ContainerDefinition `json:"containerDefinitions"`
	Compatibilities      []string              `json:"compatibilities"`
	NetworkMode          string                `json:"networkMode"`
	CPU                  string                `json:"cpu"`
	Memory               string                `json:"memory"`
	Family               string                `json:"family"`
}

type ContainerDefinition struct {
	Name         string        `json:"name"`
	Image        string        `json:"image"`
	Essential    bool          `json:"essential"`
	PortMappings []PortMapping `json:"portMappings"`
}

type PortMapping struct {
	HostPort      uint64 `json:"hostPort"`
	Protocol      string `json:"tcp"`
	ContainerPort uint64 `json:"containerPort"`
}

func findTaskDefinition(ctx context.Context, ecs *ecssvc.ECS, family, tag, val string) (*ecssvc.TaskDefinition, error) {
	var token *string = nil

	for {
		res, err := ecs.ListTaskDefinitionsWithContext(ctx, &ecssvc.ListTaskDefinitionsInput{
			FamilyPrefix: aws.String(family),
			Sort:         aws.String(ecssvc.SortOrderDesc),
			Status:       aws.String(ecssvc.TaskDefinitionStatusActive),
			NextToken:    token,
		})
		if err != nil {
			log.Fatal(fmt.Errorf("ecs: list task definitions: %w", err))
		}

		for _, arn := range res.TaskDefinitionArns {
			res, err := ecs.DescribeTaskDefinitionWithContext(context.Background(), &ecssvc.DescribeTaskDefinitionInput{
				Include:        aws.StringSlice([]string{ecssvc.TaskDefinitionFieldTags}),
				TaskDefinition: arn,
			})
			if err != nil {
				log.Fatal(fmt.Errorf("ecs: describe task definition (%s): %w", aws.StringValue(arn), err))
			}

			for _, t := range res.Tags {
				if aws.StringValue(t.Key) == tag && aws.StringValue(t.Value) == val {
					return res.TaskDefinition, nil
				}
			}
		}

		if res.NextToken == nil {
			break
		}

		token = res.NextToken
	}

	return nil, fmt.Errorf("not found")
}
