// Package cli provides a CLI UI for the gentmpl command line tool.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mmbros/gentmpl/internal/version"
	"github.com/mmbros/gentmpl/run"
	"github.com/mmbros/gentmpl/run/types"
	"github.com/pelletier/go-toml/v2"
)

const (
	appName = "gentmpl"

	// name of the command line parameters
	clConfig       = "c"
	clDebug        = "d"
	clGenConfig    = "g"
	clHelp         = "h"
	clOutput       = "o"
	clPrefix       = "p"
	clVersion      = "v"
	clAssetManager = "asset-manager"

	// default values
	defaultConfigFile   = appName + ".conf"
	defaultOutputFile   = "" // if empty use StdOut
	defaultAssetManager = types.AssetManagerNone
)

type cmdlineInfo struct {
	config       string
	debug        bool
	genConfig    bool
	help         bool
	output       string
	prefix       string
	version      bool
	assetManager types.AssetManager

	hasAssetManager bool
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
	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return unmarshalConfig(buf)
}

func cmdHelp(w io.Writer, fs *flag.FlagSet) {

	fs.SetOutput(w)

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
`, appName)

	fs.PrintDefaults()

	fmt.Fprintf(w, `
Examples:

  Generate the templates package
    %[1]s -c %[2]s -o templates.go

  Generate a demo configuration file
    %[1]s -g -o %[2]s
`, appName, defaultConfigFile)
}

func (c *cmdlineInfo) newFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.StringVar(&c.config, clConfig, defaultConfigFile, "Configuration file used to generate the package.")
	fs.StringVar(&c.output, clOutput, defaultOutputFile, "Optional output file for package/config file. If empty stdout will be used.")
	fs.BoolVar(&c.debug, clDebug, false, "Debug mode. Overwrite configuration setting:\ndo not cache templates, do not use asset manager and do not format generated code.")
	fs.BoolVar(&c.help, clHelp, false, "Show command usage information.")
	fs.BoolVar(&c.genConfig, clGenConfig, false, "Generate the configuration file instead of the package.")
	fs.StringVar(&c.prefix, clPrefix, "", "Optional common prefix of the templates files.\nIf present, overwrites the \"template_base_dir\" config parameter.")
	fs.BoolVar(&c.version, clVersion, false, "Show version informations.")

	fs.Var(&c.assetManager, clAssetManager,
		fmt.Sprintf(`Asset manager for the templates files: %q or %q (default=%q)
If present, overwrites the "asset_manager" config parameter.`,
			types.AssetManagerNone,
			types.AssetManagerEmbed,
			defaultAssetManager))

	return fs
}

func (c *cmdlineInfo) newFlagSetAndParse(args []string) (*flag.FlagSet, error) {
	fs := c.newFlagSet()
	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}
	c.hasAssetManager = isFlagPassed(clAssetManager, fs)
	return fs, nil
}

// isFlagPassed checks if flag was provided
func isFlagPassed(name string, fs *flag.FlagSet) bool {
	found := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// parseConfig returns the configuration from command line parameters,
// config parameters and defaults
func parseConfig(args *cmdlineInfo) (*config, error) {

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
		cfg.AssetManager = types.AssetManagerNone
	}

	// cache cannot be disabled if asset manager is used
	if cfg.AssetManager != types.AssetManagerNone {
		cfg.NoCache = false
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
			return fmt.Errorf("error writing output file: %s", err.Error())
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
func cmdGenConfig(args *cmdlineInfo) error {
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
	if args.hasAssetManager {
		cfg.AssetManager = args.assetManager
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

	var clinfo cmdlineInfo

	fs, err := clinfo.newFlagSetAndParse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 2
	}

	if clinfo.help {
		cmdHelp(os.Stdout, fs)
		return 0
	}

	if clinfo.version {
		version.Print(os.Stdout, appName)
		return 0
	}

	if clinfo.genConfig {
		err := cmdGenConfig(&clinfo)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 2
		}
		return 0
	}

	// check config file exists
	if _, err := os.Stat(clinfo.config); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Configuration file %q not found.\n%s\n", clinfo.config, msghelp)
		return 2
	}

	// read config file
	cfg, err := parseConfig(&clinfo)
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
