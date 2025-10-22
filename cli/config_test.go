package cli

import "testing"

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

func TestUnmarshalConfig(t *testing.T) {
	cfg, err := unmarshalConfig([]byte(txtConfig))
	if err != nil {
		t.Error(err)
	}
	want := "tmpl/"
	if cfg.TemplateBaseDir != want {
		t.Errorf("unmarshalConfig: TemplateBaseDir = %q, want %q", cfg.TemplateBaseDir, want)
	}

	if len(cfg.Pages) != 5 {
		t.Errorf("unmarshalConfig: invalid number of pages: expected=%d, actual=%d", 5, len(cfg.Pages))
	}
	if len(cfg.Templates) != 4 {
		t.Errorf("unmarshalConfig: invalid number of templates: expected=%d, actual=%d", 4, len(cfg.Templates))
	}
}
