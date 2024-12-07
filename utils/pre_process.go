package utils

import (
	"fmt"
	"os"
	"strings"
)

func GetRunnerByArgs(args []string) Runner {
	var data Runner
	for _, v := range args {
		if strings.Contains(v, "-") {
			// set id for test case
			parts := strings.Split(v, "-")
			if len(parts) == 3 {
				parts[1] = fmt.Sprintf("%s-%s", parts[1], parts[2])
			}
			data.Tests = append(data.Tests, TestRunner{File: parts[0], Id: []interface{}{parts[1]}})
		} else {
			data.Tests = append(data.Tests, TestRunner{File: v, Id: []interface{}{"*"}})
		}
	}
	return data
}

func checkfile(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		PrintRed(fmt.Sprintf("Error: File %s does not exist.\n", file))
		return fmt.Errorf("file %s does not exist", file)
	}
	return nil
}
