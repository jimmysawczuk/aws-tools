package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ssmsvc "github.com/aws/aws-sdk-go/service/secretsmanager"
)

func main() {
	var out string
	flag.StringVar(&out, "out", "", "filename to direct output (stdout if left blank)")

	flag.Parse()
	secretName := flag.Arg(0)

	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(fmt.Errorf("session: new session: %w", err))
	}

	ssm := ssmsvc.New(sess)

	rd, err := getSecret(context.Background(), ssm, secretName)
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

func getSecret(ctx context.Context, ssm *ssmsvc.SecretsManager, name string) (io.Reader, error) {
	resp, err := ssm.GetSecretValueWithContext(ctx, &ssmsvc.GetSecretValueInput{
		SecretId: aws.String(name),
	})
	if err != nil {
		return nil, fmt.Errorf("secrets manager: %w", err)
	}

	return strings.NewReader(aws.StringValue(resp.SecretString)), nil
}
