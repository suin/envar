package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thcyron/graphs"
	"gopkg.in/yaml.v2"
	"reflect"
)

type Environments []string

type YamlData struct {
	Environments Environments
	Variables    map[string]interface{}
}

type Config struct {
	Environments Environments
	Variables    map[string]map[string]interface{}
}

func (config *Config) SetSameValues(variableName string, value interface{}) {
	values := map[string]interface{}{}
	for _, environmentName := range config.Environments {
		values[environmentName] = value
	}
	config.Variables[variableName] = values
}

func (config *Config) SetValues(variableName string, values map[string]interface{}) {
	config.Variables[variableName] = values
}

type Symbol struct {
	Name string
}

type Value struct {
	Name      string
	Value     interface{}
	DependsOn string
}

func NewValue(name string, value interface{}) *Value {
	dependsOn := ""
	if symbol, ok := value.(Symbol); ok {
		value = nil
		dependsOn = symbol.Name
	}

	return &Value{
		Name:      name,
		Value:     value,
		DependsOn: dependsOn,
	}
}

func (value *Value) HasDependency() bool {
	return value.DependsOn != ""
}

func (value *Value) GetDependency() string {
	return value.DependsOn
}

func makeSymbol(value interface{}) (symbol Symbol, err error) {
	if symbolMap := value.(map[interface{}]interface{}); len(symbolMap) == 1 {
		var symbolName string
		var symbolValue interface{}
		for k, v := range symbolMap {
			symbolName = fmt.Sprintf("%v", k)
			symbolValue = v
			break
		}
		if symbolValue == nil {
			return Symbol{symbolName}, nil
		} else {
			return symbol, errors.New(fmt.Sprintf("invalid environment symbol found"))
		}
	} else {
		return symbol, errors.New(fmt.Sprintf("invalid environment symbol found"))
	}
}

func makeEnvironmentValueMap(environments Environments, values []interface{}) (map[string]*Value, []error) {
	newValues := map[string]*Value{}
	errs := []error{}
	if len(environments) != len(values) {
		errs = append(errs, errors.New(fmt.Sprintf("array length must be %d, but %d", len(environments), len(values))))
		return newValues, errs
	}

	for index, value := range values {
		environmentName := environments[index]
		if isScalar(value) {
			newValues[environmentName] = NewValue(environmentName, value)
		} else if isMap(value) {
			symbol, err := makeSymbol(value)
			if err != nil {
				errs = append(errs, err)
			} else {
				if !environmentExists(environments, symbol.Name) {
					errs = append(errs, errors.New(fmt.Sprintf("no such an environment: {%s}", symbol.Name)))
				} else {
					newValues[environmentName] = NewValue(environmentName, symbol)
				}
			}
		} else {
			data, _ := json.Marshal(value)
			errs = append(errs, errors.New(fmt.Sprintf("variable value must be Bool, Number, String or Array: %s", string(data))))
		}
	}
	return newValues, errs
}

func environmentExists(environments Environments, environmentName interface{}) bool {
	for _, name := range environments {
		if name == environmentName {
			return true
		}
	}
	return false
}

func resolveValues(values map[string]*Value) (map[string]interface{}, error) {
	output := map[string]interface{}{}

	// make graph
	graph := graphs.NewDigraph()
	for _, value := range values {
		if value.HasDependency() {
			graph.AddEdge(value.GetDependency(), value.Name, 0)
		}
	}

	// sort
	order, _, err := graphs.TopologicalSort(graph)
	if err != nil {
		return output, err
	}

	// create output
	for _, value := range values {
		if value.HasDependency() == false {
			output[value.Name] = value.Value
		}
	}
	for e := order.Front(); e != nil; e = e.Next() {
		environmentName := e.Value.(string)
		value := values[environmentName]
		environmentValue := value.Value
		if value.HasDependency() {
			environmentValue = output[value.GetDependency()]
		}
		output[environmentName] = environmentValue
	}
	return output, nil
}

func makeConfig(yamlData YamlData) (*Config, []error) {
	errs := []error{}
	config := &Config{
		Environments: yamlData.Environments,
		Variables:    map[string]map[string]interface{}{},
	}

	for variableName, variableValue := range yamlData.Variables {
		if isScalar(variableValue) {
			config.SetSameValues(variableName, variableValue)
		} else if isSlice(variableValue) {
			environmentValueMap, mapErrors := makeEnvironmentValueMap(config.Environments, variableValue.([]interface{}))
			if len(mapErrors) > 0 {
				for _, err := range mapErrors {
					errs = append(errs, errors.New(fmt.Sprintf("%s: %s", variableName, err)))
				}
			} else {
				resolvedValues, err := resolveValues(environmentValueMap)
				if err != nil {
					errs = append(errs, errors.New(fmt.Sprintf("%s: Cyclic environment symobls are detected", variableName)))
				} else {
					config.SetValues(variableName, resolvedValues)
				}
			}
		} else {
			errs = append(errs, errors.New(fmt.Sprintf("%s value must be type of Number, String, Boolean, null or Array", variableName)))
		}
	}

	return config, errs
}

func isSlice(v interface{}) bool {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Slice:
		return true
	}
	return false
}

func isMap(v interface{}) bool {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Map:
		return true
	}
	return false
}

func isScalar(v interface{}) bool {
	if v == nil {
		return true
	}
	switch reflect.ValueOf(v).Kind() {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		return true
	}
	return false
}

func parse(data string) (config *Config, errs []error) {
	yamlData := YamlData{}
	err := yaml.Unmarshal([]byte(data), &yamlData)
	if err != nil {
		return config, []error{err}
	}
	return makeConfig(yamlData)
}
