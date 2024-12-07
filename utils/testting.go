package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Testting struct {
	CaptureValues map[string]interface{}
	SumaryReport  Sumary
	Report        TestReport
	Step          int
	IsForce       bool
	StartTime     time.Time
}

func (t *Testting) init() {
	t.CaptureValues = make(map[string]interface{})
	t.SumaryReport = Sumary{}
	t.Report = TestReport{}
	t.Step = 1
	t.StartTime = time.Now()
}
func (t *Testting) RunnerFile(runnerFile string) {
	fmt.Println("\nRunner files mode")
	runner, err := ReadRunnerFile(runnerFile)
	if err != nil {
		PrintRed(fmt.Sprintf("\nError: %s", err))
		return
	}
	switch runner.Mode {
	case "flow":
		t.FlowMode(runner.Tests)
	default:
		t.UnitestMode(runner.Tests)
	}
}
func (t *Testting) UnitestMode(tests []TestRunner) {
	fmt.Println("\n๐ Unitest mode")
	t.init()
	t.Report.InitReport()
	t.Report.Title = "Unitest Report"
	t.StartTime = time.Now()
	for _, test := range tests {
		if test.DisplayName != "" {
			fmt.Printf("\n\n%s\n\n", test.DisplayName)
		}
		config, err := LoadConfig(test.File)
		if err != nil {
			PrintRed(fmt.Sprintf("\nError: %s", err))
			return
		}

		newConfig := SelectTestCase(config, test.Id)
		inCaptureValues := make(map[string]interface{})
		err, _ = t.Test(newConfig, inCaptureValues, test.File, TestOptions{InFlow: true})
		if err != nil {
			PrintRed(fmt.Sprintf("\nError: %s", err))
			return
		}
	}
	t.SumaryReport.SetValueWithReport(t.Report)
	t.SumaryReport.PrintSumary()
	elapsed := time.Since(t.StartTime).Seconds()
	fmt.Printf("\n\nTesting time: %f seconds \n\n", elapsed)
	t.Report.SaveReportToFile(elapsed)
}
func (t *Testting) FlowMode(tests []TestRunner) {
	fmt.Println("\n๐ Flow mode")
	t.init()
	t.Report.InitReport()
	t.Report.Title = "Flow Report"
	t.StartTime = time.Now()
	for _, test := range tests {
		if test.DisplayName != "" {
			fmt.Printf("\n\n%s\n\n", test.DisplayName)
		}
		config, err := LoadConfigFlowProcess(test.File, t.CaptureValues)
		if err != nil {
			PrintRed(fmt.Sprintf("\nError: %s", err))
			return
		}
		newConfig := SelectTestCase(config, test.Id)
		inCaptureValues := make(map[string]interface{})
		err, inCaptureValues = t.Test(newConfig, inCaptureValues, test.File, TestOptions{InFlow: true})
		if err != nil {
			PrintRed(fmt.Sprintf("\nError: %s", err))
			return
		}
		for k, v := range inCaptureValues {
			t.CaptureValues[k] = v
		}
	}
	t.SumaryReport.SetValueWithReport(t.Report)
	t.SumaryReport.PrintSumary()
	elapsed := time.Since(t.StartTime).Seconds()
	fmt.Printf("\n\nTesting time: %f seconds \n\n", elapsed)
	t.Report.SaveReportToFile(elapsed)

}

type TestOptions struct {
	InFlow bool
}

func (t *Testting) Test(testConfig *TestConfig, inCaptureValues map[string]interface{}, file string, options ...TestOptions) (error, map[string]interface{}) {
	inFlow := false
	if len(options) > 0 {
		inFlow = options[0].InFlow
	}
	subStep := 1
	for _, testcase := range testConfig.Cases {
		step := fmt.Sprintf("%d", t.Step)
		if len(testConfig.Cases) > 1 {
			step = fmt.Sprintf("%d.%d", t.Step, subStep)
		}
		fmt.Printf("\n\n--------------------------------------------------\n\n")
		fmt.Printf("\n\nStep:%s\n\nTest case %d: %s\n", step, testcase.ID, testcase.Name)
		if testcase.Description != "" {
			fmt.Printf("\nDescription: %s\n", testcase.Description)
		}
		testStepReport := TestStep{}
		testStepReport.Init(*testConfig, testcase, step, file)
		if t.IsForce {
			PrintRed("\n\nSkip this test case\n")
			testStepReport.TestStatus = "SKIP"
			subStep++
			t.Report.TestStep = append(t.Report.TestStep, testStepReport)
			continue
		}

		resp, err := sendRequest(testConfig.URL, testConfig.Method, testConfig.Headers, testConfig.GetRequest(testcase), testConfig.TypeAPI == "xml")
		if err != nil {
			AddErrorReport(fmt.Sprintf("Test case %d failed: %v", testcase.ID, err))
			if inFlow {
				t.IsForce = true
			}
			t.Report.TestStep = append(t.Report.TestStep, testStepReport)
			continue
		}
		testStepReport.saveRequest(testConfig, testcase, resp, step)
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			AddErrorReport(fmt.Sprintf("Test case %d failed: %v", testcase.ID, err))
			if inFlow {
				t.IsForce = true
			}
			t.Report.TestStep = append(t.Report.TestStep, testStepReport)
			continue
		}
		resp.Body.Close()
		SaveResponseToFile(bodyBytes, GetFileExtensionFromContentType(resp.Header.Get("Content-Type")), step, &testStepReport)
		err = checkResponse(bodyBytes, resp, testcase.Status, testcase.AssertResponse, testConfig.Captures, inCaptureValues, testConfig.TypeAPI == "xml", &testStepReport, testcase.SkipCapture)
		if inFlow && err != nil {
			testStepReport.TestStatus = "FAIL"
			t.IsForce = true
		}
		subStep++
		t.Report.TestStep = append(t.Report.TestStep, testStepReport)
	}
	t.Step++
	return nil, inCaptureValues
}

func (t *TestStep) saveRequest(testConfig *TestConfig, testcase TestCase, resp *http.Response, step string) {
	saveConfigFile := testConfig.Copy()
	saveConfigFile.Cases = []TestCase{testcase}

	var configBuffer bytes.Buffer
	encoder := json.NewEncoder(&configBuffer)
	encoder.SetEscapeHTML(false) // ปิดการ escape HTML
	if err := encoder.Encode(saveConfigFile); err != nil {
		fmt.Printf("error encoding JSON: %v\n", err)
		return
	}
	SaveRequestToFile(configBuffer.Bytes(), ".json", step, t)
	t.RequestBody = saveConfigFile.GetRequest(testcase)
	t.RequestHeader = resp.Request.Header.Clone()
	t.ResponseHeader = resp.Header.Clone()
	t.ResponseCode = resp.StatusCode
	contentType := resp.Header.Get("Content-Type")
	fileType := GetFileExtensionFromContentType(contentType)
	t.ResponseType = strings.Replace(fileType, ".", "", -1)
}
