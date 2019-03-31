// tcrdetail.go
package summary

import (
	"fmt"
	_ "fmt"
	"time"
)

const (
	TestStepUnexpectedSuccesses = "unexpectedSuccesses"
	TestStepExpectedFailures    = "expectedFailures"
	TestStepSuccesses           = "successes"
	TestStepFailures            = "failures"
	TestStepSkipped             = "skipped"
	TestStepErrors              = "errors"
	TestStepTotal               = "total"
)

type TestStepDetail struct {
	Name         string `json:"name"`
	Result       string `json:"result"`
	RespDuration int64  `json:"respDuration"`
	DetailInfo   string `json:"detailInfo"`
}

func (tsd *TestStepDetail) SaveResult(result string, respDurateion int64) {
	tsd.Result = result
	tsd.RespDuration = respDurateion
}

// testcase result detail
type TCRDetail struct {
	CaseName       string            `json:"casename"`
	Result         string            `json:"result"`
	startAt        int64             `json:"-"`
	Duration       int64             `json:"duration"`
	Error          string            `json:"error"`
	Teststeps      map[string]int    `json:"teststeps"`
	TeststepDetail []*TestStepDetail `json:"teststepDetail"`
}

func (tcrDetail *TCRDetail) SaveResult(result string, errorInfo string) {
	tcrDetail.Result = result
	tcrDetail.Error = errorInfo
	tcrDetail.Duration = time.Now().UnixNano() - tcrDetail.startAt
}

func (tcrDetail *TCRDetail) NewTeststepDetail(stepName string) *TestStepDetail {
	detail := &TestStepDetail{
		Name: stepName,
	}

	tcrDetail.Teststeps[TestStepTotal]++
	tcrDetail.TeststepDetail = append(tcrDetail.TeststepDetail, detail)
	return detail
}

func (tcrDetail *TCRDetail) StatIncStepResult(stepResult string) {
	if _, exist := tcrDetail.Teststeps[stepResult]; exist {
		tcrDetail.Teststeps[stepResult]++
	} else {
		panic(fmt.Sprintf("invalid Result :%s", stepResult))
	}
}

func NewTCRDetail(caseName string) *TCRDetail {
	stepStat := make(map[string]int)
	stepStat[TestStepUnexpectedSuccesses] = 0
	stepStat[TestStepExpectedFailures] = 0
	stepStat[TestStepSuccesses] = 0
	stepStat[TestStepFailures] = 0
	stepStat[TestStepSkipped] = 0
	stepStat[TestStepErrors] = 0

	return &TCRDetail{
		CaseName:  caseName,
		Teststeps: stepStat,
		startAt:   time.Now().UnixNano(),
	}
}
