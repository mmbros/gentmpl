package main

import (
	"bytes"
	"testing"
)

func TestPageExecute(t *testing.T) {
	var testCases = []struct {
		page PageEnum
	}{
		{PageInh1},
		{PageInh2},
	}

	wr := new(bytes.Buffer)

	for _, tc := range testCases {
		if err := tc.page.Execute(wr, nil); err != nil {
			t.Errorf("page.Execute: %s", err)
		}

	}

}
