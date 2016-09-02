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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/contiv/modelgen/generators"
)

// Simple test to parse json schema
func TestParseJsonSchema(t *testing.T) {
	if err := generators.ParseTemplates(); err != nil {
		t.Fatal(err)
	}

	dir, err := os.Open("testdata")
	if err != nil {
		t.Fatal(err)
	}

	dirnames, err := dir.Readdirnames(0)
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range dirnames {
		t.Logf("Parsing suite %s", name)
		var schema *Schema
		basepath := filepath.Join("testdata", name)

		// Parse all files in input directory
		err := filepath.Walk(basepath, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Ignore non-json files and ignore files in subdirs
			if filepath.Ext(path) != ".json" {
				return nil
			}

			t.Logf("Parsing file: %q\n", path)

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
			t.Fatal(err)
		}

        if schema == nil {
            t.Fatal("Could not find schema, aborting.")
        }

		// Generate the code
		goStr, err := schema.GenerateGo()
		if err != nil {
			t.Fatalf("Error generating go code. Err: %v", err)
		}

		output, err := ioutil.ReadFile(filepath.Join(basepath, "output.go"))
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(goStr, output) {
			if err := ioutil.WriteFile("/tmp/generated.go", []byte(goStr), 0666); err != nil {
				t.Fatalf("Err writing debug file `/tmp/generated.go` during failed test: %v", err)
			}
			fmt.Println("Generated code written to /tmp/generated.go")
			t.Fatal("Generated string from input was not equal to output string")
		}
	}
}
