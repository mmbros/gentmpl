package main

import (
	"os"
	"testing"
)

func TestPageExecute(t *testing.T) {
	var testCases = []struct {
		page PageEnum
	}{
		{PageInh1},
		{PageInh2},
	}

	wr := os.Stdout

	for _, tc := range testCases {
		if err := tc.page.Execute(wr, nil); err != nil {
			t.Errorf("page.Execute: %s", err)
		}

	}

}
