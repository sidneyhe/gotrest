// serve.go
package cmd

import (
	"encoding/json"
	"fmt"
	"gotrest/parse"
	"gotrest/testsuite"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

func RuntcHandler(w http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadAll(req.Body)
	testCase := &parse.TestCase{}
	err = json.Unmarshal(data, testCase)
	fmt.Printf("err =%v", err)

	var tcs []*parse.TestCase
	tcs = append(tcs, testCase)

	ts := testsuite.NewTestSession()
	ts.Start(tcs)
	summary := ts.Finish()
	data, err = json.Marshal(summary)
	w.Write(data)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "gotrest use as a service",
	Long:  `gotrest use as a service`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World!")
		http.HandleFunc("/api/runtc", RuntcHandler)

		http.ListenAndServe(":8888", nil)
	},
}
