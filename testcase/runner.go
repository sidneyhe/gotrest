// runner.go
package testcase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gotrest/comparator"
	"gotrest/parse"
	"gotrest/summary"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	CookiePrefix = "cookie."
	HeaderPrefix = "header."
	BodyPrefix   = "body."
)

// Runner runs parsed tests.
type Runner struct {
	vars     map[string]interface{}
	testCase *parse.TestCase
	// DoRequest makes the request and returns the response.
	// By default uses http.DefaultClient.Do.
	DoRequest func(r *http.Request) (*http.Response, error)
	// ParseBody is the function to use to attempt to parse
	// response bodies to make data available for assertions.
	//ParseBody func(r io.Reader) (interface{}, error)
	// Log is the function to log to.
	//Log func(string)
	// Verbose is the function that logs verbose debug information.
	// Verbose func(...interface{})
	// NewRequest makes a new http.Request. By default, uses http.NewRequest.
	NewRequest func(method, urlStr string, body io.Reader) (*http.Request, error)

	tcrDetail *summary.TCRDetail
}

func NewRunner(tc *parse.TestCase, detail *summary.TCRDetail) *Runner {
	runner := &Runner{
		vars:       make(map[string]interface{}),
		testCase:   tc,
		DoRequest:  http.DefaultTransport.RoundTrip,
		NewRequest: http.NewRequest,
		tcrDetail:  detail,
	}

	for key, value := range tc.Variables {
		runner.vars[key] = value
	}

	fmt.Printf("runner vars =%v, tc vars =%v", runner.vars, tc.Variables)

	return runner
}

func (r *Runner) PreRun() error {
	return nil
}

func (r *Runner) getRealValue(value interface{}, step *parse.TeseStep) string {
	realValue := ""
	if value == nil {
		return realValue
	}

	if strVal, ok := value.(string); ok {
		if strings.HasPrefix(strVal, "$") {
			// env variable
			if strings.HasPrefix(strVal, "${ENV(") || strings.HasPrefix(strVal, "${env(") {
				envVar := strVal[6:(len(strVal) - 2)]
				realValue = os.Getenv(envVar)

				fmt.Printf("envVar =%v, realValue=%v", envVar, realValue)
			} else {
				useVar := strVal[1:]
				// first, find it in step.Vars
				tmp, _ := step.Variables[useVar]
				if tmp != nil {
					realValue = fmt.Sprintf("%v", tmp)
				} else {
					// second, find it in global Vars
					tmp, _ = r.vars[useVar]
					fmt.Printf("useVar =%v, tmp=%v", useVar, tmp)
					if tmp != nil {
						realValue = fmt.Sprintf("%v", tmp)
					}
				}
			}
		} else {
			realValue = strVal
		}
	} else {
		realValue = fmt.Sprintf("%v", value)
	}

	return realValue
}

func (r *Runner) buildRequset(step *parse.TeseStep) (*http.Request, error) {
	request := step.Request
	url := request.URL
	if strings.HasPrefix(url, "http") == false {
		url = r.testCase.BaseUrl + url
	}

	body := bytes.NewBufferString("")

	// make request
	httpReq, err := r.NewRequest(request.Method, url, body)
	if err != nil {
		fmt.Printf("httpReq =%v\n", err)
		//panic("invalid request")
		return nil, err
	}

	// fill request
	// parameters
	query := httpReq.URL.Query()
	param := request.Parameter
	for key, value := range param {
		varValue := r.getRealValue(value, step)
		query.Add(key, varValue)
	}
	httpReq.URL.RawQuery = query.Encode()
	fmt.Printf("param =%v, query=%v", param, httpReq.URL.RawQuery)

	// fill header
	for key, value := range request.Header {
		varValue := r.getRealValue(value, step)
		httpReq.Header.Add(key, varValue)
	}

	// fill body

	return httpReq, err
}

