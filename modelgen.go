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

	log "github.com/Sirupsen/logrus"
	"github.com/contiv/objmodel/tools/modelgen/generators"
)

var (
	source = flag.String("s", "./", "Location of json schema")
	output = flag.String("o", "", "Output directory")
)

func main() {
	flag.Parse()
	if err := generators.ParseTemplates(); err != nil {
		panic(err)
	}

	var schema *Schema

	// Parse all files in input directory
	err := filepath.Walk(*source, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore non-json files
		if filepath.Ext(path) != ".json" || filepath.Dir(path) != filepath.Dir(*source) {
			return nil
		}

		fmt.Printf("Parsing file: %s\n", path)

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

	outPath := "./"
	if *output != "" {
		outPath = *output
		if err := os.MkdirAll(outPath, 0755); err != nil {
			log.Fatalf("Error creating output directory: %v", err)
		}
	}

	// Generate Go code
	outStr, err := schema.GenerateGo()
	if err != nil {
		log.Errorf("Error generating go structs. Err: %v", err)
		// XXX fallthrough so we can write the files
	}

	// Write the Go file output
	goFileName := path.Join(outPath, schema.Name+".go")
	fmt.Printf("Writing to file: %s\n", goFileName)
	err = ioutil.WriteFile(goFileName, []byte(outStr), 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Generate javascript
	outStr, err = schema.GenerateJs()
	if err != nil {
		log.Fatalf("Error generating javascript. Err: %v", err)
	}

	// Write javascript file
	jsFileName := path.Join(outPath, schema.Name+".js")
	fmt.Printf("Writing to file: %s\n", jsFileName)
	err = ioutil.WriteFile(jsFileName, []byte(outStr), 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Generate go client
	outStr, err = schema.GenerateClient()
	if err != nil {
		log.Fatalf("Error generating go client. Err: %v", err)
	}

	// Write go client file
	goClientFileName := path.Join(outPath, schema.Name+"Client.go")
	fmt.Printf("Writing to file: %s\n", goClientFileName)
	err = ioutil.WriteFile(goClientFileName, []byte(outStr), 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Generate python client
	outStr, err = schema.GeneratePythonClient()
	if err != nil {
		log.Fatalf("Error generating go client. Err: %v", err)
	}

	// Write python file
	pyClientFileName := path.Join(outPath, schema.Name+"Client.py")
	fmt.Printf("Writing to file: %s\n", pyClientFileName)
	err = ioutil.WriteFile(pyClientFileName, []byte(outStr), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
