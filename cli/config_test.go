package cli

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	txtConfig = `template_base_dir = "tmpl/"
[templates]
flat = ["flat/footer.tmpl", "flat/header.tmpl", "flat/page1.tmpl", "flat/page2and3.tmpl"]
inhbase = ["inheritance/base.tmpl"]
inh1 = ["inhbase", "inheritance/content1.tmpl"]
inh2 = ["inhbase", "inheritance/content2.tmpl"]

[pages]
Pag1 = {template="flat", base="page-1"}
Pag2 = {template="flat", base="page-2"}
Pag3 = {template="flat", base="page-3"}
Inh1 = {template="inh1"}
Inh2 = {template="inh2"}
`
)

func Test_unmarshalConfig(t *testing.T) {

	cfg := &config{}
	cfg.TemplateBaseDir = "tmpl/"
	cfg.Templates = map[string][]string{
		"flat":    {"flat/footer.tmpl", "flat/header.tmpl", "flat/page1.tmpl", "flat/page2and3.tmpl"},
		"inh1":    {"inhbase", "inheritance/content1.tmpl"},
		"inh2":    {"inhbase", "inheritance/content2.tmpl"},
		"inhbase": {"inheritance/base.tmpl"},
	}
	cfg.Pages = map[string]struct {
		Template string
		Base     string
	}{
		"Inh1": {Template: "inh1"},
		"Inh2": {Template: "inh2"},
		"Pag1": {Template: "flat", Base: "page-1"},
		"Pag2": {Template: "flat", Base: "page-2"},
		"Pag3": {Template: "flat", Base: "page-3"},
	}

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *config
		wantErr bool
	}{
		{
			name:    "parse-success",
			args:    args{data: []byte(txtConfig)},
			want:    cfg,
			wantErr: false,
		},
		{
			name:    "parse-error",
			args:    args{data: []byte(`[pages] Pag1 = template="flat", base=page-1"}`)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unmarshalConfig(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ToSlice() mismatch (-want +got):\n%s", diff)
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("unmarshalConfig() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func Test_loadConfigFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadConfigFromFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadConfigFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cmdHelp(t *testing.T) {

	tests := []struct {
		name string
		want string
	}{
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

	var clinfo cmdlineInfo

	fs := clinfo.newFlagSet()
	fs.Parse([]string{"-h"})

	if !clinfo.help {
		t.Error("used command line option -h, but cmdHelp() is not called")
		t.FailNow()
	}

	cmdHelp(w, fs)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(w.String(), tt.want) {
				t.Errorf("cmdHelp() message does not contain substring %q", tt.want)
				t.Log(w.String())
			}
		})
	}
}

// func Test_parseArgs(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want *cmdlineInfo
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got, _ := parseArgs(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("parseArgs() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func Test_parseConfig(t *testing.T) {
	type args struct {
		args *cmdlineInfo
	}
	tests := []struct {
		name    string
		args    args
		want    *config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConfig(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeOutput(t *testing.T) {
	type args struct {
		path string
		fn   func(io.Writer) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeOutput(tt.args.path, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("writeOutput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cmdGenPackage(t *testing.T) {
	type args struct {
		cfg *config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cmdGenPackage(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("cmdGenPackage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cmdGenConfig(t *testing.T) {
	type args struct {
		args *cmdlineInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cmdGenConfig(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("cmdGenConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Run(); got != tt.want {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