func extractFromMap(dotKey string, mapResult map[string]interface{}) (interface{}, error) {
	fmt.Println("extractFromMap, mapResult =", mapResult)
	tmpMap := mapResult
	keyList := strings.Split(dotKey, ".")
	lastKey := keyList[len(keyList)-1]
	ok := true
	for i := 0; i < (len(keyList) - 1); i++ {
		tmpMap, ok = tmpMap[keyList[i]].(map[string]interface{})
		if !ok {
			// return err
			return nil, errors.New("ExtractorError")
		}
	}

	if val, ok := tmpMap[lastKey]; ok {
		return val, nil
	} else {
		return nil, errors.New("ExtractorError")
	}
}

func (r *Runner) parseResponse(resp *http.Response, step *parse.TeseStep) error {
	ckMap := make(map[string]string)
	cks := resp.Cookies()
	for _, ck := range cks {
		ckMap[ck.Name] = ck.Value
	}

	fmt.Println("ckMap =", ckMap)
	mapResult := make(map[string]interface{})

	if contentType := resp.Header.Get("Content-Type"); strings.Contains(contentType, "json") {

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("err=,", err)
			return err
		}

		err = json.Unmarshal(data, &mapResult)
		if err != nil {
			fmt.Println("json.Unmarshal,err=,", err)
			return err
		}
	}

	// extract
	for exDst, exSrc := range step.Extract {
		if strings.HasPrefix(exSrc, CookiePrefix) {
			ckKey := strings.TrimPrefix(exSrc, CookiePrefix)
			if val, ok := ckMap[ckKey]; ok {
				r.vars[exDst] = val
			}
		} else if strings.HasPrefix(exSrc, BodyPrefix) { // start with bod
			dotKey := strings.TrimPrefix(exSrc, BodyPrefix)
			val, err := extractFromMap(dotKey, mapResult)
			if err != nil {
				return err
			}

			r.vars[exDst] = val
		}
	}

	// validation
	for _, checker := range step.Validate {
		checkRes := false
		var err error = nil
		var got interface{} = nil
		comp := comparator.CreateComparator(checker.Comparator)
		if checker.Check == "status_code" {
			got = resp.StatusCode
		} else if strings.HasPrefix(checker.Check, CookiePrefix) {
			key := strings.TrimPrefix(checker.Check, CookiePrefix)
			got, _ = ckMap[key]
		} else if strings.HasPrefix(checker.Check, BodyPrefix) {
			dotKey := strings.TrimPrefix(checker.Check, BodyPrefix)
			got, err = extractFromMap(dotKey, mapResult)
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(checker.Check, HeaderPrefix) {
			hKey := strings.TrimPrefix(checker.Check, HeaderPrefix)
			got = resp.Header.Get(hKey)
		}

		checkRes = comp.Compare(got, checker.Expect)
		fmt.Println("checkRes =", checkRes)
		if checkRes == false {
			return errors.New(fmt.Sprintf("ValidateFail, check =%v,expect=%v, got=%v", checker.Check, checker.Expect, got))
		}
	}

	return nil
}

func (r *Runner) Run() error {

	for _, step := range r.testCase.Teststeps {
		stepDetail := r.tcrDetail.NewTeststepDetail(step.Name)

		// make request
		httpReq, err := r.buildRequset(&step)
		if err != nil {
			fmt.Printf("httpReq =%v\n", err)
			//panic("invalid request")
			return err
		}

		// perform request
		requsetAt := time.Now().UnixNano()
		httpRes, err := r.DoRequest(httpReq)
		if err != nil {
			fmt.Printf("httpReq =%v\n", err)
			//panic("invalid request")
			return err
		}

		respDuration := time.Now().UnixNano() - requsetAt

		defer httpRes.Body.Close()

		err = r.parseResponse(httpRes, &step)

		result := summary.TestStepSuccesses
		if err != nil {
			fmt.Println("err =%v", err)
			result = summary.TestStepFailures
		}

		stepDetail.SaveResult(result, respDuration)
		r.tcrDetail.StatIncStepResult(result)

		fmt.Printf("httpRes code =%s, checkRes=%v", httpRes.Status, err)
	}

	r.tcrDetail.SaveResult(summary.TestCaseSuccess, "")

	return nil
}

func (r *Runner) PostRun() error {
	return nil
}
