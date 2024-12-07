package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func IndexOf(array []interface{}, value interface{}) int {
	for i, v := range array {
		if v == value {
			return i
		}
	}
	return -1
}

func ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("error parsing string to float64: %w", err)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("unsupported type for conversion to float64: %T", v)
	}
}

func ToStr(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func SelectTestCase(config *TestConfig, ids []interface{}) *TestConfig {
	newConfig := config.Copy()
	newConfig.Cases = []TestCase{}
	if len(ids) == 0 {
		return config
	}
	for _, id := range ids {
		id := ToStr(id)
		if id == "*" {
			return config
		}
		rangeId := strings.Split(id, "-")
		if len(rangeId) == 2 {
			startIdInt, _ := strconv.Atoi(rangeId[0])
			endIdInt, _ := strconv.Atoi(rangeId[1])
			for i := startIdInt; i <= endIdInt; i++ {
				newConfig = selectTestCaseById(newConfig, config, i)
			}
		} else {
			idInt, _ := strconv.Atoi(ToStr(id))
			newConfig = selectTestCaseById(newConfig, config, idInt)
		}
	}
	return newConfig
}
func selectTestCaseById(newConfig *TestConfig, config *TestConfig, idInt int) *TestConfig {
	for _, test := range config.Cases {
		if test.ID == idInt {
			newConfig.Cases = append(newConfig.Cases, test)
		}
	}
	return newConfig
}
func MakeRunner(files []string) Runner {
	runner := Runner{}
	for _, file := range files {
		test := TestRunner{}
		test.File = file
		test.Id = []interface{}{}
		runner.Tests = append(runner.Tests, test)
	}
	return runner
}
