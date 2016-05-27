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

// Schema represents a JSON Schema.
type Schema struct {
	Name    string   `json:"name,omitempty"`
	Objects []Object `json:"objects,omitempty"`
}

// Object represents a schema object
type Object struct {
	Name       	string               `json:"name,omitempty"`
	Version       	string               `json:"version,omitempty"`
	Type       	string               `json:"type,omitempty"`
	Key        	[]string             `json:"key,omitempty"`
	CfgProperties 	map[string]*Property `json:"cfgProperties,omitempty"`
	OperProperties 	map[string]*Property `json:"operProperties,omitempty"`
	LinkSets   	map[string]*LinkSet  `json:"link-sets,omitempty"`
	Links      	map[string]*Link     `json:"links,omitempty"`
	OperLinkSets   	map[string]*LinkSet  `json:"link-sets,omitempty"`
	OperLinks      	map[string]*Link     `json:"links,omitempty"`
}

// Property represents a property of an object
type Property struct {
	Name        string  `json:"-"`
	Type        string  `json:"type,omitempty"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Items       string  `json:"items,omitempty"`
	Min         float64 `json:"min,omitempty"`
	Max         float64 `json:"max,omitempty"`
	Length      uint    `json:"length,omitempty"`
	Default     string  `json:"default,omitempty"`
	Format      string  `json:"format,omitempty"`
	ShowSummary bool    `json:"showSummary,omitempty"`
}

type LinkSet struct {
	Name string `json:"-"`
	Ref  string `json:"ref,omitempty"`
}

type Link struct {
	Name string `json:"-"`
	Ref  string `json:"ref,omitempty"`
}
