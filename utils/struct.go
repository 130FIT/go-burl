package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Capture struct {
	Name string `json:"pass_path"`
	Path string `json:"capture_path"`
}

type TestCase struct {
	ID             int                    `json:"id"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Tags           []string               `json:"tags"`
	Status         int                    `json:"status"`
	SkipCapture    bool                   `json:"skip_capture"`
	Request        map[string]interface{} `json:"request"`
	RequestXml     string                 `json:"request_xml"`
	AssertResponse map[string]interface{} `json:"assert_response"`
}

type TestConfig struct {
	URL         string                 `json:"url"`
	Method      string                 `json:"method"`
	Headers     map[string]string      `json:"headers"`
	BaseRequest map[string]interface{} `json:"base_request"`
	Captures    []Capture              `json:"captures"`
	Cases       []TestCase             `json:"cases"`
	TypeAPI     string                 `json:"type"`
}

func (tc *TestConfig) GetRequest(testCase TestCase) map[string]interface{} {
	requestBody := make(map[string]interface{})
	for k, v := range tc.BaseRequest {
		requestBody[k] = v
	}
	for k, v := range testCase.Request {
		requestBody[k] = v
	}
	return requestBody
}

func (tc *TestConfig) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"url":          tc.URL,
		"method":       tc.Method,
		"headers":      tc.Headers,
		"base_request": tc.BaseRequest,
		"captures":     tc.Captures,
		"cases":        tc.Cases,
		"type":         tc.TypeAPI,
	}
}
func (tc *TestConfig) ToJSON() string {
	return fmt.Sprintf("%v", tc.ToMap())
}
func (tc *TestConfig) ToByteList() []byte {
	return []byte(tc.ToJSON())
}
func (tc *TestConfig) Copy() *TestConfig {
	newHeaders := make(map[string]string)
	for k, v := range tc.Headers {
		newHeaders[k] = v
	}

	newBaseRequest := make(map[string]interface{})
	for k, v := range tc.BaseRequest {
		newBaseRequest[k] = v
	}

	newCaptures := make([]Capture, len(tc.Captures))
	copy(newCaptures, tc.Captures)

	newCases := make([]TestCase, len(tc.Cases))
	copy(newCases, tc.Cases)

	return &TestConfig{
		URL:         tc.URL,
		Method:      tc.Method,
		Headers:     newHeaders,
		BaseRequest: newBaseRequest,
		Captures:    newCaptures,
		Cases:       newCases,
		TypeAPI:     tc.TypeAPI,
	}
}
func FindTestIndex(cases []TestCase, caseID int) int {
	for i, c := range cases {
		if c.ID == caseID {
			return i
		}
	}
	return -1
}

type TestStep struct {
	Step            string                   `json:"step"`
	File            string                   `json:"file"`
	TestType        string                   `json:"test_type"`
	CaseID          int                      `json:"case_id"`
	CaseName        string                   `json:"case_name"`
	CaseDescription string                   `json:"case_description"`
	CaseTags        []string                 `json:"case_tags"`
	RequestURL      string                   `json:"request_url"`
	RequestMethod   string                   `json:"request_method"`
	RequestType     string                   `json:"request_type"`
	RequestHeader   map[string][]string      `json:"request_header"`
	RequestBody     map[string]interface{}   `json:"request_body"`
	ResponseCode    int                      `json:"response_code"`
	ResponseHeader  map[string][]string      `json:"response_header"`
	ResponseType    string                   `json:"response_type"`
	ResponseBody    map[string]interface{}   `json:"response_body"`
	CapturedValues  map[string]interface{}   `json:"captured_values"`
	TestStatus      string                   `json:"test_status"`
	Error           string                   `json:"error"`
	AssertResponse  []map[string]interface{} `json:"assert_response"`
	AssertStatus    int                      `json:"assert_status"`
	Sources         string                   `json:"sources"`
	ResponseSource  string                   `json:"response_source"`
}

func (ts *TestStep) Init(testConfig TestConfig, testCase TestCase, step string, file string) {
	ts.Step = step
	ts.File = file
	ts.CaseID = testCase.ID
	ts.CaseName = testCase.Name
	ts.CaseDescription = testCase.Description
	ts.CaseTags = testCase.Tags
	ts.RequestURL = testConfig.URL
	ts.RequestMethod = testConfig.Method
	ts.RequestType = testConfig.TypeAPI
	ts.AssertStatus = testCase.Status
	rootFile, _ := filepath.Abs(ts.File)
	ts.Sources = fmt.Sprintf("file://%s", rootFile)
}
func MakeInitTestSteps(testConfigs []*TestConfig) []TestStep {
	testSteps := []TestStep{}
	step := 1
	for _, testConfig := range testConfigs {
		subStep := 1
		for _, testCase := range testConfig.Cases {
			tStep := fmt.Sprintf("%d", step)
			if len(testConfig.Cases) > 1 {
				tStep = fmt.Sprintf("%d.%d", step, subStep)
			}
			testStep := TestStep{}
			testStep.Step = tStep
			testStep.RequestURL = testConfig.URL
			testStep.RequestMethod = testConfig.Method
			testStep.RequestBody = testConfig.GetRequest(testCase)
			testStep.CaseID = testCase.ID
			testStep.CaseName = testCase.Name
			testStep.TestType = testConfig.TypeAPI
			testSteps = append(testSteps, testStep)
			subStep++
		}
		step++
	}
	return testSteps
}
func (ts *TestStep) SetTitle(step int, fileName, url, method string, caseId int, caseName string, assert map[string]interface{}, assertStatus int, testType string) {
	ts.Step = fmt.Sprintf("%d", step)
	ts.File = fileName
	ts.CaseID = caseId
	ts.CaseName = caseName
	ts.RequestURL = url
	ts.RequestMethod = method
	ts.AssertResponse = []map[string]interface{}{assert}
	ts.AssertStatus = assertStatus
	ts.TestType = testType
}

type TestReport struct {
	Title    string     `json:"title"`
	Status   string     `json:"status"`
	Duration string     `json:"duration"`
	Error    string     `json:"error"`
	DateTime string     `json:"date_time"`
	TestStep []TestStep `json:"test_step"`
}

func (r *TestReport) InitReport() {
	r.Title = "Test Report"
	r.Status = "initial"
	r.Duration = "0 ms"
	r.DateTime = time.Now().Format("2006-01-02-15-04-05")
	r.TestStep = []TestStep{}
	r.Error = ""
}
func (r *TestReport) GetCountByStatus(ststue string) int {
	count := 0
	for _, ts := range r.TestStep {
		if ts.TestStatus == ststue {
			count++
		}
	}
	return count
}
func (r *TestReport) SaveReportToFile(elapsed float64) error {
	for _, ts := range r.TestStep {
		if ts.Error != "" {
			r.Error += "- " + ts.Error + "\n"
		}
	}
	isPass := r.Error == ""
	if isPass {
		r.Status = "pass"
	} else {
		r.Status = "fail"
	}
	r.Duration = fmt.Sprintf("%f seconds", elapsed)
	report, err := json.Marshal(r)
	if err != nil {
		AddErrorReport(fmt.Sprintf("Error marshaling report: %v", err))
		return err
	}

	dirName := fmt.Sprintf("reports/report-%s", TimeSave)
	pwd, _ := os.Getwd()
	dirName = filepath.Join(pwd, dirName)
	nameFile := filepath.Join(dirName, "report.json")

	if err := os.MkdirAll(dirName, 0755); err != nil {
		return fmt.Errorf("error creating directories: %w", err)
	}

	file, err := os.Create(nameFile)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()
	report = []byte(replaceSpecialCharacters(string(report)))
	// Write the report ensuring UTF-8 encoding
	_, err = file.Write(report)
	if err != nil {
		return fmt.Errorf("error writing report to file: %w", err)
	}
	fmt.Printf("\n\nReport saved to %s\n\n", nameFile)
	return nil
}

func AddErrorReport(err string) {
	PrintRed(err)
}

type TestRunner struct {
	DisplayName string        `json:"display_name"`
	File        string        `json:"file"`
	Id          []interface{} `json:"id"`
}

type Runner struct {
	Mode  string       `json:"mode"`
	Tests []TestRunner `json:"tests"`
}

type FilesForTest struct {
	fileName string
	caseIDs  []string
}
type FileListForTest struct {
	fileList []FilesForTest
}

func (fl *FileListForTest) AddFile(fileName string, caseID []string) {
	fl.fileList = append(fl.fileList, FilesForTest{fileName: fileName, caseIDs: caseID})
}

func (fl *FileListForTest) AddCaseId(fileName string, caseID string) {
	for i, f := range fl.fileList {
		if f.fileName == fileName {
			fl.fileList[i].caseIDs = append(fl.fileList[i].caseIDs, caseID)
			return
		}
	}
	fl.AddFile(fileName, []string{caseID})
}
func (fl *FileListForTest) GetCountTest() int {
	count := 0
	for _, f := range fl.fileList {
		count += len(f.caseIDs)
	}
	return count
}

type Sumary struct {
	Totals  int
	Passed  int
	Failed  int
	Skipped int
}

func (s *Sumary) SetValueWithReport(report TestReport) {
	s.Totals = len(report.TestStep)
	s.Passed = report.GetCountByStatus("PASS")
	s.Failed = report.GetCountByStatus("FAIL")
	s.Skipped = report.GetCountByStatus("SKIP")
}
func (s *Sumary) PrintSumary() {
	fmt.Printf("\n\n\tTotal: %d\n\tPassed: %d\n\tFailed: %d\n\tSkipped: %d\n\n", s.Totals, s.Passed, s.Failed, s.Skipped)
}
func (s *Sumary) AddTotal(count int) {
	s.Totals += count
}
