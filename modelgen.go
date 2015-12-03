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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/contiv/modelgen/generators"
)

const defaultOutPath = "output"

var (
	sourceDir = flag.String("s", "./", "Location of json schema")
	outputDir = flag.String("o", "output", "Output directory")
)

func main() {
	flag.Parse()
	if err := generators.ParseTemplates(); err != nil {
		panic(err)
	}

	var schema *Schema

	// Parse all files in input directory
	err := filepath.Walk(*sourceDir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore non-json files
		if filepath.Ext(path) != ".json" {
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

	if *outputDir == "" {
		*outputDir = defaultOutPath
	}

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	outputs := map[string][]byte{}
	funcs := map[string]func() ([]byte, error){
		"go":        schema.GenerateGo,
		"js":        schema.GenerateJs,
		"client.go": schema.GenerateClient,
		"py":        schema.GeneratePythonClient,
	}

	for ext, fun := range funcs {
		str, err := fun()
		if err != nil {
			log.Fatalf("Error generating output for target %q: %v", ext, err)
		}

		outputs[strings.Join([]string{schema.Name, ext}, ".")] = []byte(str)
	}

	for fn, content := range outputs {
		target := path.Join(*outputDir, fn)
		fmt.Printf("Generating file: %q\n", target)

		if err := ioutil.WriteFile(target, content, 0666); err != nil {
			log.Fatal(err)
		}
	}
}
