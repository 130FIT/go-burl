package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	xj "github.com/basgys/goxml2json"
)

// jsonToXML converts JSON data to XML format
func jsonToXML(jsonData []byte) ([]byte, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	xmlData, err := mapToXML(jsonMap)
	if err != nil {
		return nil, fmt.Errorf("error converting map to XML: %w", err)
	}

	return xmlData, nil
}

// mapToXML converts a map to XML format
func mapToXML(data map[string]interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	buffer.WriteString(xml.Header)
	buffer.WriteString("<root>")

	err := writeXML(buffer, data)
	if err != nil {
		return nil, err
	}

	buffer.WriteString("</root>")
	return buffer.Bytes(), nil
}

// writeXML writes a map to an XML writer
func writeXML(buffer *bytes.Buffer, data map[string]interface{}) error {
	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			buffer.WriteString(fmt.Sprintf("<%s>", key))
			err := writeXML(buffer, v)
			if err != nil {
				return err
			}
			buffer.WriteString(fmt.Sprintf("</%s>", key))
		case []interface{}:
			for _, item := range v {
				buffer.WriteString(fmt.Sprintf("<%s>", key))
				if m, ok := item.(map[string]interface{}); ok {
					err := writeXML(buffer, m)
					if err != nil {
						return err
					}
				} else {
					buffer.WriteString(fmt.Sprintf("%v", item))
				}
				buffer.WriteString(fmt.Sprintf("</%s>", key))
			}
		default:
			buffer.WriteString(fmt.Sprintf("<%s>%v</%s>", key, v, key))
		}
	}
	return nil
}

// xmlToJSON converts XML data to JSON format
func xmlToJSON(xmlData []byte) ([]byte, error) {
	xml := strings.NewReader(string(xmlData))
	jsonData, err := xj.Convert(xml)
	if err != nil {
		return nil, fmt.Errorf("error converting XML to JSON: %w", err)
	}
	return jsonData.Bytes(), nil
}

func XmlStrToJson(xmlStr string) string {
	jsonData, err := xj.Convert(strings.NewReader(xmlStr))
	if err != nil {
		fmt.Println("Error converting XML to JSON: ", err)
		return ""
	}
	return jsonData.String()
}
