package cli

import (
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
		})
	}
}

// func Test_writeOutput(t *testing.T) {
// 	type args struct {
// 		path string
// 		fn   func(io.Writer) error
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := writeOutput(tt.args.path, tt.args.fn); (err != nil) != tt.wantErr {
// 				t.Errorf("writeOutput() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
