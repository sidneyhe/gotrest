// testCase_test.go
package parse

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestTestCaseToJson(t *testing.T) {
	tc := &TestCase{}
	tc.Config.Name = "test"
	testStep1 := TeseStep{
		Name: "test01",
		Request: map[string]interface{}{
			"url":    "http://127.0.0.1:5000/api/hello",
			"method": "GET",
		},
		Validate: []ValidateRule{
			{
				Check:      "status_code",
				Comparator: "eq",
				Expect:     200,
			},
		},
	}
	tc.Teststeps = append(tc.Teststeps, testStep1)
	data, err := json.Marshal(tc)
	fmt.Printf("data =%v, err =%v", string(data), err)
}
