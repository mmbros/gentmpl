// cmdline package implements the command line parameters passed to the application.
package cmdline

import (
	"flag"
)

const (
	// name of the command line parameters
	clBaseDir   = "b"
	clConfig    = "c"
	clDebug     = "d"
	clGenConfig = "g"
	clHelp      = "h"
	clOutput    = "o"
	clVersion   = "v"

	// default values
	defaultOutputFile = "" // if empty use StdOut
)

// Args struct is used to manage the command line parameters.
type Args struct {
	baseDir   string
	config    string
	debug     bool
	genConfig bool
	help      bool
	output    string
	version   bool

	appName string
	fs      *flag.FlagSet
}

func defaultConfigFile(appName string) string { return appName + ".conf" }

// NewArgs creates a new Args struct.
func NewArgs(appName string, errorHandling flag.ErrorHandling) *Args {
	fs := flag.NewFlagSet(appName, errorHandling)

	a := Args{fs: fs, appName: appName}

	fs.StringVar(&a.config, clConfig, defaultConfigFile(appName), "Configuration file used to generate the package.")
	fs.StringVar(&a.output, clOutput, defaultOutputFile, "Optional output file for package/config file. If empty stdout will be used.")
	fs.BoolVar(&a.debug, clDebug, false, "Debug mode. Overwrite configuration setting:\ndo not cache templates, do not use asset manager and do not format generated code.")
	fs.BoolVar(&a.help, clHelp, false, "Show command usage information.")
	fs.BoolVar(&a.genConfig, clGenConfig, false, "Generate the configuration file instead of the package.")
	fs.StringVar(&a.baseDir, clBaseDir, "", "Base directory of the templates files.\nIf present, overwrites the \"template_base_dir\" config parameter.")
	fs.BoolVar(&a.version, clVersion, false, "Show version informations.")

	// 	fs.Var(&a.assetManager, clAssetManager,
	// 		fmt.Sprintf(`Asset manager for the templates files: %q or %q (default=%q)
	// If present, overwrites the "asset_manager" config parameter.`,
	// 			types.AssetManagerNone,
	// 			types.AssetManagerEmbed,
	// 			defaultAssetManager))

	return &a
}

// Parse parses flag definitions from the argument list, which should not
// include the command name.
func (a *Args) Parse(arguments []string) error {
	return a.fs.Parse(arguments)
}

// isFlagPassed checks if flag was provided.
func (a *Args) isFlagPassed(name string) bool {
	found := false
	a.fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// public methods of the Args struct.

// Debug returns true if debug flag was setted.
func (a *Args) Debug() bool { return a.debug }

// GenConfig returns true if gen-config flag was setted.
func (a *Args) GenConfig() bool { return a.genConfig }

// Help returns true if help flag was setted.
func (a *Args) Help() bool { return a.help }

// Version returns true if version flag was setted.
func (a *Args) Version() bool { return a.version }

// Config returns the path of the configuration file.
// If the config flag was not specified, the default <appname>.conf is used.
func (a *Args) Config() string { return a.config }

// OutputFile returns the path of the output file.
func (a *Args) OutputFile() string { return a.output }

// TemplateBaseDir returns the value of the base directory of the templates passed in the command line.
func (a *Args) TemplateBaseDir() string { return a.baseDir }

// IsPassedOutputFile returns true if the Output command line option was passed.
func (a *Args) IsPassedOutputFile() bool { return a.isFlagPassed(clOutput) }

// IsPassedTemplateBaseDir returns true if the TemplateBaseDir command line option was passed.
func (a *Args) IsPassedTemplateBaseDir() bool { return a.isFlagPassed(clBaseDir) }
