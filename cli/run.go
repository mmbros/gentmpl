// cli package contains the main entry function of the application.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mmbros/gentmpl/internal/cmdline"
	"github.com/mmbros/gentmpl/internal/config"
	"github.com/mmbros/gentmpl/internal/version"
)

// Run parses the command line arguments and executes the corrisponding command:
//   - PrintHelp: print usage information
//   - PrintVersion: print version information
//   - CreateConfig: generate the package based on the provided configuration parameters
//   - CreatePackage: generate the template package
//
// It returns the code that should be used for os.Exit.
func Run(appName string) int {

	args := cmdline.NewArgs(appName, flag.ExitOnError)

	err := args.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 2
	}

	if args.Help() {
		args.PrintHelp(os.Stdout)
		return 0
	}

	if args.Version() {
		version.PrintVersion(os.Stdout, appName)
		return 0
	}

	if args.GenConfig() {
		err := cmdGenConfig(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 2
		}
		return 0
	}

	// check config file exists
	if _, err := os.Stat(args.Config()); os.IsNotExist(err) {
		// print error  message and hint
		fmt.Fprintf(os.Stderr, "Configuration file %q not found.\n"+
			"Try '%s -h' for more information.\n", args.Config(), appName)
		return 2
	}

	// read config file
	cfg, err := config.Parse(args)
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

// writeOutput apply the fn func to the io.Writer defined by path.
// If path is empty the Stdout will be used;
// else a new file with the give path will be used.
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
func cmdGenPackage(cfg *config.Config) error {
	ctx := cfg.Context
	return writeOutput(cfg.OutputFile, ctx.WritePackage)
}

// cmdGenConfig generate a demo configuration file for the gentmpl tool.
func cmdGenConfig(args *cmdline.Args) error {
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
	cfg, err := config.Unmarshal([]byte(text))
	if err != nil {
		return err
	}

	ctx := cfg.Context
	return writeOutput(args.OutputFile(), ctx.WriteConfig)
}
