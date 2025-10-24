package cmdline

import (
	"flag"
	"strings"
	"testing"
)

func TestArgs_PrintHelp(t *testing.T) {

	const appName = "TestAppName"

	tests := []struct {
		name     string
		want     string
		dontWant string
	}{
		{
			name: "contains appName:",
			want: appName,
		},
		{
			name:     "does not contain gentmpl:",
			dontWant: "gentmpl",
		},
		{
			name: "contains Generate",
			want: "Generate the configuration file",
		},
		{
			name: "contains utility",
			want: "an utility that generates a go package",
		},
		{
			name: "contains row Options:",
			want: "Options:\n",
		},
	}

	w := &strings.Builder{}

	args := NewArgs(appName, flag.ContinueOnError)

	err := args.Parse([]string{"-h"})
	if err != nil {
		t.Errorf("PrintHelp: unexpected error: %v", err)
	}

	if !args.Help() {
		t.Error("used command line option -h, but PrintHelp is not called")
		t.FailNow()
	}

	args.PrintHelp(w)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.want != "") && !strings.Contains(w.String(), tt.want) {
				t.Errorf("cmdHelp() message does not contain substring %q", tt.want)
				t.Log(w.String())
			}
			if (tt.dontWant != "") && strings.Contains(w.String(), tt.dontWant) {
				t.Errorf("cmdHelp() message contain unexpected substring %q", tt.dontWant)
				t.Log(w.String())
			}
		})
	}
}
