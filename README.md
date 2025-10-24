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

gentmpl generates a package that automatically handle the creation of the
templates, loading and parsing the files specified in the configuration.
Moreover for each page of name Name gentmpl defines a constant PageName so that
to render the page all you have to do is:
  err := PageName.Execute(w, data)

Options:

  -b string
        Base directory of the templates files.
        If present, overwrites the "template_base_dir" config parameter.
  -c string
        Configuration file used to generate the package. (default "gentmpl.conf")
  -d    Debug mode. Overwrite configuration setting:
        do not cache templates, do not use asset manager and do not format generated code.
  -g    Generate the configuration file instead of the package.
  -h    Show command usage information.
  -o string
        Optional output file for package/config file. If empty stdout will be used.
  -v    Show version informations.

Examples:

  Generate the templates package
    gentmpl -c gentmpl.conf -o templates.go

  Generate a demo configuration file
    gentmpl -g -o gentmpl.conf
```

### Description

gentmpl generates a package that automatically handle the creation of the
templates, loading and parsing the files specified in the configuration.
Moreover for each page of name _Name_ gentmpl defines a 
constant `PageName` so that to render the page all you have to do is:

    err := PageName.Execute(w, data)

In case the `-g` option is given, gentmpl generates a demo configuration file,
instead of the package.

In case the `-v` option is given, gentmpl print version information and exit.

### Examples:

Generate the templates package (using the default configuration file):
```
gentmpl -o templates.go
```

Generate the templates package with debug mode using an esplicit configuration
file:
```
gentmpl -d -c tmpl.conf -o tmpl.go
```

Generate a demo configuration file:
```
gentmpl -g -o demo.conf
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
  (default) |  "go-bindata".

- `func_map`: string (default ""). Name of the template.FuncMap variable used
  in template creation. The variable must be defined in another file of the
  same package (ex: "templates/func-map.go"). If empty, no funcMap will be
  used.

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

- `text_remplate`: bool (default false). Use text/template instead of
  html/template.

## Generated Package

The generated package exports an enum type `PageEnum` and a list of constant of
the same type for each page to be rendered.

Example:
```
// PageEnum is the type of the Pages
PageEnum uint8

// PageEnum constants
const (
	PageInh1 PageEnum = iota
	PageInh2
	PagePag1
	PagePag2
	PagePag3
)

```

The following methods are defined on the `PageEnum` type:

  - `Execute(io.Writer, interface{}) error`: execute the page's template to the
    specified data object.
  - `Base() string`: returns the base name used to render the page's template
  - `Template() template.Template`: returns the template
  - `Files() []string`: returns the files used by the page's template

## Development notes

The project layout follows [Writing Go Applications with Reusable Logic](https://npf.io/2016/10/reusable-commands/).

See also:

- [Gorram](https://github.com/natefinch/gorram)

