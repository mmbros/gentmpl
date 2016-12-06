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

	// default values
	defaultConfigFile = "gentmpl.toml"
	defaultOutputFile = "" // if empty use StdOut
)

type clArgs struct {
	config    string
	debug     bool
	genConfig bool
	help      bool
	output    string
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
Utility that generate a go package for load and render text/html templates.

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

	return cfg, nil
}

func genOutput(output string, fn func(io.Writer) error) error {
	var w io.Writer
	if output == "" {
		w = os.Stdout
	} else {
		// create a new file
		file, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
		if err != nil {
			return fmt.Errorf("Error writing output file: %s", err.Error())
		}
		w = file
		defer file.Close()
	}
	return fn(w)
}

func genPackage(cfg *config) error {
	ctx := cfg.Context
	return genOutput(cfg.OutputFile, ctx.WritePackage)
}

func GenConfig(args *clArgs) error {
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
	return genOutput(args.output, ctx.WriteConfig)
}

// run parse cmd line, read the config file and generate the package
func Run() int {
	const msghelp = "Try 'gentmpl -h' for more information."

	args := parseArgs()

	if args.help {
		cmdHelp()
		return 0
	}

	if args.genConfig {
		err := GenConfig(args)
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
	err = genPackage(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	return 0
}
