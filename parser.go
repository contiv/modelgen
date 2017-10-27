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
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/contiv/modelgen/texthelpers"
)

var validObjects map[string]bool
var propList []*Property

func init() {
	validObjects = make(map[string]bool)

	//Adding standard types as valid objects
	validObjects["string"] = true
	validObjects["int"] = true
	validObjects["number"] = true
	validObjects["bool"] = true
}

// ParseSchema parses the json schema and returns a Schema object
func ParseSchema(input []byte) (*Schema, error) {
	schema := new(Schema)

	// Decode json
	err := json.Unmarshal(input, schema)
	if err != nil {
		log.Errorf("Error decoding json. Err: %v", err)
		return nil, err
	}

	// Collect valid list of objects so they can be included in others
	for _, obj := range schema.Objects {
		if obj.Name != "" {
			validObjects[obj.Name] = true
		}
	}

	// Perform error checking on the schema
	for _, obj := range schema.Objects {
		// check there is a name for the object
		if obj.Name == "" {
			return nil, errors.New("Object has no name")
		}

		// check there is a version for the object
		if obj.Version == "" {
			log.Errorf("object = %s", obj.Name)
			return nil, errors.New("Object has no version")
		}
		if len(obj.CfgProperties) == 0 && len(obj.OperProperties) == 0 {
			log.Errorf("Neither CfgProperties nor OperProperties defined in obj: %#v", obj)
			return nil, errors.New("Neither CfgProperties nor OperProperties defined")
		}

		// Check the type
		if (obj.Type == "") || (obj.Type != "object") {
			return nil, errors.New("Invalid object type")
		}

		// Check each property definition
		for propName, prop := range obj.CfgProperties {
			// set the property name
			prop.Name = propName

			if !isValidProperty(prop) {
				// Add property type to list of objects to be verified
				propList = append(propList, prop)

				prop.CfgObject = true
			}
		}
		for propName, prop := range obj.OperProperties {
			// set the property name
			prop.Name = propName

			if !isValidProperty(prop) {
				// Add property type to list of objects to be verified
				propList = append(propList, prop)

				prop.CfgObject = false
			}
		}

		// Make sure key properties exists
		if len(obj.CfgProperties) > 0 {
			for _, keyField := range obj.Key {
				if obj.CfgProperties[keyField] == nil {
					log.Errorf("Key = %s, object = %s", keyField, obj.Name)
					return nil, errors.New("Key field does not exist in cfg")
				}
			}
		} else if len(obj.OperProperties) > 0 {
			for _, keyField := range obj.Key {
				if obj.OperProperties[keyField] == nil {
					log.Errorf("Key = %s, object = %s", keyField, obj.Name)
					return nil, errors.New("Key field does not exist in oper")
				}
			}
		}

		// parse links and linksets
		for lsName, ls := range obj.LinkSets {
			// set the name
			ls.Name = texthelpers.InitialCap(lsName)

			// FIXME: perform error checking
		}
		for lName, link := range obj.Links {
			// set the name
			link.Name = texthelpers.InitialCap(lName)

			// FIXME: perform error checking
		}
	}
	return schema, nil
}

func isValidProperty(prop *Property) bool {
	// Check the property type
	switch prop.Type {
	case "string":
		break
	case "number":
		break
	case "int":
		break
	case "bool":
		break
	case "array":
		if prop.Items == "" {
			log.Errorf("Array property %s needs items field", prop.Name)
			return false
		}
		// Add property type to list of objects to be verified
		propList = append(propList, prop)
	default:
		return false
	}

	return true
}

// VerifyObjects check if the objects referenced are valid object types
func VerifyObjects(properties []*Property, validObjects map[string]bool) error {
	var propType string

	for _, prop := range properties {
		if prop.Type == "array" {
			propType = prop.Items
		} else {
			propType = prop.Type
		}

		if _, ok := validObjects[propType]; !ok {
			errString := fmt.Sprintf("Unknown object type %s for %s", propType, prop.Name)
			return errors.New(errString)
		}
	}
	return nil
}

// MergeSchema merges two schemas and returns the result
func MergeSchema(first, second *Schema) *Schema {
	if first == nil {
		return second
	}

	// Merge objects from both Schemas
	first.Objects = append(first.Objects, second.Objects...)

	return first
}
