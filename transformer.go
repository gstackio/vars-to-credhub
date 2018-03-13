package main

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

// Importable represents a credential to be loaded into Credhub via `bulk-import`
type Importable struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

// BulkImport represents what will be actually sent to Credhub
type BulkImport struct {
	Credentials []Importable `yaml:"credentials"`
}

// Transform performs the translation from vars file to import file
func Transform(prefix string, input io.Reader) (BulkImport, error) {
	decoder := yaml.NewDecoder(input)

	var pipelineVars map[interface{}]interface{}
	decoder.Decode(&pipelineVars)

	vals := make([]Importable, 0, len(pipelineVars))
	for key, val := range pipelineVars {

		// Let's require, for now, simple types only in the var field
		var valStr string
		switch v := val.(type) {
		default:
			return BulkImport{}, fmt.Errorf("Invalid value type in vars file %T. Currently only primitive values are supported", v)
		case bool, float32, float64, int, int16, int32, int64, string, uint, uint16, uint32, uint64:
			valStr = fmt.Sprint(val)
		}

		vals = append(vals, Importable{
			Name:  fmt.Sprintf("%s/%s", prefix, key),
			Type:  "value",
			Value: valStr,
		})
	}

	return BulkImport{
		Credentials: vals,
	}, nil
}
