package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ecssvc "github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
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

	ctx := context.Background()

	awscfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	ecs := ecssvc.NewFromConfig(awscfg)

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

func findTaskDefinition(ctx context.Context, ecs *ecssvc.Client, family, tag, val string) (*ecstypes.TaskDefinition, error) {
	var token *string = nil

	for {
		res, err := ecs.ListTaskDefinitions(ctx, &ecssvc.ListTaskDefinitionsInput{
			FamilyPrefix: aws.String(family),
			Sort:         ecstypes.SortOrderDesc,
			Status:       ecstypes.TaskDefinitionStatusActive,
			NextToken:    token,
		})
		if err != nil {
			log.Fatal(fmt.Errorf("ecs: list task definitions: %w", err))
		}

		for _, arn := range res.TaskDefinitionArns {
			res, err := ecs.DescribeTaskDefinition(ctx, &ecssvc.DescribeTaskDefinitionInput{
				Include:        []ecstypes.TaskDefinitionField{ecstypes.TaskDefinitionFieldTags},
				TaskDefinition: aws.String(arn),
			})
			if err != nil {
				log.Fatal(fmt.Errorf("ecs: describe task definition (%s): %w", arn, err))
			}

			for _, t := range res.Tags {
				if aws.ToString(t.Key) == tag && aws.ToString(t.Value) == val {
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

func buildDefinition(taskDef *ecstypes.TaskDefinition, placeholder string) TaskDefinition {
	def := TaskDefinition{
		TaskDefinitionARN: aws.ToString(taskDef.TaskDefinitionArn),
		ExecutionRoleARN:  aws.ToString(taskDef.ExecutionRoleArn),
		TaskRoleARN:       aws.ToString(taskDef.TaskRoleArn),
		Compatibilities:   taskDef.Compatibilities,
		NetworkMode:       string(taskDef.NetworkMode),
		CPU:               aws.ToString(taskDef.Cpu),
		Memory:            aws.ToString(taskDef.Memory),
		Family:            aws.ToString(taskDef.Family),
		PidMode:           string(taskDef.PidMode),
	}

	for i, c := range taskDef.ContainerDefinitions {
		cdef := ContainerDefinition{
			Name:              aws.ToString(c.Name),
			Image:             aws.ToString(c.Image),
			Essential:         aws.ToBool(c.Essential),
			CPU:               uint64(c.Cpu),
			Memory:            uint64(aws.ToInt32(c.Memory)),
			MemoryReservation: uint64(aws.ToInt32(c.MemoryReservation)),
		}

		if i == 0 {
			cdef.Image = placeholder
		}

		for _, p := range c.PortMappings {
			cdef.PortMappings = append(cdef.PortMappings, PortMapping{
				HostPort:      uint64(aws.ToInt32(p.HostPort)),
				Protocol:      string(p.Protocol),
				ContainerPort: uint64(aws.ToInt32(p.ContainerPort)),
			})
		}

		for _, e := range c.Environment {
			cdef.Environment = append(cdef.Environment, Environment{
				Name:  aws.ToString(e.Name),
				Value: aws.ToString(e.Value),
			})
		}

		if c.LogConfiguration != nil {
			cdef.LogConfiguration.LogDriver = string(c.LogConfiguration.LogDriver)
			cdef.LogConfiguration.Options = c.LogConfiguration.Options
		}

		if c.FirelensConfiguration != nil {
			cdef.FirelensConfiguration.Type = string(c.FirelensConfiguration.Type)
			cdef.FirelensConfiguration.Options = c.FirelensConfiguration.Options
		}

		def.ContainerDefinitions = append(def.ContainerDefinitions, cdef)
	}

	return def
}

type TaskDefinition struct {
	TaskDefinitionARN    string                   `json:"taskDefinitionArn"`
	ExecutionRoleARN     string                   `json:"executionRoleArn"`
	TaskRoleARN          string                   `json:"taskRoleArn"`
	ContainerDefinitions []ContainerDefinition    `json:"containerDefinitions"`
	Compatibilities      []ecstypes.Compatibility `json:"compatibilities"`
	NetworkMode          string                   `json:"networkMode"`
	CPU                  string                   `json:"cpu"`
	Memory               string                   `json:"memory"`
	Family               string                   `json:"family"`
	PidMode              string                   `json:"pidMode,omitempty"`
}

type ContainerDefinition struct {
	Name                  string                `json:"name"`
	Image                 string                `json:"image"`
	Essential             bool                  `json:"essential"`
	CPU                   uint64                `json:"cpu,omitzero"`
	Memory                uint64                `json:"memory,omitzero"`
	MemoryReservation     uint64                `json:"memoryReservation,omitzero"`
	PortMappings          []PortMapping         `json:"portMappings"`
	Environment           []Environment         `json:"environment"`
	LogConfiguration      LogConfiguration      `json:"logConfiguration,omitzero"`
	FirelensConfiguration FirelensConfiguration `json:"firelensConfiguration,omitzero"`
}

type PortMapping struct {
	HostPort      uint64 `json:"hostPort"`
	Protocol      string `json:"tcp"`
	ContainerPort uint64 `json:"containerPort"`
}

type Environment struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type LogConfiguration struct {
	LogDriver string            `json:"logDriver"`
	Options   map[string]string `json:"options"`
}

type FirelensConfiguration struct {
	Type    string            `json:"type"`
	Options map[string]string `json:"options"`
}
