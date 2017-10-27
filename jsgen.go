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
	"text/template"

	"github.com/contiv/modelgen/generators"
	"github.com/contiv/modelgen/texthelpers"
)

var funcMap = template.FuncMap{
	"initialCap":        texthelpers.InitialCap,
	"initialLow":        texthelpers.InitialLow,
	"depunct":           texthelpers.Depunct,
	"capFirst":          texthelpers.CapFirst,
	"translateCfgType":  texthelpers.TranslateCfgPropertyType,
	"translateOperType": texthelpers.TranslateOperPropertyType,
}

func (s *Schema) GenerateJs() ([]byte, error) {
	goBytes := []byte(`//
// This file is auto generated by modelgen tool
// Do not edit this file manually
`)

	//  Generate all object Views
	for _, obj := range s.Objects {
		objBytes, err := obj.GenerateJsViews()
		if err != nil {
			return nil, err
		}

		goBytes = append(goBytes, objBytes...)
	}

	return goBytes, nil
}

func (obj *Object) GenerateJsViews() ([]byte, error) {
	goBytes, err := generators.RunTemplate("jsView", obj)
	if err != nil {
		return nil, err
	}

	return goBytes, nil
}
