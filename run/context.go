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
	// name of the template file used to generate the module
	templateFile = "templates.tmpl"

	defaultPackageName      = "templates"
	defaultPageEnumType     = "PageEnum"
	defaultTemplateEnumType = "templateEnum"
)

type Context struct {
	Debug        bool `toml:"-"`
	FormatOutput bool `toml:"-"`

	// Package name to use in the generated code.
	PackageName string

	// Asset manager to use. Possible values:
	// - none (default)
	// - go-bindata
	AssetManager string

	// Name of the template.FuncMap variable used in template creation.
	// The variable must be defined in another file of the same package
	// (es. "templates/func-map.go").
	// If empty, no funcMap will be used
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

type dataType struct {
	ProgramName      string
	Timestamp        time.Time
	Debug            bool
	FormatOutput     bool
	AssetManager     string
	PackageName      string
	FuncMap          string
	TemplateBaseDir  string
	TemplateEnumType string
	PageEnumType     string

	Pages     []string            // page names (sorted)
	Bases     []string            // base names
	Templates []string            // used template names (sorted)
	PI2BI     []int               // page-index to base-index
	PI2TI     []int               // page-index to template-index
	T2F       map[string][]string // mapping from template name -> (file1, file2, ...)

	assetMngr AssetMngr

	pageEnumPrefix string
	pageEnumSuffix string
}

func nvl(a, b string) string {
	if a == "" {
		return b
	}
	return a
}

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
	t2f, err := resolveIncludes(ctx.Templates, templates.ToSlice())
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

	data := &dataType{
		ProgramName:      "gentmpl",
		Timestamp:        time.Now(),
		Debug:            ctx.Debug,
		FormatOutput:     ctx.FormatOutput,
		AssetManager:     ctx.AssetManager,
		PackageName:      nvl(ctx.PackageName, defaultPackageName),
		TemplateEnumType: nvl(ctx.TemplateEnumType, defaultTemplateEnumType),
		PageEnumType:     nvl(ctx.PageEnumType, defaultPageEnumType),
		FuncMap:          ctx.FuncMap,
		TemplateBaseDir:  ctx.TemplateBaseDir,

		Pages:     pages.ToSlice(),
		Templates: templates.ToSlice(),
		Bases:     bases.ToSlice(),
		PI2TI:     pi2ti,
		PI2BI:     pi2bi,
		T2F:       t2f,

		assetMngr:      assetMngr,
		pageEnumPrefix: ctx.PageEnumPrefix,
		pageEnumSuffix: ctx.PageEnumSuffix,
	}

	return data, nil
}

// PageName returns the PageEnum constant of the page with given name
func (d *dataType) PageName(name string) string {
	return d.pageEnumPrefix + name + d.pageEnumSuffix
}

func (d *dataType) UseGoBindata() bool {
	return d.assetMngr == AssetMngrGoBindata
}

func (d *dataType) PrintInitVars(sFiles, sIdxs string) string {
	// array of files
	files := NewOrderedSetString()
	// mapping from template-index to files-index
	ti2fi := make([][]int, len(d.Templates))

	for tmplIdx, tmplName := range d.Templates {
		fileNames := d.T2F[tmplName]
		fileIdxs := make([]int, len(fileNames))
		for j, fileName := range fileNames {
			fileIdxs[j] = files.Add(fileName)
		}
		ti2fi[tmplIdx] = fileIdxs
	}

	w := new(bytes.Buffer)
	fmt.Fprintln(w, "// files paths")
	fmt.Fprintf(w, "var %s = [...]string{%s}\n", sFiles, astr2str(files.ToSlice()))

	fmt.Fprintln(w, "// indexes of the files of each template")
	fmt.Fprintf(w, "var %s = [...][]%s{\n", sIdxs, uint(files.Len()))
	for tmplIdx, fileIdxs := range ti2fi {
		fmt.Fprintf(w, "  {%s}, // %s\n", aint2str(fileIdxs), d.Templates[tmplIdx])
	}
	fmt.Fprintf(w, "}")
	return w.String()
}

func (ctx *Context) WriteModule(w io.Writer) error {

	var (
		buf bytes.Buffer // A Buffer needs no initialization.
		p   []byte
	)

	// check and prepare context
	data, err := ctx.checkAndPrepare()
	if err != nil {
		return err
	}
	// define the functions available in the template
	templateFuncMap := template.FuncMap{
		"uint":     uint,
		"astr2str": astr2str,
		"aint2str": aint2str,
	}
	// create a new template and parse templateFile into it
	t, err := template.New("").Funcs(templateFuncMap).ParseFiles(templateFile)
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(&buf, "module", data)
	if err != nil {
		return err
	}

	if ctx.FormatOutput {
		// execute format source
		p, err = format.Source(buf.Bytes())
		if err != nil {
			return fmt.Errorf("Formatting source: %s", err.Error())
		}
	} else {
		p = buf.Bytes()
	}

	_, err = w.Write(p)

	return err
}
