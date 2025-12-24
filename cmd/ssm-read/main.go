package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	ssmsvc "github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/jimmysawczuk/aws-tools/internal/ssm"
)

func main() {
	var path string
	var out string

	flag.StringVar(&path, "path", "", "path prefix for ssm")
	flag.StringVar(&out, "out", "", "output (leave blank for stdout)")

	flag.Parse()

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	ssmClient := ssmsvc.NewFromConfig(cfg)

	if path == "" || !strings.HasPrefix(path, "/") {
		log.Fatal("path must be present and start with /")
	}

	params, err := ssm.GetParametersFromPath(context.Background(), ssmClient, path)
	if err != nil {
		log.Fatal("ssm: get parameters from path", err)
	}

	log.Println(len(params), "parameters loaded")

	var w io.Writer = os.Stdout

	if out != "" {
		fp, err := os.Create(out)
		if err != nil {
			log.Fatal("os: open file (write):", err)
		}

		defer fp.Close()

		w = fp
	}

	for _, v := range params {
		fmt.Fprintf(w, "%s=%q\n", v.Name, v.Value)
	}
}
