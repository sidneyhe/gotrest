// summary.go
package summary

import (
	"time"
)

const (
	TestCaseSuccess = "Success"
	TestCaseFail    = "Fail"
)

type Stat struct {
	Total   int32 `json:"total"`
	Success int32 `json:"success"`
	Fail    int32 `json:"fail"`
}

type TestSummary struct {
	Platform    map[string]interface{} `json:"platform"`
	Duration    int64                  `json:"duration"`
	StartAt     int64                  `json:"-"`
	StartAtDate string                 `json:"startAt"`
	Stat        Stat                   `json:"stat"`
	Detail      []*TCRDetail           `json:detail`
}

func (summary *TestSummary) NewTCRDetail(caseName string) *TCRDetail {
	detail := NewTCRDetail(caseName)
	summary.Stat.Total++
	summary.Detail = append(summary.Detail, detail)

	return detail
}

func (summary *TestSummary) StatIncSuccessCase() {
	summary.Stat.Success++
}

func (summary *TestSummary) StatIncFailCase() {
	summary.Stat.Fail++
}

func (summary *TestSummary) SaveResult() {
	summary.Duration = time.Now().UnixNano() - summary.StartAt
}

func NewTestSummary() *TestSummary {
	// get Platform
	now := time.Now()
	summary := &TestSummary{
		StartAt:     time.Now().UnixNano(),
		StartAtDate: now.String(),
	}

	return summary
}
