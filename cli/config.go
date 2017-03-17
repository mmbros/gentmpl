// Package cli provides a CLI UI for the gentmpl command line tool.
package cli

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/mmbros/gentmpl/run"
	"github.com/naoina/toml"
)

const (
	// name of the command line parameters
	clConfig    = "c"
	clDebug     = "d"
	clGenConfig = "g"
	clHelp      = "h"
	clOutput    = "o"
	clPrefix    = "p"

	// default values
	defaultConfigFile = "gentmpl.conf"
	defaultOutputFile = "" // if empty use StdOut
)

type clArgs struct {
	config    string
	debug     bool
	genConfig bool
	help      bool
	output    string
	prefix    string
}

type config struct {
	run.Context
	OutputFile string
}

func unmarshalConfig(data []byte) (*config, error) {
	// parse config file
	var cfg config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// loadConfigFromFile returns the configuration from a configuration file
func loadConfigFromFile(path string) (*config, error) {
	// open config file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// read config file
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return unmarshalConfig(buf)
}

func cmdHelp() {
	fmt.Fprintln(os.Stderr, `Usage: gentmpl [OPTION]...

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
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
Examples:

  Generate the templates package
    genttmpl -c %[1]s -o templates.go

  Generate a demo configuration file
    gentmpl -g -o %[1]s
`, defaultConfigFile)
}
func parseArgs() *clArgs {
	var args clArgs

	// command line arguments
	flag.StringVar(&args.config, clConfig, defaultConfigFile, "Configuration file used to generate the package.")
	flag.StringVar(&args.output, clOutput, defaultOutputFile, "Optional output file for package/config file. If empty stdout will be used.")
	flag.BoolVar(&args.debug, clDebug, false, "Debug mode. Do not cache templates and do not format generated code.")
	flag.BoolVar(&args.help, clHelp, false, "Show command usage information.")
	flag.BoolVar(&args.genConfig, clGenConfig, false, "Generate the configuration file instead of the package.")
	flag.StringVar(&args.prefix, clPrefix, "", "Optional common prefix of the templates files.\n        If present, overwrites the \"template_base_dir\" config parameter.")

	flag.Parse()

	return &args
}

// parseConfig returns the configuration from command line parameters,
// config parameters and defaults
func parseConfig(args *clArgs) (*config, error) {

	// init config from the config file
	cfg, err := loadConfigFromFile(args.config)
	if err != nil {
		return nil, err
	}

	// update config settings with command line parameters and set defaults
	if args.output != defaultOutputFile || cfg.OutputFile == "" {
		cfg.OutputFile = args.output
	}
	if args.debug {
		cfg.NoGoFormat = true
		cfg.NoCache = true
	}
	if args.prefix != "" {
		cfg.TemplateBaseDir = args.prefix
	}

	return cfg, nil
}

// writeOutput apply the fn func to the io.Writer defined by path. If path is
// empty the Stdout will be used; else a new file with the give path will be
// used.
func writeOutput(path string, fn func(io.Writer) error) error {
	var w io.Writer
	if path == "" {
		w = os.Stdout
	} else {
		// create a new file
		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
		if err != nil {
			return fmt.Errorf("Error writing output file: %s", err.Error())
		}
		w = file
		defer file.Close()
	}
	return fn(w)
}

// cmdGenPackage generate the package based on the provided configuration
// parameters.
func cmdGenPackage(cfg *config) error {
	ctx := cfg.Context
	return writeOutput(cfg.OutputFile, ctx.WritePackage)
}

// cmdGenConfig generate a demo configuration file for the gentmpl tool.
func cmdGenConfig(args *clArgs) error {
	const text = `[templates]
flat = ["flat/footer.tmpl", "flat/header.tmpl", "flat/page1.tmpl", "flat/page2and3.tmpl"]
inhbase = ["inheritance/base.tmpl"]
inh1 = ["inhbase", "inheritance/content1.tmpl"]
inh2 = ["inhbase", "inheritance/content2.tmpl"]
[pages]
Pag1 = {template="flat", base="page-1"}
Pag2 = {template="flat", base="page-2"}
Pag3 = {template="flat", base="page-3"}
Inh1 = {template="inh1"}
Inh2 = {template="inh2"}
`
	cfg, err := unmarshalConfig([]byte(text))
	if err != nil {
		return err
	}
	ctx := cfg.Context
	return writeOutput(args.output, ctx.WriteConfig)
}

// Run parses the command line arguments of the gentmpl tool.
//
// If -g option was specified, it generate a demo configuration file for the
// gentmpl utility calling the CreateConfig method of a run.Context object.
//
// Otherwise it reads the configuration file for initialize a run.Context and
// runs its CreatePackage method.
//
// It returns the code that should be used for os.Exit.
func Run() int {
	const msghelp = "Try 'gentmpl -h' for more information."

	args := parseArgs()

	if args.help {
		cmdHelp()
		return 0
	}

	if args.genConfig {
		err := cmdGenConfig(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 2
		}
		return 0
	}

	// check config file exists
	if _, err := os.Stat(args.config); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Configuration file %q not found.\n%s\n", args.config, msghelp)
		return 2
	}

	// read config file
	cfg, err := parseConfig(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 2
	}

	// generate the package
	err = cmdGenPackage(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	return 0
}
