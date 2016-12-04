//go:generate go-bindata -pkg run -nometadata context.tmpl
package run

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io"
	"text/template"
	"time"
)

const (
	// path of the template file
	templateFile = "context.tmpl"

	defaultPackageName      = "templates"
	defaultPageEnumType     = "PageEnum"
	defaultPagePrefix       = "Page"
	defaultTemplateEnumType = "templateEnum"
)

// Context is ...
type Context struct {
	// Do not cache the templates.
	// A new template will be created on every page.ExecuteTemplate.
	NoCache bool

	// Do not format the generated code with go/format.
	NoGoFormat bool

	// Package name to use in the generated code.
	PackageName string

	// Asset manager to use. Possible values:
	// - none (default)
	// - go-bindata
	AssetManager string

	// Name of the template.FuncMap variable used in template creation.
	// The variable must be defined in another file of the same package
	// (ex: "templates/func-map.go").
	// If empty, no funcMap will be used.
	FuncMap string

	// Name of the TemplateEnum type definition
	TemplateEnumType string

	// Name of the PageEnum type definition
	PageEnumType string

	// Strings used as prefix and suffix in the PageEnum constants
	// Example:  page "CreateUser", prefix="Page", suffix="" -> PageCreateUser
	PageEnumPrefix string
	PageEnumSuffix string

	// Base folder of the templates files
	TemplateBaseDir string

	// Mapping from template name to files included in the template.
	// The files path is relative to the TemplateBaseDir path.
	Templates map[string][]string

	// Mapping from page name to template name and base values used to render the page
	Pages map[string]struct {
		Template string
		Base     string
	}
}

// dataType contains all the information passed to the template used to
// generate the package in WritePackage
type dataType struct {
	ProgramName      string
	Timestamp        time.Time
	NoCache          bool
	NoGoFormat       bool
	AssetManager     string
	PackageName      string
	FuncMap          string
	TemplateBaseDir  string
	TemplateEnumType string
	PageEnumType     string

	Pages     []string // page names (sorted)
	Bases     []string // base names
	Templates []string // used template names (sorted)
	Files     []string // used files
	PI2BI     []int    // page-index to base-index
	PI2TI     []int    // page-index to template-index
	TI2AFI    [][]int  // template-index to array of file-index

	assetMngr      AssetMngr
	pageEnumPrefix string
	pageEnumSuffix string
}

func nvl(a, b string) string {
	if a == "" {
		return b
	}
	return a
}

