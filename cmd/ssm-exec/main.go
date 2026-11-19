package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	ssmsvc "github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/jimmysawczuk/aws-tools/internal/ssm"
)

func main() {
	var path string

	flag.StringVar(&path, "path", os.Getenv("SSM_PATH"), "path prefix for ssm; falls back to $SSM_PATH")
	flag.Usage = func() {
		log.SetFlags(0)
		log.Printf("usage: %s -path /ssm/prefix program [args...]", os.Args[0])
		log.Printf("       SSM_PATH=/ssm/prefix %s program [args...]", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if path == "" || !strings.HasPrefix(path, "/") {
		log.Fatal("path must be set via -path or SSM_PATH and start with /")
	}

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("program to exec is required")
	}

	bin, err := exec.LookPath(args[0])
	if err != nil {
		log.Fatalf("exec: look path (%s): %v", args[0], err)
	}

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	ssmClient := ssmsvc.NewFromConfig(cfg)

	params, err := ssm.GetParametersFromPath(ctx, ssmClient, path)
	if err != nil {
		log.Fatal("ssm: get parameters from path: ", err)
	}

	set, err := ssm.LoadIntoEnvNoOverwrite(params)
	if err != nil {
		log.Fatal("ssm: load into env: ", err)
	}

	log.Printf("%d parameters loaded (%d set, %d already in env)", len(params), set, len(params)-set)

	if err := syscall.Exec(bin, args, os.Environ()); err != nil {
		log.Fatalf("exec (%s): %v", bin, err)
	}
}
