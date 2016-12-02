package cli

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

/*
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
*/
var templates = map[string][]string{
	"flat":    {"flat/footer.tmpl", "flat/header.tmpl", "flat/page1.tmpl", "flat/page2and3.tmpl"},
	"inhbase": {"inheritance/base.tmpl"},
	"inh1":    {"inhbase", "inheritance/content1.tmpl"},
	"inh2":    {"inhbase", "inheritance/content2.tmpl"},
}
var pages = map[string]struct {
	Template string
	Base     string
}{
	"Pag1": {"flat", "page-1"},
	"Pag2": {"flat", "page-2"},
	"Pag3": {"flat", "page-3"},
	"Inh1": {"inh1", ""},
	"Inh2": {"inh2", ""},
}
var fileinfos = []struct {
	path    string
	content string
}{
	{"flat/header.tmpl", `{{define "header"}}<html><head></head><body>{{end}}`},
	{"flat/footer.tmpl", `{{define "footer"}}</body></html>{{end}}`},
	{"flat/page1.tmpl", `{{define "page-1"}}{{template "header"}}Page 1{{template "footer"}}{{end}}`},
	{"flat/page2and3.tmpl", `{{define "page-2"}}{{template "header"}}Page 2{{template "footer"}}{{end}}{{define "page-3"}}{{template "header"}}Page 3{{template "footer"}}{{end}}`},
	{"inheritance/base.tmpl", `<html><head></head><body>{{template "content" .}}</body></html>`},
	{"inheritance/content1.tmpl", `{{define "content"}}content 1{{end}}`},
	{"inheritance/content2.tmpl", `{{define "content"}}content 2{{end}}`},
}

func writefile(fullpath, content string) error {
	folder := filepath.Dir(fullpath)
	if err := os.MkdirAll(folder, 0777); err != nil {
		return err
	}
	if err := ioutil.WriteFile(fullpath, []byte(content), 0666); err != nil {
		return err
	}
	return nil
}

// create a tmp dir with the template files
// returns the name of the created folder
func setupTemplates() string {
	// creates a tmp root directory
	tmpdir, err := ioutil.TempDir("", "gentmpltest_")
	if err != nil {
		panic(err)
	}
	// creates the templates files
	for _, fi := range fileinfos {
		if err := writefile(filepath.Join(tmpdir, fi.path), fi.content); err != nil {
			panic(err)
		}
	}
	// return the name of the tmp root directory
	return tmpdir
}

func setupConfig(folder string) *Config {
	outfile := filepath.Join(folder, defaultPackageName+".go")
	cfg := &Config{
		OutputFile:      outfile,
		Folder:          folder,
		FormatOutput:    true,
		Debug:           defaultDebug,
		PackageName:     "gentmpltest",
		DefaultPageBase: defaultDefaultPageBase,
		PageEnumType:    defaultPageEnumType,
		PageEnumPrefix:  "Page",
		Templates:       templates,
		Pages:           pages,
	}
	return cfg
}

func TestWrite(t *testing.T) {
	// setup
	tmpdir := setupTemplates()
	cfg := setupConfig(tmpdir)

	// do test
	err := cfg.WriteModule()
	if err != nil {
		t.Fatalf(err.Error())
	}

	// clean up
	//defer os.RemoveAll(tmpdir)
}
