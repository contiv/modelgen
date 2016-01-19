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

import "github.com/contiv/modelgen/generators"

// GenerateClient generates go client for the model
func (s *Schema) GenerateClient() ([]byte, error) {
	// Generate file headers
	goBytes, err := generators.RunTemplate("clientHdr", s)
	if err != nil {
		return nil, err
	}

	//  Generate all object definitions
	for _, obj := range s.Objects {
		objBytes, err := obj.GenerateGoStructs()
		if err != nil {
			return nil, err
		}

		goBytes = append(goBytes, objBytes...)
	}

	//  Generate clients for all objects
	for _, obj := range s.Objects {
		objBytes, err := obj.GenerateClientObjs()
		if err != nil {
			return nil, err
		}

		goBytes = append(goBytes, objBytes...)
	}

	return goBytes, nil
}

// GenerateClientObjs generates the client for individual objects
func (obj *Object) GenerateClientObjs() ([]byte, error) {
	goBytes, err := generators.RunTemplate("clientObj", obj)
	if err != nil {
		return nil, err
	}

	return goBytes, nil
}
