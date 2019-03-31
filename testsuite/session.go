// session.go
package testsuite

import (
	"gotrest/parse"
	"gotrest/summary"
	"gotrest/testcase"
)

type TestSession struct {
	summary *summary.TestSummary
}

func (ts *TestSession) Run(tc *parse.TestCase) error {
	tcrDetail := ts.summary.NewTCRDetail(tc.Name)
	tcRunner := testcase.NewRunner(tc, tcrDetail)
	if err := tcRunner.PreRun(); err != nil {
		ts.summary.StatIncFailCase()
		return err
	}

	if err := tcRunner.Run(); err != nil {
		ts.summary.StatIncFailCase()
		return err
	}

	if err := tcRunner.PostRun(); err != nil {
		ts.summary.StatIncFailCase()
		return err
	}

	ts.summary.StatIncSuccessCase()

	return nil
}

func (ts *TestSession) Start(tcs []*parse.TestCase) error {
	for _, tc := range tcs {
		// tcrDetail := ts.summary.NewTCRDetail(tc.Name)
		// tcRunner := testcase.NewRunner(tc, tcrDetail)
		// if tcRunner.PreRun() != nil {
		// 	ts.summary.StatIncFailCase()
		// 	continue
		// }

		// if tcRunner.Run() != nil {
		// 	ts.summary.StatIncFailCase()
		// 	continue
		// }

		// if tcRunner.PostRun() != nil {
		// 	ts.summary.StatIncFailCase()
		// 	continue
		// }

		// ts.summary.StatIncSuccessCase()
		ts.Run(tc)
	}

	return nil
}

func (ts *TestSession) Finish() *summary.TestSummary {
	ts.summary.SaveResult()
	return ts.summary
}

func NewTestSession() *TestSession {
	return &TestSession{
		summary: summary.NewTestSummary(),
	}
}
