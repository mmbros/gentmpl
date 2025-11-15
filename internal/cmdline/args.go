// cmdline package implements the command line parameters passed to the application.
package cmdline

import (
	"flag"
	"fmt"

	"github.com/mmbros/gentmpl/run/types"
)

const (
	// name of the command line parameters
	clAssetManager = "a"
	clBaseDir      = "b"
	clConfig       = "c"
	clGenConfig    = "g"
	clHelp         = "h"
	clPkgName      = "n"
	clNoCache      = "no-cache"
	clNoGoFormat   = "no-go-format"
	clOutput       = "o"
	clVersion      = "v"

	// default values
	defaultOutputFile   = "" // if empty use StdOut
	defaultAssetManager = types.AssetManagerNone
	defaultPkgName      = "templates"
)

// Args struct is used to manage the command line parameters.
type Args struct {
	baseDir      string
	config       string
	genConfig    bool
	help         bool
	output       string
	version      bool
	assetManager types.AssetManager
	noCache      bool
	noGoFormat   bool
	pkgName      string

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
	fs.BoolVar(&a.help, clHelp, false, "Show command usage information.")
	fs.BoolVar(&a.genConfig, clGenConfig, false, "Generate the configuration file instead of the package.")
	fs.StringVar(&a.baseDir, clBaseDir, "", "Base directory of the templates files.\nIf present, overwrites the \"template_base_dir\" config parameter.")
	fs.StringVar(&a.pkgName, clPkgName, defaultPkgName, "Package name of the generated file.\nIf present, overwrites the \"package_name\" config parameter.")
	fs.BoolVar(&a.version, clVersion, false, "Show version informations.")

	fs.BoolVar(&a.noCache, clNoCache, false, "Do not cache templates. Ignored if asset manager is used.\nIf present, overwrites the \"no_cache\" config parameter.")
	fs.BoolVar(&a.noGoFormat, clNoGoFormat, false, "Do not format the generated code with go/format.\nIf present, overwrites the \"no_go_format\" config parameter.")

	fs.Var(&a.assetManager, clAssetManager,
		fmt.Sprintf(`Asset manager for the templates files: %q or %q (default=%q)
If present, overwrites the "asset_manager" config parameter.`,
			types.AssetManagerNone,
			types.AssetManagerEmbed,
			defaultAssetManager))

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

// AssetManager the asset manager flag.
// If the flag was not specified, the default asset manager is used.
func (a *Args) AssetManager() types.AssetManager { return a.assetManager }

// Config returns the path of the configuration file.
// If the config flag was not specified, the default <appname>.conf is used.
func (a *Args) Config() string { return a.config }

// GenConfig returns true if GenConfig flag was setted.
func (a *Args) GenConfig() bool { return a.genConfig }

// Help returns true if Help flag was setted.
func (a *Args) Help() bool { return a.help }

// NoCache returns the NoCache flag.
func (a *Args) NoCache() bool { return a.noCache }

// NoGoFormat returns the NoGoFormat flag.
func (a *Args) NoGoFormat() bool { return a.noGoFormat }

// OutputFile returns the path of the output file.
func (a *Args) OutputFile() string { return a.output }

// PkgName returns the PkgName flag.
func (a *Args) PkgName() string { return a.pkgName }

// Version returns true if version flag was setted.
func (a *Args) Version() bool { return a.version }

// TemplateBaseDir returns the value of the base directory of the templates passed in the command line.
func (a *Args) TemplateBaseDir() string { return a.baseDir }

// IsPassedAssetManager returns true if the AssetManager flag was passed.
func (a *Args) IsPassedAssetManager() bool { return a.isFlagPassed(clAssetManager) }

// IsPassedNoCache returns true if the NoCache flag was passed.
func (a *Args) IsPassedNoCache() bool { return a.isFlagPassed(clNoCache) }

// IsPassedNoGoFormat returns true if the NoGoFormat flag was passed.
func (a *Args) IsPassedNoGoFormat() bool { return a.isFlagPassed(clNoGoFormat) }

// IsPassedOutputFile returns true if the OutputFile flag was passed.
func (a *Args) IsPassedOutputFile() bool { return a.isFlagPassed(clOutput) }

// IsPassedTemplateBaseDir returns true if the TemplateBaseDir flag was passed.
func (a *Args) IsPassedTemplateBaseDir() bool { return a.isFlagPassed(clBaseDir) }

// IsPassedPkgName returns true if the PkgName flag was passed.
func (a *Args) IsPassedPkgName() bool { return a.isFlagPassed(clPkgName) }
