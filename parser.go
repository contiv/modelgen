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

	log "github.com/Sirupsen/logrus"
	"github.com/contiv/modelgen/texthelpers"
)

// ParseSchema parses the json schema and returns a Schema object
func ParseSchema(input []byte) (*Schema, error) {
	schema := new(Schema)

	// Decode json
	err := json.Unmarshal(input, schema)
	if err != nil {
		log.Errorf("Error decoding json. Err: %v", err)
		return nil, err
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
			return nil, errors.New("Niether CfgProperties nor OperProperties defined")
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
				return nil, errors.New("Invalid property type")
			}
		}
		for propName, prop := range obj.OperProperties {
			// set the property name
			prop.Name = propName

			if !isValidProperty(prop) {
				return nil, errors.New("Invalid property type")
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
	default:
		log.Errorf("Unknown proprty type %s for %s", prop.Type, prop.Name)
		return false
	}

	return true
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
