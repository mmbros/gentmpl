// config package implements the methods to load the configuration of the tool.
package config

import (
	"io"
	"os"

	"github.com/mmbros/gentmpl/internal/cmdline"
	"github.com/mmbros/gentmpl/run"
	"github.com/mmbros/gentmpl/run/types"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	run.Context
	OutputFile string
}

// Unmarshal creates a new Config from an array of bytes.
func Unmarshal(data []byte) (*Config, error) {
	// parse config file
	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// FromFile creates a new Config loading the specified configuration file.
func FromFile(path string) (*Config, error) {
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
	return Unmarshal(buf)
}

// Parse create a new Config from
//  1. command line parameters
//  2. configuration file
//  3. default values
func Parse(args *cmdline.Args) (*Config, error) {

	// init config from the config file
	cfg, err := FromFile(args.Config())
	if err != nil {
		return nil, err
	}

	// update config settings with command line parameters and set defaults
	if args.IsPassedOutputFile() {
		cfg.OutputFile = args.OutputFile()
	}

	if args.IsPassedTemplateBaseDir() {
		cfg.TemplateBaseDir = args.TemplateBaseDir()
	}

	if args.Debug() {
		cfg.NoGoFormat = true
		cfg.NoCache = true
		cfg.AssetManager = types.AssetManagerNone
	}

	// cache cannot be disabled if asset manager is used
	if cfg.AssetManager != types.AssetManagerNone {
		cfg.NoCache = false
	}

	return cfg, nil
}
