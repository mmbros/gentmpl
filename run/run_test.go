package run

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jteeuwen/go-bindata"
)

const (
	// tmp folder prefix
	testDirPrefix = "gentmpl_"

	//  template base directory
	templateBaseDir = "tmpl"
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
		out := filepath.Join(tmpdir, templateBaseDir, fi.path)
		if err := writeFile(out, fi.content); err != nil {
			panic(err)
		}
	}
	// return the name of the tmp root directory
	return tmpdir
}

// ctx2str returns a short string that represents the context.
//   - as -> AssetManager : 0=none,  1=GoBindata
//   - nc -> NoCache      : 0=false, 1=true
//   - fm -> FuncMap      : 0=false, 1=true
//   - nf -> NoGoFormat   : 0=false, 1=true
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

// writeTemplates create the templates generated file
func writeTemplates(ctx *Context, dir string) error {
	path := filepath.Join(dir, "templates.go")

	buf := new(bytes.Buffer)
	err := ctx.WritePackage(buf)
	if err == nil {
		err = writeFile(path, buf.String())
	}
	return err
}

// writeBindata create a go-bindata file based on the Context
func writeBindata(ctx *Context, dir string) error {
	if ctx.AssetManager == "" {
		return nil
	}
	path := filepath.Join(dir, "bindata.go")
	prefix := filepath.Clean(filepath.Join(dir, ctx.TemplateBaseDir))

	c := bindata.NewConfig()
	c.Debug = ctx.NoCache
	c.Output = path
	c.Prefix = prefix
	c.Package = ctx.PackageName
	c.Input = []bindata.InputConfig{
		bindata.InputConfig{
			Path:      prefix,
			Recursive: true,
		},
	}
	return bindata.Translate(c)
}

// create a FuncMap file
func writeFuncmap(ctx *Context, dir string) error {
	if ctx.FuncMap == "" {
		return nil
	}
	path := filepath.Join(dir, "funcmap.go")

	const text = `package %s
import "html/template"
var %s = template.FuncMap{}
`
	content := fmt.Sprintf(text, ctx.PackageName, ctx.FuncMap)
	return writeFile(path, content)
}

// create a main file
func writeMain(ctx *Context, dir string) error {
	if ctx.PackageName != "main" {
		return nil
	}
	path := filepath.Join(dir, "main.go")
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
	return writeFile(path, text)
}

func execGoRun(dir string) error {
	var out []byte
	var err error

	//cmdline := fmt.Sprintf("go run %s", filepath.Join(dir, "*.go"))
	cmdline := fmt.Sprintf("cd %s && go run *.go", dir)

	cmd := exec.Command("sh", "-c", cmdline)
	out, err = cmd.CombinedOutput()
	sout := string(out)

	if err != nil {
		return errors.New(sout)
	}

	return nil
}

func subtestRun(ctx *Context, folder, root string, t *testing.T) {

	var mapFuncs = map[string]func(*Context, string) error{
		"templates": writeTemplates,
		"funcmap":   writeFuncmap,
		"bindata":   writeBindata,
		"main":      writeMain,
	}
	var numerr int
	dir := filepath.Join(root, folder)

	// create the needed files in the dir
	for title, fn := range mapFuncs {
		if err := fn(ctx, dir); err != nil {
			numerr++
			t.Errorf("%s/%s: %s", folder, title, err.Error())
		}
	}

	// exec "go run *.go"
	if numerr == 0 {
		if err := execGoRun(dir); err != nil {
			t.Errorf("%s/exec %s", folder, err.Error())
		}
	}
}

func TestRun(t *testing.T) {
	if testing.Short() {
		t.Skip("TestRun: skipping test in short mode")
	}

	// setup
	root := setupDirTemplates()
	t.Logf("TestRun: created TempDir %q", root)

	// clean up
	defer func() {
		if t.Failed() {
			t.Logf("TestRun failed: don't delete TempDir %q", root)
		} else {
			os.RemoveAll(root)
			t.Logf("TestRun: deleted TempDir %q", root)
		}
	}()

	// init context constant properties
	ctx := &Context{
		PackageName:     "main",
		TemplateBaseDir: filepath.Join("..", templateBaseDir),
		Pages:           pages,
		Templates:       templates,
	}

	// loops over context parameters
	for _, nocache := range []bool{false, true} {
		ctx.NoCache = nocache

		for _, assetmngr := range []string{"", "go-bindata"} {
			ctx.AssetManager = assetmngr

			for _, funcmap := range []string{"", "funcMap"} {
				ctx.FuncMap = funcmap

				name := ctx2str(ctx)
				t.Run(name, func(t *testing.T) {
					subtestRun(ctx, name, root, t)
				})
			}
		}
	}
}
