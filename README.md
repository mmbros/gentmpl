# gentmpl

gentmpl is an utility that generate a go package for load and render text or
html templates.

[![Doc](https://godoc.org/github.com/mmbros/gentmpl)](https://godoc.org/github.com/mmbros/gentmpl)

### Generated Package

The generated package exports an enum type `PageEnum` and a list of constant of
the same type for each page to be rendered.
The following methods are defined on the `PageEnum` type:

- `Execute(io.Writer, interface{}) error`: execute the page's template to the
   specified data object.
- `Base() string`: returns the base name used to render the page's template
- `Template() template.Template`: returns the template
- `Files() []string`: returns the files used by the page's template

### Configuration

gentmpl reads from a configuration file the parameters used to generate the
code.

The mandatory informations of the configuration file are the `templates` and
`pages` sections.

#### Templates

The `templates` section defines the templates used to render the Pages.
Each template must have name and an array of string item.
Each string item can be a:

  - path of a file to load in the template creation. The file path is
      relative to the `template_base_dir` folder.
	    - name of another template to include in the current template.

#### Pages

The `pages` section defines the pages to render.
Each page must have a name, a template name and optionally a base name.
If defined, the base name will be used in `template.ExecuteTemplate` as the
name of the template. Otherwise will be called `template.Execute`.


