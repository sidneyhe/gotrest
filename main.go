// main.go
package main

import (
	"gotrest/cmd"
)

// func RuntcHandler(w http.ResponseWriter, req *http.Request) {
// 	data, err := ioutil.ReadAll(req.Body)
// 	testCase := &parse.TestCase{}
// 	err = json.Unmarshal(data, testCase)
// 	fmt.Printf("err =%v", err)

// 	var tcs []*parse.TestCase
// 	tcs = append(tcs, testCase)

// 	ts := testsuite.NewTestSession()
// 	ts.Start(tcs)
// 	summary := ts.GetSummary()
// 	data, err = json.Marshal(summary)
// 	w.Write(data)
// }

func main() {
	// fmt.Println("Hello World!")
	// http.HandleFunc("/api/runtc", RuntcHandler)

	// http.ListenAndServe(":8888", nil)

	cmd.Execute()
}
