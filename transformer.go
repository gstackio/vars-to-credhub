package main

import (
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v2"
)

// Importable represents a credential to be loaded into Credhub via `bulk-import`
type Importable struct {
	Name  string      `yaml:"name"`
	Type  string      `yaml:"type"`
	Value interface{} `yaml:"value,omitempty"`
}

// BulkImport represents what will be actually sent to Credhub
type BulkImport struct {
	Credentials []Importable `yaml:"credentials"`
}

// Utility to test if something is a map
func isMap(x interface{}) bool {
	t := fmt.Sprintf("%T", x)
	return strings.HasPrefix(t, "map[")
}

// Utility to test if something is a an array
func isArray(x interface{}) bool {
	t := fmt.Sprintf("%T", x)
	return strings.HasPrefix(t, "[]inter")
}

func handleMapSlice(mapVal yaml.MapSlice, prefix string, parentKey string) Importable {
	idxToDelete := -1
	for idx, item := range mapVal {
		if item.Value == "public_key_fingerprint" {
			idxToDelete = idx
			break
		}
	}
	value := mapVal
	if idxToDelete >= 0 {
		value = append(mapVal[:idxToDelete], mapVal[idxToDelete+1:]...)
	}
	return Importable{
		Name:  fmt.Sprintf("%s/%s", prefix, parentKey),
		Type:  getType(parentKey, fmt.Sprint(mapVal)),
		Value: value,
	}
}

func handleMap(mapVal map[interface{}]interface{}, prefix string, parentKey string) Importable {
	//RSA key types if output from bosh will have a fingerprint that credhub
	//can't deal with, so we'll just remove it. Delete is harmless if the
	//key DNE - so just do it on every map.
	delete(mapVal, "public_key_fingerprint")
	return Importable{
		Name:  fmt.Sprintf("%s/%s", prefix, parentKey),
		Type:  getType(parentKey, fmt.Sprint(mapVal)),
		Value: mapVal,
	}
}

func handleArray(arrayVal []interface{}, prefix string, key string) Importable {
	return Importable{
		Name:  fmt.Sprintf("%s/%s", prefix, key),
		Type:  "json",
		Value: arrayVal,
	}
}

func getType(key string, valStr string) string {
	//attempt to guess the type for the item, default to "value"
	if strings.Contains(key, "password") || strings.Contains(key, "secret") {
		return "password"
		//the cert check should be above rsa since a cert may also contain a pk
	} else if strings.Contains(valStr, "CERTIFICATE---") {
		return "certificate"
	} else if strings.Contains(valStr, "KEY---") {
		return "rsa"
	}
	return "value"
}

func makeImportable(prefix string, key string, valStr string) Importable {
	return Importable{
		Name:  fmt.Sprintf("%s/%s", prefix, key),
		Type:  getType(key, valStr),
		Value: valStr,
	}
}

// Transform performs the translation from vars file to import file
func Transform(prefix string, input io.Reader) (BulkImport, error) {
	decoder := yaml.NewDecoder(input)

	var pipelineVars yaml.MapSlice
	err := decoder.Decode(&pipelineVars)
	if err != nil {
		return BulkImport{
			Credentials: nil,
		}, err
	}

	vals := make([]Importable, 0, len(pipelineVars))
	for _, mapItem := range pipelineVars {
		key := mapItem.Key
		val := mapItem.Value
		// Let's require, for now, simple types & maps in the var field
		switch valType := val.(type) {
		case yaml.MapSlice:
			vals = append(vals, handleMapSlice(val.(yaml.MapSlice), prefix, key.(string)))
		default:
			if isMap(val) {
				vals = append(vals, handleMap(val.(map[interface{}]interface{}), prefix, key.(string)))
			} else if isArray(val) {
				vals = append(vals, handleArray(val.([]interface{}), prefix, key.(string)))
			} else {
				return BulkImport{}, fmt.Errorf("Invalid value type in vars file %T. Currently only primitive values & maps are supported", valType)
			}
		case bool, float32, float64, int, int16, int32, int64, string, uint, uint16, uint32, uint64, nil:
			valStr := fmt.Sprint(val)
			vals = append(vals, makeImportable(prefix, key.(string), valStr))
		}
	}

	return BulkImport{
		Credentials: vals,
	}, nil
}
