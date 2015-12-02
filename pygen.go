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

import "github.com/contiv/objmodel/tools/modelgen/generators"

// GenerateClient generates python client for the model
func (s *Schema) GeneratePythonClient() (string, error) {
	// Generate file headers
	goStr, err := generators.RunTemplate("pyclientHdr", s)
	if err != nil {
		return "", err
	}

	//  Generate clients for all objects
	for _, obj := range s.Objects {
		objStr, err := obj.GeneratePythonClientObjs()
		if err == nil {
			goStr = goStr + objStr
		}
	}

	return goStr, nil
}

// GeneratePythonClientObjs generates the python client for individual objects
func (obj *Object) GeneratePythonClientObjs() (string, error) {
	goStr, err := generators.RunTemplate("pyclientObj", obj)
	if err != nil {
		return "", err
	}

	return goStr, nil
}
