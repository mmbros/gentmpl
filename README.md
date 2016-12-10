# gentmpl

gentmpl is an utility that generate a go package for load and render text or
html templates.

It is base on a configuration file that must define two sections:

- `templates`: defines the templates used to render the pages
- `pages`: defines the template and base names to render each page

gentmpl generate a package in which for each page X it is defined a constant
PageX of type PageEnum. To render the page all you have to do it's to call:

    err := PageX.Execute(w, data)

[![GoDoc](https://godoc.org/github.com/mmbros/gentmpl?status.svg)](https://godoc.org/github.com/mmbros/gentmpl)

### Configuration

gentmpl reads from a TOML configuration file the parameters used to generate
the code.

The mandatory informations of the configuration file are the `templates` and
`pages` sections.

#### Templates

The `templates` section defines the templates used to render the Pages.
Each template must have name and an array of string item.
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

#### Pages

The `pages` section defines the pages to render.
Each page must have a name, a template name and optionally a base name.
If defined, the base name will be used in `template.ExecuteTemplate` as the
name of the template. Otherwise will be called `template.Execute`.

Example:
```
[pages]
Inh1 = {template="inh1"}
Inh2 = {template="inh2"}
Pag1 = {template="flat", base="page-1"}
Pag2 = {template="flat", base="page-2"}
Pag3 = {template="flat", base="page-3"}
```


### Generated Package

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


