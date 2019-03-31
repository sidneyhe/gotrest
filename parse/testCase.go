// testCase.go
package parse

import (
	_ "fmt"
)

// {"config": {
//                       "name": "desc1",
//                       "path": "testcase1_path",
//                       "variables": [],                    # optional
//                   },
//                   "teststeps": [
//                       # test data structure
//                       {
//                           'name': 'test step desc1',
//                           'variables': [],    # optional
//                           'extract': [],      # optional
//                           'validate': [],
//                           'request': {}
//                       },
//                       test_dict_2   # another test dict
//                   ]
//               }

// type TestCaseConfig struct {
// }

type ValidateRule struct {
	Check      string      `json:"check"`
	Comparator string      `json:"comparator"`
	Expect     interface{} `json:"expect"`
}

type TestStepRequest struct {
	Method    string                 `json:"method"`
	URL       string                 `json:"url"`
	Parameter map[string]interface{} `json:"parameter"`
	Header    map[string]interface{} `json:"header"`
}

type TeseStep struct {
	Name      string                 `json:"name"`
	Variables map[string]interface{} `json:"variables"`
	Extract   map[string]string      `json:"extract"`
	Validate  []ValidateRule         `json:"validate"`
	Request   TestStepRequest        `json:"request"`
}

type TestCase struct {
	Version   string                 `json:"version"` // test case format version
	Name      string                 `json:"name"`
	BaseUrl   string                 `json:"baseUrl"`
	Variables map[string]interface{} `json:"variables"`
	Teststeps []TeseStep             `json:"tesesteps"`
}
