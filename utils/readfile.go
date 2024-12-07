package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

var Directory = ""

func LoadConfigFlowProcess(filePath string, captureValues map[string]interface{}) (*TestConfig, error) {
	if captureValues == nil {
		return LoadConfig(filePath)
	}

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return nil, err
	}
	for key, value := range captureValues {
		err = updateValue(data, key, value)
		if err != nil {
			return nil, err
		}
	}
	var config TestConfig
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dataBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadConfig(filePath string) (*TestConfig, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config TestConfig
	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func ReadRunnerFile(filename string) (Runner, error) {
	var data Runner
	fmt.Println("\n\nReading  file:", filename)
	Directory = filepath.Dir(filename)

	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return data, err
	}
	if data.Mode == "" {
		data.Mode = "unittest" 
	}
	fmt.Println("\n\n\tMode:", data.Mode)
	for i, v := range data.Tests {
		data.Tests[i].File = filepath.Join(Directory, v.File)
	}

	return data, nil
}

func updateValue(data map[string]interface{}, key string, value interface{}) error {
	keys := strings.Split(key, ".")

	// Traverse the map/slice to get to the desired location
	var current interface{} = data
	for i, k := range keys {
		if strings.Contains(k, "[") && strings.Contains(k, "]") {
			// Handle slice index notation
			sliceKey := k[:strings.Index(k, "[")]
			indexStr := k[strings.Index(k, "[")+1 : strings.Index(k, "]")]
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return fmt.Errorf("invalid index in key: %s", k)
			}

			currentSlice, ok := current.(map[string]interface{})[sliceKey].([]interface{})
			if !ok || len(currentSlice) <= index {
				return fmt.Errorf("invalid slice or index out of range for key: %s", k)
			}

			if i == len(keys)-1 {
				// Update the value in the slice
				if mapElem, ok := currentSlice[index].(map[string]interface{}); ok {
					mapElem[keys[i+1]] = value
					return nil
				}
				return fmt.Errorf("element at index %d is not a map", index)
			}
			current = currentSlice[index]

		} else {
			if i == len(keys)-1 {
				// Update the value in the map
				if m, ok := current.(map[string]interface{}); ok {
					m[k] = value
					return nil
				}
				return errors.New("final element is not a map")
			}

			// Navigate deeper into the map
			if m, ok := current.(map[string]interface{}); ok {
				current = m[k]
			} else {
				return errors.New("unexpected type, expected map[string]interface{}")
			}
		}
	}

	return nil
}
