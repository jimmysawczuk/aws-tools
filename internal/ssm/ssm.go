package ssm

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type Param struct {
	Name   string
	Value  string
	Secure bool
}

func LoadParametersIntoPath(ctx context.Context, cl *ssm.Client, path string, params []Param) error {
	for _, param := range params {
		ty := types.ParameterTypeString
		if param.Secure {
			ty = types.ParameterTypeSecureString
		}

		if param.Value == "" {
			continue
		}

		_, err := cl.PutParameter(ctx, &ssm.PutParameterInput{
			Name:      aws.String(path + "/" + param.Name),
			Value:     aws.String(param.Value),
			Type:      ty,
			Overwrite: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("ssm: put parameter (%s): %w", path+"/"+param.Name, err)
		}
	}

	return nil
}

func GetParametersFromPath(ctx context.Context, ssmClient *ssm.Client, path string) ([]Param, error) {
	var tok *string
	var params []types.Parameter

	for {

		res, err := ssmClient.GetParametersByPath(ctx, &ssm.GetParametersByPathInput{
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
			Name:  strings.TrimLeft(strings.Replace(aws.ToString(p.Name), path, "", 1), "/"),
			Value: aws.ToString(p.Value),
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

func DeleteParameters(ctx context.Context, ssmClient *ssm.Client, paths []string) error {
	res, err := ssmClient.DeleteParameters(ctx, &ssm.DeleteParametersInput{
		Names: paths,
	})
	if err != nil {
		return fmt.Errorf("ssm: delete parameters: %w", err)
	}

	if len(res.InvalidParameters) > 0 {
		return fmt.Errorf("invalid parameters: %v", res.InvalidParameters)
	}

	return nil
}
