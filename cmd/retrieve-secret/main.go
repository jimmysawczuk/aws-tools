package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func main() {
	var out string
	flag.StringVar(&out, "out", "", "filename to direct output (stdout if left blank)")

	flag.Parse()
	secretName := flag.Arg(0)

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	// Create Secrets Manager client
	sm := secretsmanager.NewFromConfig(cfg)

	rd, err := getSecret(ctx, sm, secretName)
	if err != nil {
		log.Fatalf("couldn't get secret: %s", err)
	}

	var wr io.Writer = os.Stdout
	if out != "" {
		fp, err := os.Create(out)
		if err != nil {
			log.Fatalf("couldn't open file for writing: %s", fmt.Errorf("os: create: %w", err))
		}

		wr = fp
	}

	if _, err := io.Copy(wr, rd); err != nil {
		log.Fatalf("couldn't copy to stdout: %s", err)
	}
}

func getSecret(ctx context.Context, sm *secretsmanager.Client, name string) (io.Reader, error) {
	resp, err := sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	})
	if err != nil {
		return nil, fmt.Errorf("secrets manager: %w", err)
	}

	return strings.NewReader(aws.ToString(resp.SecretString)), nil
}
