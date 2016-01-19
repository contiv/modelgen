/***
Copyright 2014 Cisco Systems Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/contiv/modelgen/generators"
)

var app *cli.App

func main() {
	app = cli.NewApp()
	app.Usage = "Generate REST scaffolding from a json-described object model"
	app.ArgsUsage = "[source directory] [output directory]"
	app.HideHelp = true
	app.Version = ""
	app.Action = run
	app.Run(os.Args)
}

func run(ctx *cli.Context) {
	if len(ctx.Args()) != 2 {
		cli.ShowAppHelp(ctx)
		os.Exit(1)
	}

	sourceDir := ctx.Args()[0]
	if sourceDir == "" {
		sourceDir = "."
	}

	outputDir := ctx.Args()[1]
	if outputDir == "" {
		outputDir = "output"
	}

	if err := generators.ParseTemplates(); err != nil {
		panic(err)
	}

	var schema *Schema

	// Parse all files in input directory
	err := filepath.Walk(sourceDir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore non-json files and ignore files in subdirs
		if filepath.Ext(path) != ".json" || filepath.Dir(path) != filepath.Dir(sourceDir) {
			return nil
		}

		fmt.Printf("Parsing file: %q\n", path)

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Parse the schema
		sch, err := ParseSchema(b)
		if err != nil {
			return err
		}

		// Append to global schema
		schema = MergeSchema(schema, sch)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if schema == nil {
		log.Fatal("Could not find schema, aborting.")
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	// create client dir
	clientDir := path.Join(outputDir, "client")
	if err := os.MkdirAll(clientDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	outputs := map[string][]byte{}
	funcs := map[string]func() ([]byte, error){
		"%s.go":              schema.GenerateGo,
		"client/%s.js":       schema.GenerateJs,
		"client/%sClient.go": schema.GenerateClient,
		"client/%sClient.py": schema.GeneratePythonClient,
	}

	for frmt, fun := range funcs {
		str, err := fun()
		if err != nil {
			log.Fatalf("Error generating output for target %q: %v", frmt, err)
		}

		outFilename := fmt.Sprintf(frmt, schema.Name)
		outputs[outFilename] = []byte(str)
	}

	for outFilename, content := range outputs {
		target := path.Join(outputDir, outFilename)
		fmt.Printf("Generating file: %q\n", target)

		if err := ioutil.WriteFile(target, content, 0666); err != nil {
			log.Fatal(err)
		}
	}
}
