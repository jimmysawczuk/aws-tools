package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	ssmsvc "github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/jimmysawczuk/aws-tools/internal/ssm"
)

func main() {
	var path string
	var dryRun bool

	flag.StringVar(&path, "path", "", "path prefix for ssm")
	flag.BoolVar(&dryRun, "dry-run", true, "set to true to actually delete params")

	flag.Parse()

	if path == "" || !strings.HasPrefix(path, "/") {
		log.Fatal("path must be present and start with /")
	}

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	ssmClient := ssmsvc.NewFromConfig(cfg)

	params, err := ssm.GetParametersFromPath(context.Background(), ssmClient, path)
	if err != nil {
		log.Fatal("ssm: get parameters from path", err)
	}

	log.Println(len(params), "parameters found")
	for _, v := range params {
		log.Println(v.Name, v.Secure)
	}

	if dryRun {
		return
	}

	for i := 0; i < len(params); i += 10 {
		sl := params[i:]
		if len(sl) > 10 {
			sl = sl[:10]
		}

		// log.Println(names("", sl))

		if err := ssm.DeleteParameters(context.Background(), ssmClient, names(path+"/", sl)); err != nil {
			panic(err)
		}
	}
}

func names(prefix string, in []ssm.Param) []string {
	out := make([]string, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = prefix + in[i].Name
	}

	return out
}
