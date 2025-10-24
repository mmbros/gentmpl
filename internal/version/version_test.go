package version

import (
	"bytes"
	"strings"
	"testing"
)

func Test_Print(t *testing.T) {
	type args struct {
		appname string
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		{
			name:  "test",
			args:  args{appname: "gentmpl"},
			wantW: "app version",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			Print(w, tt.args.appname)
			if gotW := w.String(); !strings.Contains(gotW, tt.wantW) {
				t.Errorf("cmdVersion() = %v, do not contain want %v", gotW, tt.wantW)
			}
		})
	}
}
