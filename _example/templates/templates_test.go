package templates

import (
	"bytes"
	"testing"
)

func TestPageExecute(t *testing.T) {
	var testCases = []struct {
		page PageEnum
	}{
		{PagePag1},
		{PagePag2},
		{PagePag3},
		{PageInh1},
		{PageInh2},
	}

	InitTemplates()

	wr := new(bytes.Buffer)

	for _, tc := range testCases {
		if err := tc.page.Execute(wr, nil); err != nil {
			t.Errorf("page.Execute: %s", err)
		}

	}

}

func Test_file2path(t *testing.T) {
	tests := []struct {
		name string
		file string
		want string
	}{
		{
			name: "dot prefix",
			file: ".xxx",
			want: ".xxx",
		},
		{
			name: "empty file",
			file: "",
			want: "",
		},
		{
			name: "path separator prefix",
			file: "/xxx",
			want: "/xxx",
		},
		{
			name: "standard path",
			file: "folder/file.ext",
			want: "tmpl/folder/file.ext",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := file2path(tt.file)
			if got != tt.want {
				t.Errorf("file2path(%q), got %q, want %q", tt.file, got, tt.want)
			}
		})
	}
}
