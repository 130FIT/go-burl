package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func checkResponse(bodyBytes []byte, response *http.Response, expectedStatus int, assertions map[string]interface{}, captures []Capture, capturedValues map[string]interface{}, useXML bool, testStep *TestStep, isSkipCapture bool) error {

	if len(bodyBytes) == 0 {
		testStep.Error = "response body is empty"
		return fmt.Errorf("response body is empty")
	}

	var jsonBytes []byte
	var err error
	if useXML {
		jsonBytes, err = xmlToJSON(bodyBytes)
		if err != nil {
			testStep.Error = "error converting XML to JSON"
			testStep.TestStatus = "FAIL"
			return fmt.Errorf("error converting XML to JSON: %w", err)
		}
	} else {
		jsonBytes = bodyBytes
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &responseData); err != nil {
		testStep.Error = "error unmarshaling JSON response"
		testStep.TestStatus = "FAIL"
		return fmt.Errorf("error unmarshaling JSON response: %w", err)
	}
	testStep.ResponseBody = responseData
	fmt.Printf("\n\nResponse status: %d\nAssert status: %d\nResult: %v\n\n", response.StatusCode, expectedStatus, response.StatusCode == expectedStatus)
	if response.StatusCode != expectedStatus {
		testStep.Error = fmt.Sprintf("expected status %d, got %d", expectedStatus, response.StatusCode)
		testStep.TestStatus = "FAIL"
		return fmt.Errorf("expected status %d, got %d", expectedStatus, response.StatusCode)
	}
	isAssert := true
	for _key, expectedValue := range assertions {
		key, mode := assertActualTool(_key)
		actualValue, exists := getValue(responseData, key)
		if !exists {
			testStep.Error = fmt.Sprintf("key '%s' not found in response", key)
			testStep.TestStatus = "FAIL"
			testStep.AssertResponse = append(testStep.AssertResponse, map[string]interface{}{"key": key, "expectedValue": expectedValue, "actualValue": nil, "isAssert": false, "mode": mode})
			PrintRed(fmt.Sprintf("Assert Failed\n \n\t Key: %s \n\t Expected: %v \n\t Actual: %v\n\n", key, expectedValue, nil))
			return fmt.Errorf("key '%s' not found in response", key)
		}
		compareMode := mode
		if strings.Contains(mode, "(count)") {
			actualValue = len(actualValue.([]interface{}))
			compareModes := strings.Split(mode, "|")
			compareMode = compareModes[len(compareModes)-1]
		}
		isAssert = assertCompare(actualValue, expectedValue, compareMode)

		testStep.AssertResponse = append(testStep.AssertResponse, map[string]interface{}{"key": key, "expectedValue": expectedValue, "actualValue": actualValue, "isAssert": isAssert, "mode": mode})
		if !isAssert {
			PrintRed(fmt.Sprintf("\nResult : %v\n", isAssert))
			errorStr := fmt.Sprintf("\nAssert\n\nkey :%v\n \n\t Expected: (%T) %v \n\t Actual: (%T) %v \n\t Compare with %v\n", key, expectedValue, expectedValue, actualValue, actualValue, mode)
			testStep.Error = errorStr
			testStep.TestStatus = "FAIL"
			return fmt.Errorf(errorStr)
		}
	}
	PrintGreen(fmt.Sprintf("\nResult : %v\n", isAssert))
	if !isSkipCapture {
		for _, capture := range captures {
			value, exists := getValue(responseData, capture.Path)
			if !exists {
				testStep.Error = fmt.Sprintf("key %s not found in response", capture.Path)
				PrintRed(fmt.Sprintf("\nkey %s not found in response", capture.Path))
			} else {
				capturedValues[capture.Name] = value
			}
		}
	}
	testStep.TestStatus = "PASS"
	return nil
}
func assertCompare(actualValue, expectedValue interface{}, mode string) bool {
	fmt.Printf("\n\nAssert\n \n\t Expected: (%T) %v \n\t Actual: (%T) %v \n\t Compare with %v\n", expectedValue, expectedValue, actualValue, actualValue, mode)
	switch mode {
	case "(<)", "(>)", "(<=)", "(>=)":
		return compare(actualValue, expectedValue, mode)
	case "(!=)":
		return actualValue != expectedValue
	case "(==)":
		return actualValue == expectedValue
	case "(contains)":
		return strings.Contains(actualValue.(string), expectedValue.(string))
	case "(notcontains)":
		return !strings.Contains(actualValue.(string), expectedValue.(string))
	default:
		return actualValue == expectedValue
	}
}

func compare(actual, expected interface{}, mode string) bool {
	// Helper function to convert a value to float64 if possible

	var actualFloat, expectedFloat float64
	var err error
	if actualFloat, err = ToFloat64(actual); err != nil {
		return false
	}
	if expectedFloat, err = ToFloat64(expected); err != nil {
		return false
	}
	// Perform comparison based on mode
	switch mode {
	case "(<)":
		return actualFloat < expectedFloat
	case "(>)":
		return actualFloat > expectedFloat
	case "(<=)":
		return actualFloat <= expectedFloat
	case "(>=)":
		return actualFloat >= expectedFloat
	default:
		fmt.Printf("unsupported mode: %s\n", mode)
		return false
	}
}
func assertActualTool(key string) (string, string) {
	var assertKey string
	var mode string
	switch {
	case strings.HasPrefix(key, "(count)"):
		mode = "(count)"
	case strings.HasPrefix(key, "(<)"):
		mode = "(<)"
	case strings.HasPrefix(key, "(>)"):
		mode = "(>)"
	case strings.HasPrefix(key, "(<=)"):
		mode = "(<=)"
	case strings.HasPrefix(key, "(>=)"):
		mode = "(>=)"
	case strings.HasPrefix(key, "(!=)"):
		mode = "(!=)"
	case strings.HasPrefix(key, "(==)"):
		mode = "(==)"
	case strings.HasPrefix(key, "(contains)"):
		mode = "(contains)"
	case strings.HasPrefix(key, "(notcontains)"):
		mode = "(notcontains)"
	default:
		mode = "(==)"
	}
	assertKey = strings.Replace(key, mode, "", -1)
	assertKey = strings.TrimSpace(assertKey)
	if mode == "(count)" {
		newKey, nextMode := assertActualTool(assertKey)
		assertKey = newKey
		mode = mode + "|" + nextMode
	}
	return assertKey, mode
}

func getValue(data map[string]interface{}, key string) (interface{}, bool) {
	keys := strings.Split(key, ".")
	return getNestedValue(data, keys)
}

// Helper function to handle nested keys
func getNestedValue(data interface{}, keys []string) (interface{}, bool) {
	if len(keys) == 0 {
		return data, true
	}

	currentKey := keys[0]
	remainingKeys := keys[1:]

	switch data.(type) {
	case map[string]interface{}:
		dataMap := data.(map[string]interface{})
		if strings.Contains(currentKey, "[") && strings.Contains(currentKey, "]") {
			// Extract the map key and index
			parts := strings.Split(currentKey, "[")
			mapKey := parts[0]
			indexStr := strings.Trim(parts[1], "]")
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil, false
			}

			value, ok := dataMap[mapKey].([]interface{})
			if !ok || index >= len(value) {
				return nil, false
			}
			return getNestedValue(value[index], remainingKeys)
		} else {
			// Get the value of the map key
			value, ok := dataMap[currentKey]
			if !ok {
				return nil, false
			}
			return getNestedValue(value, remainingKeys)
		}

	case []interface{}:
		return nil, false

	default:
		return nil, false
	}
}