// checkAndPrepare check for errors in the Context's parameters.
// If no error is found, returns the dataTaype object created based on the Context
func (ctx *Context) checkAndPrepare() (*dataType, error) {
	var (
		err       error
		assetMngr AssetMngr
	)

	// asset manager
	assetMngr, err = parseAssetMngr(ctx.AssetManager)
	if err != nil {
		return nil, err
	}

	// pages
	if len(ctx.Pages) == 0 {
		return nil, errors.New("No pages found")
	}
	pages := NewOrderedSetString()
	for pageName := range ctx.Pages {
		pages.Add(pageName)
	}
	pages.Sort()

	// templates used by the pages
	templates := NewOrderedSetString()
	for _, pageName := range pages.ToSlice() {
		templateName := ctx.Pages[pageName].Template
		if templateName == "" {
			return nil, fmt.Errorf("Page must have a template: page=%s", pageName)
		}
		_, ok := ctx.Templates[templateName]
		if !ok {
			return nil, fmt.Errorf("Template not found for page: page=%s, template=%s", pageName, templateName)
		}
		templates.Add(templateName)
	}
	templates.Sort()

	// page-index -> template-idx
	// Note: must evaluate after pages.Sort and templates.Sort
	pi2ti := make([]int, pages.Len())
	for pageIdx, pageName := range pages.ToSlice() {
		templateName := ctx.Pages[pageName].Template
		templateIdx, _ := templates.Index(templateName)
		pi2ti[pageIdx] = templateIdx
	}

	// resolve used templates
	// mapping from template name -> (file1, file2, ...)
	t2af, err := resolveIncludes(ctx.Templates, templates.ToSlice())
	if err != nil {
		return nil, err
	}

	// bases
	bases := NewOrderedSetString()
	// mapping from page-index to base-index
	pi2bi := make([]int, pages.Len())

	for pageIdx, pageName := range pages.ToSlice() {
		baseName := ctx.Pages[pageName].Base
		pi2bi[pageIdx] = bases.Add(baseName)
	}

	// files
	files := NewOrderedSetString()
	// mapping from template-index to array of file-index
	ti2afi := make([][]int, templates.Len())

	for tmplIdx, tmplName := range templates.ToSlice() {
		fileNames := t2af[tmplName]
		fileIdxs := make([]int, len(fileNames))
		for j, fileName := range fileNames {
			fileIdxs[j] = files.Add(fileName)
		}
		ti2afi[tmplIdx] = fileIdxs
	}

	data := &dataType{
		ProgramName:      "gentmpl",
		Timestamp:        time.Now(),
		NoCache:          ctx.NoCache,
		NoGoFormat:       ctx.NoGoFormat,
		AssetManager:     ctx.AssetManager,
		PackageName:      nvl(ctx.PackageName, defaultPackageName),
		TemplateEnumType: nvl(ctx.TemplateEnumType, defaultTemplateEnumType),
		PageEnumType:     nvl(ctx.PageEnumType, defaultPageEnumType),
		FuncMap:          ctx.FuncMap,
		TemplateBaseDir:  ctx.TemplateBaseDir,

		Pages:     pages.ToSlice(),
		Templates: templates.ToSlice(),
		Bases:     bases.ToSlice(),
		Files:     files.ToSlice(),
		PI2TI:     pi2ti,
		PI2BI:     pi2bi,
		TI2AFI:    ti2afi,

		assetMngr:      assetMngr,
		pageEnumPrefix: nvl(ctx.PageEnumPrefix, defaultPagePrefix),
		pageEnumSuffix: ctx.PageEnumSuffix,
	}

	return data, nil
}

// PageName returns the PageEnum constant of the page with given name.
func (d *dataType) PageName(name string) string {
	return d.pageEnumPrefix + name + d.pageEnumSuffix
}

func (d *dataType) UseGoBindata() bool {
	return d.assetMngr == AssetMngrGoBindata
}

// getTemplate init the template used to write the package
func getTemplate() *template.Template {
	// define the functions available in the template
	templateFuncMap := template.FuncMap{
		"uint":     uint,
		"astr2str": astr2str,
		"aint2str": aint2str,
	}
	// getTemplate create a new template and parse templateFile into it
	t := template.New("").Funcs(templateFuncMap)
	t.Parse(string(MustAsset(templateFile)))

	return t
}

// WritePackage prints the generated package to writer.
func (ctx *Context) WritePackage(w io.Writer) error {

	var (
		buf bytes.Buffer // A Buffer needs no initialization.
		p   []byte
	)

	// check and prepare context
	data, err := ctx.checkAndPrepare()
	if err != nil {
		return err
	}

	// execute the named template
	t := getTemplate()
	err = t.ExecuteTemplate(&buf, "package", data)
	if err != nil {
		return err
	}

	if ctx.NoGoFormat {
		p = buf.Bytes()
	} else {
		// execute format source
		p, err = format.Source(buf.Bytes())
		if err != nil {
			return fmt.Errorf("Formatting source: %s", err.Error())
		}
	}

	// write bytes to writer
	_, err = w.Write(p)

	return err
}

func (ctx *Context) WriteConfig(w io.Writer) error {
	t := getTemplate()
	return t.ExecuteTemplate(w, "toml", ctx)
}

// Check check for errors in the Context's parameters.
func (ctx *Context) Check() error {
	_, err := ctx.checkAndPrepare()
	return err
}
