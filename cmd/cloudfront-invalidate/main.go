package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	cloudfrontsvc "github.com/aws/aws-sdk-go/service/cloudfront"
)

var (
	sess       *session.Session
	cloudfront *cloudfrontsvc.CloudFront
)

func main() {
	flag.Parse()

	sess = session.Must(session.NewSession())
	cloudfront = cloudfrontsvc.New(sess)
	ctx := context.Background()

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
		resp, err := cloudfront.GetInvalidationWithContext(ctx, &cloudfrontsvc.GetInvalidationInput{
			DistributionId: aws.String(dist.ID),
			Id:             aws.String(invalidation),
		})
		if err != nil {
			log.Printf("couldn't get status: %s", err)
			time.Sleep(5 * time.Second)
		}

		if aws.StringValue(resp.Invalidation.Status) == "Completed" {
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
	resp, err := cloudfront.GetDistributionWithContext(ctx, &cloudfrontsvc.GetDistributionInput{
		Id: aws.String(id),
	})
	if err != nil {
		return nil, fmt.Errorf("aws: cloudfront: get distribution: %w", err)
	}

	dist := Distribution{
		ID:      aws.StringValue(resp.Distribution.Id),
		ARN:     aws.StringValue(resp.Distribution.ARN),
		Comment: aws.StringValue(resp.Distribution.DistributionConfig.Comment),
	}

	for _, s := range resp.Distribution.DistributionConfig.Aliases.Items {
		dist.Aliases = append(dist.Aliases, aws.StringValue(s))
	}

	return &dist, nil
}

func invalidateDistribution(ctx context.Context, id string, path string) (string, error) {
	resp, err := cloudfront.CreateInvalidationWithContext(ctx, &cloudfrontsvc.CreateInvalidationInput{
		DistributionId: aws.String(id),
		InvalidationBatch: &cloudfrontsvc.InvalidationBatch{
			CallerReference: aws.String(time.Now().Format("20060102150405")),
			Paths: &cloudfrontsvc.Paths{
				Items:    aws.StringSlice([]string{path}),
				Quantity: aws.Int64(1),
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("aws: cloudfront: create invalidation: %w", err)
	}

	return aws.StringValue(resp.Invalidation.Id), nil
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

	var aerr awserr.Error
	if errors.As(err, &aerr) && aerr.Code() != cloudfrontsvc.ErrCodeNoSuchDistribution {
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
	dists, err := cloudfront.ListDistributionsWithContext(ctx, &cloudfrontsvc.ListDistributionsInput{})
	if err != nil {
		return "", fmt.Errorf("list distributions: %w", err)
	}

	for _, d := range dists.DistributionList.Items {
		for _, a := range d.Aliases.Items {
			if aws.StringValue(a) == domain {
				return aws.StringValue(d.Id), nil
			}
		}
	}

	return "", fmt.Errorf("not found")
}
