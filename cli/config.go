//go:generate stringer -type=AssetMngr
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
	clConfig  = "c"
	clOutput  = "o"
	clPackage = "p"
	clDebug   = "d"

	defaultConfigFile  = "templates.config.toml"
	defaultOutputFile  = "" // if empty use StdOut
	defaultPackageName = "templates"
)

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

// parseConfig returns the configuration from command line parameters,
// config parameters and defaults
func parseConfig() (*config, error) {
	// command line arguments
	pConfigFile := flag.String(clConfig, defaultConfigFile, "Templates configuration file.")
	pOutFile := flag.String(clOutput, defaultOutputFile, "Output file.")
	pPkgName := flag.String(clPackage, defaultPackageName, "Package name to use in the generated code.")
	pDebug := flag.Bool(clDebug, false, "Do not embed the assets, but provide the embedding API. Contents will still be loaded from disk.")
	flag.Parse()

	// init config from the config file
	cfg, err := loadConfigFromFile(*pConfigFile)
	if err != nil {
		return nil, fmt.Errorf("LoadConfig: %s", err.Error())
	}

	// update config settings with command line parameters and set defaults
	if *pOutFile != defaultOutputFile || cfg.OutputFile == "" {
		cfg.OutputFile = *pOutFile
	}
	if *pPkgName != defaultPackageName || cfg.PackageName == "" {
		cfg.PackageName = *pPkgName
	}
	if *pDebug {
		cfg.NoGoFormat = true
		cfg.NoCache = true
	}

	return cfg, nil
}

func (cfg *config) WriteModule() error {
	var w io.Writer

	ctx := cfg.Context

	if cfg.OutputFile == "" {
		w = os.Stdout
	} else {
		// write buffer to output file
		file, err := os.OpenFile(cfg.OutputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
		if err != nil {
			return fmt.Errorf("Error writing output file: %s", err.Error())
		}
		defer file.Close()
		w = file
	}

	return ctx.WriteModule(w)
}

// run parse cmd line, read the config file and generate the package
func Run() int {

	cfg, err := parseConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	err = cfg.WriteModule()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	return 0
}
