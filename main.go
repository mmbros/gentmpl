// gentmpl is an utility that generates a go package for parse and render html or
// text templates.
//
// gentmpl reads a configuration file with two mandatory sections:
//   - templates: defines the templates used to render the pages
//   - pages: defines the template and base names to render each page
//
// gentmpl generates a package that automatically handle the creation of the
// templates, loading and parsing the files specified in the configuration.
// Moreover for each page of name Name gentmpl defines a constant PageName so that
// to render the page all you have to do is:
//
//	err := PageName.Execute(w, data)
package main

import (
	"os"

	"github.com/mmbros/gentmpl/cli"
)

func main() {
	os.Exit(cli.Run("gentmpl"))
}
