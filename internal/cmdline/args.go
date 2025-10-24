package cmdline

import (
	"flag"
)

const (
	// appName = "gentmpl"

	// name of the command line parameters
	clBaseDir   = "b"
	clConfig    = "c"
	clDebug     = "d"
	clGenConfig = "g"
	clHelp      = "h"
	clOutput    = "o"
	clVersion   = "v"
	// clAssetManager = "asset-manager"

	// default values
	// defaultConfigFile   = appName + ".conf"
	defaultOutputFile = "" // if empty use StdOut
	// defaultAssetManager = types.AssetManagerNone
)

type Args struct {
	baseDir   string
	config    string
	debug     bool
	genConfig bool
	help      bool
	output    string
	version   bool
	// assetManager types.AssetManager

	appName string
	fs      *flag.FlagSet
}

func defaultConfigFile(appName string) string { return appName + ".conf" }

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

func (a *Args) Parse(arguments []string) error {
	return a.fs.Parse(arguments)
}

// IsFlagPassed checks if flag was provided
func (a *Args) isFlagPassed(name string) bool {
	found := false
	a.fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func (a *Args) Debug() bool             { return a.debug }
func (a *Args) GenConfig() bool         { return a.genConfig }
func (a *Args) Help() bool              { return a.help }
func (a *Args) Version() bool           { return a.version }
func (a *Args) Config() string          { return a.config }
func (a *Args) OutputFile() string      { return a.output }
func (a *Args) TemplateBaseDir() string { return a.baseDir }

func (a *Args) IsPassedOutputFile() bool      { return a.isFlagPassed(clOutput) }
func (a *Args) IsPassedTemplateBaseDir() bool { return a.isFlagPassed(clBaseDir) }
