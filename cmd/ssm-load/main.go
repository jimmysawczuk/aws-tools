package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/jimmysawczuk/aws-tools/internal/ssm"
	"github.com/joho/godotenv"
)

func main() {
	var path string
	var dryRun bool

	flag.StringVar(&path, "path", "", "path prefix for ssm")
	flag.BoolVar(&dryRun, "dry-run", true, "set to false to actually write to parameter store")

	flag.Parse()

	if path == "" || !strings.HasPrefix(path, "/") {
		log.Fatal("path must be present and start with /")
	}

	envPath := flag.Arg(0)

	fp, err := os.Open(envPath)
	if err != nil {
		log.Fatal("os: open", err)
	}

	defer fp.Close()

	res, err := godotenv.Parse(fp)
	if err != nil {
		log.Fatal("godotenv: parse", err)
	}

	params := []ssm.Param{}
	for k, v := range res {
		params = append(params, ssm.Param{
			Name:   k,
			Value:  v,
			Secure: true,
		})
	}

	if !dryRun {
		if err := ssm.LoadParametersIntoPath(context.Background(), path, params); err != nil {
			log.Fatal("ssm: load parameters into path", err)
		}
		return
	}

	for _, p := range params {
		log.Println("setting", path+"/"+p.Name+":", strings.Repeat("*", len(p.Value)))
	}
}
