package run

import (
	"bytes"
	"strings"
	"testing"
)

var templates1 = map[string][]string{
	"flat":    {"flat/footer.tmpl", "flat/header.tmpl", "flat/page1.tmpl", "flat/page2and3.tmpl"},
	"inhbase": {"inheritance/base.tmpl"},
	"inh1":    {"inhbase", "inheritance/content1.tmpl"},
	"inh2":    {"inhbase", "inheritance/content2.tmpl"},
}
var pages1 = map[string]struct {
	Template string
	Base     string
}{
	"Pag1": {"flat", "page-1"},
	"Pag2": {"flat", "page-2"},
	"Pag3": {"flat", "page-3"},
	"Inh1": {"inh1", ""},
	"Inh2": {"inh2", ""},
}

func errorLike(err error, msg string) bool {
	return strings.Contains(err.Error(), msg)
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
		tmpl[key] = templates1[key]
	}
	ctx = &Context{Pages: pages1, Templates: tmpl}
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
	ctx := &Context{Pages: pages1, Templates: templates1}
	w := new(bytes.Buffer)
	err := ctx.WritePackage(w)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestWriteConfig(t *testing.T) {
	ctx := &Context{Pages: pages1, Templates: templates1}
	w := new(bytes.Buffer)
	err := ctx.WriteConfig(w)
	if err != nil {
		t.Errorf(err.Error())
	}
}
