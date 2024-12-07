package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func sendRequest(url, method string, headers map[string]string, reqData map[string]interface{}, useXML bool) (*http.Response, error) {
	client := &http.Client{}
	var body io.Reader

	if len(reqData) > 0 {
		var bodyBytes []byte
		var err error

		if useXML {
			jsonBytes, err := json.Marshal(reqData)
			if err != nil {
				return nil, err
			}
			bodyBytes, err = jsonToXML(jsonBytes)
			if err != nil {
				return nil, err
			}
		} else {
			bodyBytes, err = json.Marshal(reqData)
			if err != nil {
				return nil, err
			}
		}

		body = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}



	return client.Do(req)
}
