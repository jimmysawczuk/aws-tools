package ssm

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type Param struct {
	Name   string
	Value  string
	Secure bool
}

func LoadParametersIntoPath(ctx context.Context, path string, params []Param) error {
	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("session: new session: %w", err)
	}

	ssmClient := ssm.New(sess)

	for _, param := range params {
		ty := ssm.ParameterTypeString
		if param.Secure {
			ty = ssm.ParameterTypeSecureString
		}

		_, err := ssmClient.PutParameterWithContext(ctx, &ssm.PutParameterInput{
			Name:      aws.String(path + "/" + param.Name),
			Value:     aws.String(param.Value),
			Type:      aws.String(ty),
			Overwrite: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("ssm: put parameter (%s): %w", path+"/"+param.Name, err)
		}
	}

	return nil
}

func GetParametersFromPath(ctx context.Context, path string) ([]Param, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("session: new session: %w", err)
	}

	ssmClient := ssm.New(sess)

	var tok *string
	var params []*ssm.Parameter

	for {

		res, err := ssmClient.GetParametersByPathWithContext(ctx, &ssm.GetParametersByPathInput{
			NextToken:      tok,
			Path:           aws.String(path),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			return nil, fmt.Errorf("ssm: get parameters by path: %w", err)
		}

		params = append(params, res.Parameters...)

		if res.NextToken == nil {
			break
		}

		tok = res.NextToken
	}

	tbr := make([]Param, len(params))
	for i, p := range params {
		tbr[i] = Param{
			Name:  strings.TrimLeft(strings.Replace(aws.StringValue(p.Name), path, "", 1), "/"),
			Value: aws.StringValue(p.Value),
		}
	}

	return tbr, nil
}

func LoadIntoEnv(in []Param) error {
	for _, v := range in {
		if err := os.Setenv(v.Name, v.Value); err != nil {
			return fmt.Errorf("os: setenv: %w", err)
		}
	}

	return nil
}
