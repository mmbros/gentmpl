package cmdline

import (
	"fmt"
	"io"
)

// PrintHelp prints usage information about the application.
func (a *Args) PrintHelp(w io.Writer) {

	a.fs.SetOutput(w)

	fmt.Fprintf(w, `Usage: %[1]s [OPTION]...

%[1]s is an utility that generates a go package for parse and render html or
text templates.

%[1]s reads a configuration file with two mandatory sections:
  - templates: defines the templates used to render the pages
  - pages: defines the template and base names to render each page

%[1]s generates a package that automatically handle the creation of the
templates, loading and parsing the files specified in the configuration.
Moreover for each page of name Name %[1]s defines a constant PageName so that
to render the page all you have to do is:
  err := PageName.Execute(w, data)

Options:

`, a.appName)

	a.fs.PrintDefaults()

	fmt.Fprintf(w, `
Examples:

  Generate the templates package
    %[1]s -c %[2]s -o templates.go

  Generate a demo configuration file
    %[1]s -g -o %[2]s
`, a.appName, defaultConfigFile(a.appName))
}
