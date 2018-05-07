package main

import (
    "fmt"
    "io"
    "strings"

    "gopkg.in/yaml.v2"
)

// Importable represents a credential to be loaded into Credhub via `bulk-import`
type Importable struct {
    Name   string `yaml:"name"`
    Type   string `yaml:"type"`
    Value  string `yaml:"value,omitempty"`
    SubMap map[interface{}]interface{} `yaml:"subMapValue,omitempty"`
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

func handleMap(mapVal map[interface{}]interface{}, prefix string, parentKey string) Importable {
    return Importable{
        Name: fmt.Sprintf("%s/%s", prefix, parentKey),
        Type: getType(parentKey, fmt.Sprint(mapVal)),
        SubMap: mapVal,
    }
}

func getType(key string, valStr string) string {
    //attempt to guess the type for the item, default to "value"
    if strings.Contains(key, "password") || strings.Contains(key, "secret") {
        return "password"
    //the cert check should be above rsa since a cert may also contain a pk
    } else if strings.Contains(valStr, "CERTIFICATE") {
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

    var pipelineVars map[interface{}]interface{}
    decoder.Decode(&pipelineVars)

    vals := make([]Importable, 0, len(pipelineVars))
    for key, val := range pipelineVars {

        // Let's require, for now, simple types only in the var field
        var valStr string
        switch v := val.(type) {
        default:
            if isMap(val) {
                vals = append(vals, handleMap(val.(map[interface{}]interface{}), prefix, key.(string)))
            } else {
                return BulkImport{}, fmt.Errorf("Invalid value type in vars file %T. Currently only primitive values are supported", v)
            }
        case bool, float32, float64, int, int16, int32, int64, string, uint, uint16, uint32, uint64:
            valStr = fmt.Sprint(val)
            vals = append(vals, makeImportable(prefix, key.(string), valStr))
        }
    }

    return BulkImport{
        Credentials: vals,
    }, nil
}
