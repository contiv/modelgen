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

	log "github.com/Sirupsen/logrus"
	"github.com/contiv/objmodel/tools/modelgen/generators"
)

var validPropertyTypes = []string{
	"string",
	"bool",
	"array",
	"number",
	"int",
}

// GenerateGo generates go code for the schema
func (s *Schema) GenerateGo() (string, error) {
	// Generate file headers
	outStr, err := s.GenerateGoHdrs()
	if err != nil {
		return "", err
	}

	// Generate structs
	structStr, err := s.GenerateGoStructs()
	if err != nil {
		log.Errorf("Error generating go structs. Err: %v", err)
		return "", err
	}

	// Merge the header and struct
	outStr += structStr

	// Merge rest handler
	str, err := s.GenerateGoFuncs()
	if err != nil {
		return "", err
	}

	gobytes, err := format.Source([]byte(outStr + str))
	if err != nil {
		return outStr + str, err
	}

	return string(gobytes), nil
}

// GenerateGoStructs generates go code from a schema
func (s *Schema) GenerateGoStructs() (string, error) {
	var goStr string

	//  Generate all object definitions
	for _, obj := range s.Objects {
		objStr, err := obj.GenerateGoStructs()
		if err != nil {
			return "", err
		}

		goStr += objStr
	}

	for _, name := range []string{"gostructs", "callbacks", "init", "register"} {
		str, err := generators.RunTemplate(name, s)
		if err != nil {
			return "", err
		}

		goStr += str
	}

	return goStr, nil
}

// GenerateGoHdrs generates go file headers
func (s *Schema) GenerateGoHdrs() (string, error) {
	return generators.RunTemplate("hdr", s)
}

func (s *Schema) GenerateGoFuncs() (string, error) {
	// Output the functions and routes
	return generators.RunTemplate("routeFunc", s)
}

func (obj *Object) GenerateGoStructs() (string, error) {
	return generators.RunTemplate("objstruct", obj)
}

func (prop *Property) GenerateGoStructs() (string, error) {
	var found bool

	for _, myType := range validPropertyTypes {
		if myType == prop.Type {
			found = true
		}
	}

	if !found {
		return "", errors.New("Unknown Property")
	}

	return generators.RunTemplate("propstruct", prop)
}
