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
	"errors"
	"go/format"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/contiv/modelgen/generators"
)

var validPropertyTypes = []string{
	"string",
	"bool",
	"array",
	"number",
	"int",
}

// GenerateGo generates go code for the schema
func (s *Schema) GenerateGo() ([]byte, error) {
	// Generate file headers
	out, err := s.GenerateGoHdrs()
	if err != nil {
		return nil, err
	}

	// Generate structs
	str, err := s.GenerateGoStructs()
	if err != nil {
		log.Errorf("Error generating go structs. Err: %v", err)
		return nil, err
	}

	// Merge the header and struct
	out = append(out, str...)

	// Merge rest handler
	str, err = s.GenerateGoFuncs()
	if err != nil {
		return nil, err
	}

	out = append(out, str...)

	gobytes, err := format.Source(out)
	if err != nil {
		os.Stderr.Write(out)
		return out, err
	}

	return gobytes, nil
}

// GenerateGoStructs generates go code from a schema
func (s *Schema) GenerateGoStructs() ([]byte, error) {
	var goBytes []byte

	//  Generate all object definitions
	for _, obj := range s.Objects {
		objBytes, err := obj.GenerateGoStructs()
		if err != nil {
			return nil, err
		}

		goBytes = append(goBytes, objBytes...)
	}

	for _, name := range []string{"gostructs", "callbacks", "init", "register"} {
		byts, err := generators.RunTemplate(name, s)
		if err != nil {
			return nil, err
		}

		goBytes = append(goBytes, byts...)
	}

	return goBytes, nil
}

// GenerateGoHdrs generates go file headers
func (s *Schema) GenerateGoHdrs() ([]byte, error) {
	return generators.RunTemplate("hdr", s)
}

func (s *Schema) GenerateGoFuncs() ([]byte, error) {
	// Output the functions and routes
	return generators.RunTemplate("routeFunc", s)
}

func (obj *Object) GenerateGoStructs() ([]byte, error) {
	return generators.RunTemplate("objstruct", obj)
}

func (prop *Property) GenerateGoStructs() (string, error) {
	// this function has to return a string because it is used in templates.
	var found bool

	for _, myType := range validPropertyTypes {
		if myType == prop.Type {
			found = true
		}
	}

	if !found {
		return "", errors.New("Unknown Property")
	}
	byt, err := generators.RunTemplate("propstruct", prop)
	return string(byt), err
}
