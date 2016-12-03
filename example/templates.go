// Generated by gentmpl; *** DO NOT EDIT ***
// Created: 2016-12-03 17:50:00
// Params: no_cache=false, no_go_format=false, asset_manager="Go-Bindata", func_map="funcMap"

package main

import (
	"html/template"
	"io"
	"path/filepath"
)

// type definitions
type (
	templateEnum uint8
	PageEnum     uint8
)

// PageEnum constants
const templatesLen = 3
const (
	PageInh1 PageEnum = iota
	PageInh2
	PagePag1
	PagePag2
	PagePag3
)

// module variables
var mTemplates [templatesLen]*template.Template

func files2paths(files []string) []string {
	const templatesFolder = "tmpl"
	var path string
	paths := make([]string, len(files))
	for i, file := range files {
		switch {
		case len(file) == 0, file[0] == '.', file[0] == filepath.Separator:
			path = file
		default:
			path = filepath.Join(templatesFolder, file)
		}
		paths[i] = path
	}
	return paths
}

// getTemplateFiles returns the files used by the n-th template
func getTemplateFiles(n templateEnum) []string {
	const (
		// files paths
		files = [...]string{"flat/footer.tmpl", "flat/header.tmpl", "flat/page1.tmpl", "flat/page2and3.tmpl", "inheritance/base.tmpl", "inheritance/content1.tmpl", "inheritance/content2.tmpl"}
		// template-index to array of file-index
		ti2afi = [...][]uint8{
			{0, 1, 2, 3}, // flat
			{4, 5},       // inh1
			{4, 6},       // inh2
		}
	)
	// get the template files indexes
	idxs := ti2afi[n]
	// build the array of files
	astr := make([]string, len(idxs))
	for j, idx := range idxs {
		astr[j] = files[idx]
	}
	return astr
}

// Files returns the files used by the template of the page
func (page PageEnum) Files() []string {
	// from page to template indexes
	var p2t = [...]templateEnum{1, 2, 0, 0, 0}
	// get the template of the page
	t := p2t[page]
	return getTemplateFiles(t)
}

func init() {
	// init base templates
	for t := templateEnum(0); t < templatesLen; t++ {
		files := getTemplateFiles(t)

		// use go-bindata MustAsset func to load templates
		tmpl := template.New(filepath.Base(files[0])).Funcs(funcMap)
		for _, path := range files2paths(files) {
			tmpl.Parse(string(MustAsset(path)))
		}
		mTemplates[t] = tmpl

	}
}

// Template returns the template.Template of the page
func (page PageEnum) Template() *template.Template {
	var idx = [...]templateEnum{1, 2, 0, 0, 0}
	return mTemplates[idx[page]]
}

// Base returns the template name of the page
func (page PageEnum) Base() string {
	const bases = [...]string{"", "page-1", "page-2", "page-3"}

	const pi2bi = [...]PageEnum{0, 0, 1, 2, 3}
	return bases[pi2bi[page]]

}

// Execute applies a parsed page template to the specified data object,
// writing the output to wr.
// If an error occurs executing the template or writing its output, execution
// stops, but partial results may already have been written to the output writer.
// A template may be executed safely in parallel.
func (page PageEnum) Execute(wr io.Writer, data interface{}) error {
	tmpl := page.Template()
	name := page.Base()
	if name != "" {
		return tmpl.ExecuteTemplate(wr, name, data)
	}
	return tmpl.Execute(wr, data)
}

/*
func main(){
	var page PageEnum = PageInh1
	wr := os.Stdout

	if err := page.Execute(wr, nil); err != nil {
		fmt.Print(err)
	}
}
*/
