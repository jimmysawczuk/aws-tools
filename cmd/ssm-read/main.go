package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/jimmysawczuk/paramstore/internal/ssm"
)

func main() {
	var path string

	flag.StringVar(&path, "path", "", "path prefix for ssm")

	flag.Parse()

	if path == "" || !strings.HasPrefix(path, "/") {
		log.Fatal("path must be present and start with /")
	}

	params, err := ssm.GetParametersFromPath(context.Background(), path)
	if err != nil {
		log.Fatal("ssm: get parameters from path", err)
	}

	log.Println(len(params), "parameters loaded")

	for _, v := range params {
		log.Println(v.Name + "=" + v.Value)
	}
}
