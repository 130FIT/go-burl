package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var TimeSave = time.Now().Format("20060102-150405")

var DirPath, _ = filepath.Abs(fmt.Sprintf("reports/report-%s", TimeSave))

func SaveResponseToFile(bodyBytes []byte, fileExtension string, step string, t *TestStep) error {
	dirName := fmt.Sprintf("reports/report-%s/response", TimeSave)
	pwd, _ := os.Getwd()
	dirName = filepath.Join(pwd, dirName)
	nameFile := filepath.Join(dirName, fmt.Sprintf("response-test-step-%s%s", step, fileExtension))

	if err := os.MkdirAll(dirName, 0755); err != nil {
		return fmt.Errorf("error creating directories: %w", err)
	}

	file, err := os.Create(nameFile)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// Write the bodyBytes ensuring UTF-8 encoding
	_, err = file.Write(bodyBytes)
	if err != nil {
		return fmt.Errorf("error writing response to file: %w", err)
	}
	rootFile, _ := filepath.Abs(nameFile)
	rootFile = strings.ReplaceAll(rootFile, "\\", "/")
	filePath := fmt.Sprintf("file:///%s", rootFile)
	t.ResponseSource = filePath
	return nil
}

func SaveRequestToFile(bodyBytes []byte, fileExtension string, step string, t *TestStep) error {

	dirName := fmt.Sprintf("reports/report-%s/request", TimeSave)
	pwd, _ := os.Getwd()
	dirName = filepath.Join(pwd, dirName)
	nameFile := filepath.Join(dirName, fmt.Sprintf("request-test-step-%s%s", step, fileExtension))
	if err := os.MkdirAll(dirName, 0755); err != nil {
		return fmt.Errorf("error creating directories: %w", err)
	}

	file, err := os.Create(nameFile)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// Write the bodyBytes ensuring UTF-8 encoding
	_, err = file.Write(bodyBytes)
	if err != nil {
		return fmt.Errorf("error writing request to file: %w", err)
	}

	rootFile, _ := filepath.Abs(nameFile)
	rootFile = strings.ReplaceAll(rootFile, "\\", "/")
	filePath := fmt.Sprintf("file:///%s", rootFile)
	t.Sources = filePath
	return nil
}

func GetFileExtensionFromContentType(contentType string) string {
	switch {
	case strings.Contains(contentType, "json"):
		return ".json"
	case strings.Contains(contentType, "xml"):
		return ".xml"
	case strings.Contains(contentType, "html"):
		return ".html"
	case strings.Contains(contentType, "plain"):
		return ".txt"
	default:
		return ".bin"
	}
}

func replaceSpecialCharacters(str string) string {
	specialCharacters := map[string]string{"\\u003c": "<", "\\u003e": ">"}
	for key, value := range specialCharacters {
		str = strings.ReplaceAll(str, key, value)
	}
	return str
}
