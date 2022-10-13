package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jimmysawczuk/paramstore/internal/ssm"
)

func main() {
	var path string
	var out string

	flag.StringVar(&path, "path", "", "path prefix for ssm")
	flag.StringVar(&out, "out", "", "output (leave blank for stdout)")

	flag.Parse()

	if path == "" || !strings.HasPrefix(path, "/") {
		log.Fatal("path must be present and start with /")
	}

	params, err := ssm.GetParametersFromPath(context.Background(), path)
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
