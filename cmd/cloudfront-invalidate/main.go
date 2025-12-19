package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cloudfrontsvc "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cloudfronttypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

var cloudfront *cloudfrontsvc.Client

func main() {
	flag.Parse()

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	cloudfront = cloudfrontsvc.NewFromConfig(cfg)

	dist, path, err := parseArgs(ctx, flag.Args())
	if err != nil {
		log.Fatalf("couldn't parse args: %s", err)
	}

	log.Printf("%s: %s", dist.ID, dist.Comment)
	for _, v := range dist.Aliases {
		log.Printf(" - %s", v)
	}

	invalidation, err := invalidateDistribution(ctx, dist.ID, path)
	if err != nil {
		log.Fatalf("couldn't invalidate distribution: %s", err)
		return
	}

	log.Printf("invalidation created: %s", invalidation)

	for {
		resp, err := cloudfront.GetInvalidation(ctx, &cloudfrontsvc.GetInvalidationInput{
			DistributionId: aws.String(dist.ID),
			Id:             aws.String(invalidation),
		})
		if err != nil {
			log.Printf("couldn't get status: %s", err)
			time.Sleep(5 * time.Second)
		}

		if aws.ToString(resp.Invalidation.Status) == "Completed" {
			log.Println("invalidation complete")
			break
		}

		log.Println("waiting on invalidation to complete")
		time.Sleep(5 * time.Second)
	}
}

type Distribution struct {
	ID      string
	ARN     string
	Comment string
	Aliases []string
}

func getDistribution(ctx context.Context, id string) (*Distribution, error) {
	resp, err := cloudfront.GetDistribution(ctx, &cloudfrontsvc.GetDistributionInput{
		Id: aws.String(id),
	})
	if err != nil {
		return nil, fmt.Errorf("aws: cloudfront: get distribution: %w", err)
	}

	dist := Distribution{
		ID:      aws.ToString(resp.Distribution.Id),
		ARN:     aws.ToString(resp.Distribution.ARN),
		Comment: aws.ToString(resp.Distribution.DistributionConfig.Comment),
	}

	dist.Aliases = append(dist.Aliases, resp.Distribution.DistributionConfig.Aliases.Items...)

	return &dist, nil
}

func invalidateDistribution(ctx context.Context, id string, path string) (string, error) {
	resp, err := cloudfront.CreateInvalidation(ctx, &cloudfrontsvc.CreateInvalidationInput{
		DistributionId: aws.String(id),
		InvalidationBatch: &cloudfronttypes.InvalidationBatch{
			CallerReference: aws.String(time.Now().Format("20060102150405")),
			Paths: &cloudfronttypes.Paths{
				Items:    []string{path},
				Quantity: aws.Int32(1),
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("aws: cloudfront: create invalidation: %w", err)
	}

	return aws.ToString(resp.Invalidation.Id), nil
}

func parseArgs(ctx context.Context, args []string) (*Distribution, string, error) {
	if len(args) < 1 {
		return nil, "", fmt.Errorf("at least one argument is required")
	}

	id := args[0]
	path := "/*"
	if len(args) > 1 {
		path = args[1]
	}

	dist, err := getDistribution(ctx, id)
	if err == nil {
		return dist, path, nil
	}

	if nsd := new(cloudfronttypes.NoSuchDistribution); !errors.As(err, &nsd) {
		return nil, "", fmt.Errorf("get distribution: %w", err)
	}

	distID, err := findDistribution(ctx, id)
	if err != nil {
		return nil, "", fmt.Errorf("find distribution: %w", err)
	}

	dist, err = getDistribution(ctx, distID)
	if err != nil {
		return nil, "", fmt.Errorf("get distribution: %w", err)
	}

	return dist, path, nil
}

func findDistribution(ctx context.Context, domain string) (string, error) {
	dists, err := cloudfront.ListDistributions(ctx, &cloudfrontsvc.ListDistributionsInput{})
	if err != nil {
		return "", fmt.Errorf("list distributions: %w", err)
	}

	for _, d := range dists.DistributionList.Items {
		for _, a := range d.Aliases.Items {
			if a == domain {
				return aws.ToString(d.Id), nil
			}
		}
	}

	return "", fmt.Errorf("not found")
}
