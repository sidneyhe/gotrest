// root.go
package cmd

import (
	"encoding/json"
	"fmt"
	"gotrest/parse"
	"gotrest/testsuite"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

//LoadConfig load config
func loadTestCase(path string) (*parse.TestCase, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tc := &parse.TestCase{}
	if err = json.Unmarshal(data, tc); err != nil {
		return nil, err
	}
	return tc, nil
}

func runTest(cmd *cobra.Command, args []string) {
	ts := testsuite.NewTestSession()
	fmt.Printf("rootCmd, args =%v\n", args)
	for _, file := range args {
		tc, _ := loadTestCase(file)
		ts.Run(tc)
	}

	summary := ts.Finish()

	data, _ := json.Marshal(summary)

	fmt.Printf("\nsummary:%v\n", string(data))
}

var rootCmd = &cobra.Command{
	Use:   "gotrest",
	Short: "a tool to test restful api writed by golang, gotrest=go test rest",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run:   runTest,
}

func Execute() {
	rootCmd.AddCommand(serveCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
