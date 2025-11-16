# gentmpl

gentmpl is a command line utility that generates a go package that parse and
render html or text templates.

[![GoDoc](https://godoc.org/github.com/mmbros/gentmpl?status.svg)](https://godoc.org/github.com/mmbros/gentmpl)
[![Go Report Card](https://goreportcard.com/badge/github.com/mmbros/gentmpl)](https://goreportcard.com/report/github.com/mmbros/gentmpl)
[![Go](https://github.com/mmbros/gentmpl/actions/workflows/go.yml/badge.svg)](https://github.com/mmbros/gentmpl/actions/workflows/go.yml)

## Installation

```
go get -u github.com/mmbros/gentmpl
```

## Usage

```
Usage: gentmpl [OPTION]...

gentmpl is an utility that generates a go package for parse and render html or
text templates.

gentmpl reads a configuration file with two mandatory sections:
  - templates: defines the templates used to render the pages
  - pages: defines the template and base names to render each page

gentmpl generates a package that automatically handles the creation of the
templates, loading and parsing the files specified in the configuration.
For each page of name Name gentmpl defines a constant PageName so that
to render the page all you have to do is:
    // initialize the templates
    tmpls := New(funcMap)
    // execute a named template
    err := tmpls.Execute(w, PageName, data)

Options:

  -a value
    	Asset manager for the templates files: "none" or "embed" (default="none")
    	If present, overwrites the "asset_manager" config parameter.
  -b string
    	Base directory of the templates files.
    	If present, overwrites the "template_base_dir" config parameter.
  -c string
    	Configuration file used to generate the package. (default "gentmpl.conf")
  -g	Generate the configuration file instead of the package.
  -h	Show command usage information.
  -n string
    	Package name of the generated file.
    	If present, overwrites the "package_name" config parameter. (default "templates")
  -no-cache
    	Do not cache templates. Ignored if asset manager is used.
    	If present, overwrites the "no_cache" config parameter.
  -no-go-format
    	Do not format the generated code with go/format.
    	If present, overwrites the "no_go_format" config parameter.
  -o string
    	Optional output file for package/config file. If empty stdout will be used.
  -v	Show version informations.

Examples:

  Generate the templates package
    gentmpl -o templates.gen.go

  Generate a demo configuration file
    gentmpl -g -o gentmpl.conf
```

### Description

gentmpl generates a package that automatically handle the creation of the
templates, loading and parsing the files specified in the configuration.

For each page of name _Name_ gentmpl defines a constant `PageName` so that
to render the page all you have to do is:

    // initialize the templates
    tmpls := New(funcMap)
    // execute a named template
    err := tmpls.Execute(w, PageName, data)


In case the `-g` option is given, gentmpl generates a demo configuration file,
instead of the package.

In case the `-v` option is given, gentmpl print version information and exit.

### Examples:

Generate the templates package using only the default configuration file:
```
gentmpl -o templates.gen.go
```

Generate the templates package using the specified configuration file 
and with the passed arguments override:
```
gentmpl -o templates.gen.go -c mytemplates.conf -a none -b /tmp/tmpl -no-go-format -no-cache templates.gen.go
```

Generate a demo configuration file:
```
gentmpl -g -o demo.conf
```

Generate a demo configuration file with arguments override:
```
gentmpl -g -o demo.conf -n views -a embed
```
## Configuration file

gentmpl reads from a TOML configuration file the parameters used to generate
the code.

The mandatory informations of the configuration file are the `templates` and
`pages` sections.

### Templates

The `templates` section defines the templates used to render the pages.
Each template must have a name and a list of string items.
Each string item can be a:
- path of a file to load in the template creation.
- name of another template to include in the current template.

Example:
```
[templates]
flat = ["flat/footer.tmpl", "flat/header.tmpl", "flat/page1.tmpl", "flat/page2and3.tmpl"]
inh1 = ["inhbase", "inheritance/content1.tmpl"]
inh2 = ["inhbase", "inheritance/content2.tmpl"]
inhbase = ["inheritance/base.tmpl"]
```

### Pages

The `pages` section defines the pages to render.  Each page must have a name, a
template attribute that refers to the name of a template defined in the
`templates` section, and optionally a base name.  If defined, the base name
will be used in `template.ExecuteTemplate` as the name of the template.
Otherwise `template.Execute` will be used.

Example:
```
[pages]
Inh1 = {template="inh1"}
Inh2 = {template="inh2"}
Pag1 = {template="flat", base="page-1"}
Pag2 = {template="flat", base="page-2"}
Pag3 = {template="flat", base="page-3"}
```

### Optional configuration parameters

- `asset_manager`: string. Asset manager to use. Possible values: "none"
  (default) |  "embed".

- `no_cache`: bool (default false). Do not cache the templates. A new template
  will be created on every page.Execute.

- `no_go_format`: bool (dafault false). Do not format the generated code with
  go/format.

- `package_name`: string (default "templates"). Package name used in the
  generated code.

- `page_enum_prefix`: string (default "Page"). String used as prefix in the
  PageEnum constants.

- `page_enum_suffix`: string (default ""). String used as suffix in the
  PageEnum constants.

- `page_enum_type`: string (default "PageEnum"). Name of the PageEnum type
  definition.

- `template_base_dir`: string (default ""). Base folder of the templates files.

- `template_enum_type`: string (default "templateEnum"). Name of the
  TemplateEnum type definition.

- `text_template`: bool (default false). Use text/template instead of
  html/template.

## Generated Package

The generated package exports two types: `Templates` and `PageEnum`


### Templates

The `Templates` type contains the cached templates.

The following methods are defined on the `Templates` type:

  - `New(funcMap template.FuncMap) *Templates`: creates a new Templates object using the passed FuncMap to initialize the templates.

  - `Execute(wr io.Writer, page PageEnum, data interface{}) error`: execute the page's template to the
    specified data object.
  
  - `Template(page PageEnum) *template.Template`: returns the template of the given page.

### PageEnum

The `PageEnum` type identifies each page to be rendered.

Example:
```
// PageEnum is the type of the Pages
PageEnum uint8

// PageEnum constants
const (
	PageInh1 PageEnum = iota
	PageInh2k
	PagePag1
	PagePag2
	PagePag3
)

```

The following methods are defined on the `PageEnum` type:

  - `Files() []string`: returns the files used by the page's template
  - `Base() string`: returns the base name used to render the page's template

