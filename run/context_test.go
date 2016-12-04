package run

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jteeuwen/go-bindata"
)

const (
	// tmp folder prefix
	testDirPrefix = "gentmpl_"

	// base template directory
	baseTemplateDir = "tmpl"
)

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

var results = map[string]string{
	"Pag1": "<html><head></head><body>Page 1</body></html>",
	"Pag2": "<html><head></head><body>Page 2</body></html>",
	"Pag3": "<html><head></head><body>Page 3</body></html>",
	"Inh1": "<html><head></head><body>content 1</body></html>",
	"Inh2": "<html><head></head><body>content 2</body></html>",
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

func writeFile(fullpath, content string) error {
	folder := filepath.Dir(fullpath)
	if err := os.MkdirAll(folder, 0777); err != nil {
		return err
	}
	if err := ioutil.WriteFile(fullpath, []byte(content), 0666); err != nil {
		return err
	}
	return nil
}

func errorLike(err error, msg string) bool {
	return strings.Contains(err.Error(), msg)
}

// create a tmp dir with the template files
// returns the name of the created folder
func setupDirTemplates() string {
	// creates a tmp root directory
	tmpdir, err := ioutil.TempDir("", testDirPrefix)
	if err != nil {
		panic(err)
	}
	// creates the templates files
	for _, fi := range fileinfos {
		out := filepath.Join(tmpdir, baseTemplateDir, fi.path)
		if err := writeFile(out, fi.content); err != nil {
			panic(err)
		}
	}
	// return the name of the tmp root directory
	return tmpdir
}

func TestCheck(t *testing.T) {
	var (
		ctx *Context
		err error
	)

	checkErr := func(err error, testCase, errLike string) {
		if err == nil {
			t.Errorf("%s: expected error like %q; no error found", testCase, errLike)
		}
		if !errorLike(err, errLike) {
			t.Errorf("%s: expected error like %q; found error %q", testCase, errLike, err.Error())
		}
	}

	// test no pages
	ctx = &Context{}
	err = ctx.Check()
	checkErr(err, "no pages", "No pages found")

	// test page without a template
	ctx = &Context{
		Pages: map[string]struct {
			Template string
			Base     string
		}{
			"Pag": {},
		},
	}
	err = ctx.Check()
	checkErr(err, "page with no template", "Page must have a template")

	// test template's page not found
	tmpl := map[string][]string{}
	for _, key := range []string{"flat", "inh1"} {
		tmpl[key] = templates[key]
	}
	ctx = &Context{Pages: pages, Templates: tmpl}
	err = ctx.Check()
	checkErr(err, "no template", "Template not found for page")

	// test cyclic templates
	ctx = &Context{
		Pages: map[string]struct {
			Template string
			Base     string
		}{
			"Pag": {"t1", ""},
		},
		Templates: map[string][]string{
			"t1": {"p1", "t2"},
			"t2": {"p2", "t1"},
		},
	}
	err = ctx.Check()
	checkErr(err, "cyclic templates", "Found invalid cyclic template")

}

func TestWritePackage(t *testing.T) {
	ctx := &Context{Pages: pages, Templates: templates}
	w := new(bytes.Buffer)
	err := ctx.WritePackage(w)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestWriteConfig(t *testing.T) {
	ctx := &Context{Pages: pages, Templates: templates}
	w := new(bytes.Buffer)
	err := ctx.WriteConfig(w)
	if err != nil {
		t.Errorf(err.Error())
	}
}
func ctx2str(ctx *Context) string {
	var b bytes.Buffer
	writeBool := func(x bool) {
		ch := '0'
		if x {
			ch = '1'
		}
		b.WriteRune(ch)
	}
	writeSep := func() { b.WriteRune('-') }

	b.WriteString("as")
	writeBool(strings.ToLower(ctx.AssetManager) == "go-bindata")

	writeSep()
	b.WriteString("nc")
	writeBool(ctx.NoCache)

	writeSep()
	b.WriteString("fm")
	writeBool(ctx.FuncMap != "")

	//writeSep()
	//b.WriteString("nf")
	//writeBool(ctx.NoGoFormat)

	return b.String()

}

func writeBindata(out string, ctx *Context) error {
	c := bindata.NewConfig()
	c.Debug = ctx.NoCache
	c.Output = out
	c.Prefix = filepath.Clean(filepath.Join(filepath.Dir(out), "..", baseTemplateDir))
	c.Package = "main"
	c.Input = []bindata.InputConfig{
		bindata.InputConfig{
			Path:      c.Prefix,
			Recursive: true,
		},
	}

	return bindata.Translate(c)
}
func writeMain(out string, ctx *Context) error {
	const text = `package main

import (
	"fmt"
	"os"
)

func main(){
	var page PageEnum = PageInh1
	wr := os.Stdout

	if err := page.Execute(wr, nil); err != nil {
		fmt.Print(err)
	}
}
`
	return writeFile(out, text)

}

func TestAll(t *testing.T) {

	// setup
	tmpdir := setupDirTemplates()
	ctx := &Context{
		PackageName:     "main",
		TemplateBaseDir: filepath.Join("..", baseTemplateDir),
		Pages:           pages,
		Templates:       templates,
	}

	buf := new(bytes.Buffer)

	for _, nocache := range []bool{false, true} {
		ctx.NoCache = nocache

		for _, assetmngr := range []string{"", "go-bindata"} {
			ctx.AssetManager = assetmngr

			for _, funcmap := range []string{"", "funcMap"} {
				ctx.FuncMap = funcmap

				folder := ctx2str(ctx)

				// template.go
				out := filepath.Join(tmpdir, folder, "templates.go")
				buf.Reset()
				err := ctx.WritePackage(buf)
				if err == nil {
					err = writeFile(out, buf.String())
				}
				if err != nil {
					t.Errorf("%s/template %s", folder, err.Error())
				}

				// bindata.go
				if assetmngr != "" {
					out = filepath.Join(tmpdir, folder, "bindata.go")
					err = writeBindata(out, ctx)
					if err != nil {
						t.Errorf("%s/bindata %s", folder, err.Error())
					}
				}

				// main.go
				out = filepath.Join(tmpdir, folder, "main.go")
				err = writeMain(out, ctx)
				if err != nil {
					t.Errorf("%s/main %s", folder, err.Error())
				}
			}
		}
	}
	// clean up
	//defer os.RemoveAll(tmpdir)
}
