package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	ecssvc "github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/kelseyhightower/envconfig"
)

var cfg struct {
	ClusterName string `envconfig:"CLUSTER_NAME" required:"true"`
	ServiceName string `envconfig:"SERVICE_NAME" required:"true"`

	ContainerName string `envconfig:"CONTAINER_NAME" required:"true"`
	ContainerPort uint   `envconfig:"CONTAINER_PORT" default:"8080" required:"true"`
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

	ecs := ecssvc.NewFromConfig(awscfg)

	spec, err := buildAppspec(context.Background(), ecs, cfg.ClusterName, cfg.ServiceName)
	if err != nil {
		log.Fatal(fmt.Errorf("build appspec: %w", err))
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(spec); err != nil {
		log.Fatal(fmt.Errorf("json: encode: %w", err))
	}
}

func buildAppspec(ctx context.Context, ecs *ecssvc.Client, cluster string, service string) (*Appspec, error) {
	res, err := ecs.DescribeServices(ctx, &ecssvc.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: []string{service},
	})
	if err != nil {
		return nil, fmt.Errorf("ecs: describe services: %w (cluster: %s, service: %s)", err, cluster, service)
	}

	if len(res.Services) != 1 {
		return nil, fmt.Errorf("unexpected service count: %d (cluster: %s, service: %s)", len(res.Services), cluster, service)
	}

	svc := res.Services[0]

	resource := ResourceDef{
		Type: "AWS::ECS::Service",
		Properties: Properties{
			TaskDefinition: true,
			LoadBalancerInfo: LoadBalancerInfo{
				ContainerName: cfg.ContainerName,
				ContainerPort: cfg.ContainerPort,
			},
		},
	}

	for _, s := range svc.CapacityProviderStrategy {
		resource.Properties.CapacityProviderStrategy = append(resource.Properties.CapacityProviderStrategy, CapacityProviderStrategy{
			Base:             uint(s.Base),
			CapacityProvider: aws.ToString(s.CapacityProvider),
			Weight:           uint(s.Weight),
		})
	}

	return &Appspec{
		Resources: []map[string]ResourceDef{{
			"TargetService": resource,
		}},
	}, nil
}

type Appspec struct {
	Version   version `json:"version"`
	Resources []map[string]ResourceDef
}

type version bool

func (v version) MarshalJSON() ([]byte, error) {
	return []byte("0.0"), nil
}

type ResourceDef struct {
	Type       string
	Properties Properties
}

type Properties struct {
	TaskDefinition           TaskDefinition
	LoadBalancerInfo         LoadBalancerInfo
	CapacityProviderStrategy []CapacityProviderStrategy `json:",omitempty"`
}

type TaskDefinition bool

type LoadBalancerInfo struct {
	ContainerName string
	ContainerPort uint
}

type CapacityProviderStrategy struct {
	Base             uint
	CapacityProvider string
	Weight           uint
}

func (t TaskDefinition) MarshalJSON() ([]byte, error) {
	if !t {
		return nil, nil
	}

	return json.RawMessage(`"<TASK_DEFINITION>"`), nil
}
