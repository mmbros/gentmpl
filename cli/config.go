//go:generate stringer -type=AssetMngr
package cli

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"strings"

	"github.com/naoina/toml"
)

const (
	// name of the command line parameters
	parametersclOutput = "output"
	clDebug            = "debug"
	clConfig           = "config"
	clPackage          = "pkg"

	defaultConfigFile       = "templates.config.toml"
	defaultPageEnumType     = "PageEnum"
	defaultTemplateEnumType = "templateEnum"
	defaultOutputFile       = ""
	defaultPackageName      = "templates"
	defaultDefaultPageBase  = "base"
	defaultDebug            = false
)

type AssetMngr int

const (
	AssetMngrNone AssetMngr = iota
	AssetMngrGoBindata
)

type Config struct {
	OutputFile   string
	FormatOutput bool
	Debug        bool

	FuncMap      string
	AssetManager string `toml:"omitempty"`
	AssetMngr    int    `toml:"-"` //ignore field in reading config file

	Folder          string
	PackageName     string
	DefaultPageBase string
	PageEnumType    string
	PageEnumPrefix  string
	PageEnumSuffix  string
	Templates       map[string][]string
	Pages           map[string]struct {
		Template string
		Base     string
	}
}

func parseAssetMngr(s string) (int, error) {
	switch strings.ToLower(s) {
	case "go-bindata":
		return AssetMngrGoBindata, nil
	case "none", "":
		return AssetMngrNone, nil
	}
	return AssetMngrNone, fmt.Errorf("Invalid asset_manager value: %q", s)
}

func (cfg *Config) resolveIncludes(astr []string) (map[string][]string, error) {
	return resolveIncludes(cfg.Templates, astr)
}

// loadConfigFromFile returns the configuration from a configuration file
func loadConfigFromFile(path string) (*Config, error) {
	// error function
	exitWithErr := func(err error) (*Config, error) {
		return nil, err
	}

	// open config file
	f, err := os.Open(path)
	if err != nil {
		return exitWithErr(err)
	}
	defer f.Close()
	// read config file
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return exitWithErr(err)
	}
	// parse config file
	var config Config
	if err := toml.Unmarshal(buf, &config); err != nil {
		return exitWithErr(err)
	}

	if am, err := parseAssetMngr(config.AssetManager); err != nil {
		return exitWithErr(err)
	} else {
		config.AssetMngr = am
	}

	return &config, nil
}

// LoadConfig returns the configuration from command line parameters,
// config parameters and defaults
func LoadConfig() (*Config, error) {
	// command line arguments
	pConfigFile := flag.String(clConfig, defaultConfigFile, "Templates configuration file.")
	pOutFile := flag.String(clOutput, defaultOutputFile, "Output file.")
	filepPkgName := flag.String(clPackage, defaultPackageName, "Package name to use in the generated code.")
	pDebug := flag.Bool(clDebug, false, "Do not embed the assets, but provide the embedding API. Contents will still be loaded from disk.")
	flag.Parse()

	// init config from the config file
	config, err := loadConfigFromFile(*pConfigFile)
	if err != nil {
		return nil, fmt.Errorf("LoadConfig: %s", err.Error())
	}

	// update config settings with command line parameters and set defaults
	if *pOutFile != defaultOutputFile || config.OutputFile == "" {
		config.OutputFile = *pOutFile
	}
	if *pPkgName != defaultPackageName || config.PackageName == "" {
		config.PackageName = *pPkgName
	}
	if *pDebug != defaultDebug {
		config.Debug = *pDebug
	}
	if config.PageEnumType == "" {
		config.PageEnumType = defaultPageEnumType
	}
	if config.DefaultPageBase == "" {
		config.DefaultPageBase = defaultDefaultPageBase
	}

	return config, nil
}

func (cfg *Config) WriteModule() error {

	var (
		buf  bytes.Buffer // A Buffer needs no initialization.
		byt  []byte
		err  error
		file *os.File
	)
	ctx := NewContext(cfg)

	//w := os.Stdout
	w := &buf

	// execute write module
	if err = ctx.PrintModule(w); err != nil {
		return fmt.Errorf("Print module: %s", err.Error())
	}

	if cfg.FormatOutput {
		// execute format source
		byt, err = format.Source(buf.Bytes())
		if err != nil {
			return fmt.Errorf("Formatting source: %s", err.Error())
		}
	} else {
		byt = buf.Bytes()
	}

	if cfg.OutputFile != "" {
		// write buffer to output file
		file, err = os.OpenFile(cfg.OutputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
		if err != nil {
			return fmt.Errorf("Error writing output file: %s", err.Error())
		}
		defer file.Close()
	} else {
		file = os.Stdout
	}

	file.Write(byt)
	return nil
}
