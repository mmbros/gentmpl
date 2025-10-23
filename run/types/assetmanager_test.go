package types

import (
	"flag"
	"fmt"
	"testing"

	"github.com/naoina/toml"
)

func TestString(t *testing.T) {

	var testCases = []struct {
		assetMngr AssetManager
		expected  string
	}{

		{AssetManagerNone, "none"},
		{AssetManagerGoBindata, "go-bindata"},
		{AssetManagerGoRice, "go.rice"},
		{100, "AssetManager(100)"},
	}

	for _, tc := range testCases {
		actual := tc.assetMngr.String()
		if actual != tc.expected {
			t.Errorf("AssetManager %d: expected %q, found %q",
				tc.assetMngr, tc.expected, actual)
		}
	}
}

func TestParse(t *testing.T) {

	var testCases = []struct {
		input    string
		expected AssetManager
		ok       bool
	}{
		{"", AssetManagerNone, true},
		{"None", AssetManagerNone, true},
		{"NONE", AssetManagerNone, true},
		{"Go-Bindata", AssetManagerGoBindata, true},
		{"Bindata", AssetManagerGoBindata, true},
		{"Go.Rice", AssetManagerGoRice, true},
		{"Rice", AssetManagerGoRice, true},
		{"Ri.ce", AssetManagerGoRice, false},
	}
	for _, tc := range testCases {
		actual, err := ParseAssetManager(tc.input)

		if tc.ok {
			if err != nil {
				t.Errorf("Unexpected error for input %q: %s", tc.input, err.Error())
			} else if actual != tc.expected {
				t.Errorf("Input %q: expected %q, found %q", tc.input, tc.expected, actual)
			}

		} else {
			if err == nil {
				t.Errorf("Expected error for input %q: found %[2]d (%[2]s)", tc.input, actual)
			}
		}
	}
}

func TestUnmarshal(t *testing.T) {
	var testCases = []struct {
		input    string
		expected AssetManager
		ok       bool
	}{
		{"NO_DECL", AssetManagerNone, true},
		{"", AssetManagerNone, true},
		{"None", AssetManagerNone, true},
		{"NONE", AssetManagerNone, true},
		{"Go-Bindata", AssetManagerGoBindata, true},
		{"Bindata", AssetManagerGoBindata, true},
		{"Go.Rice", AssetManagerGoRice, true},
		{"Rice", AssetManagerGoRice, true},
		{"Ri.ce", AssetManagerGoRice, false},
		{"Embed", AssetManagerEmbed, true},
		{"go:EMBED", AssetManagerEmbed, true},
		{"go_embed", AssetManagerEmbed, false},
	}

	type config struct {
		AssetManager AssetManager
	}
	var cfg config
	var text string

	for i, tc := range testCases {
		// case 0: test toml without asset manager declaration
		if i > 0 {
			text = fmt.Sprintf("asset_manager = %q", tc.input)
		}

		err := toml.Unmarshal([]byte(text), &cfg)
		if tc.ok {
			if err != nil {
				t.Errorf("Unexpected error for input %q: %s", tc.input, err.Error())
			} else if cfg.AssetManager != tc.expected {
				t.Errorf("Input %q: expected %q, found %q", tc.input, tc.expected, cfg.AssetManager)
			}

		} else {
			if err == nil {
				t.Errorf("Expected error for input %q: found %[2]d (%[2]s)", tc.input, cfg.AssetManager)
			}
		}
	}
}

func TestMarshal(t *testing.T) {
	var testCases = []struct {
		input    AssetManager
		expected string
		ok       bool
	}{
		{AssetManagerNone, "none", true},
		{AssetManagerGoBindata, "go-bindata", true},
		{AssetManagerGoRice, "go.rice", true},
		{AssetManager(100), "", false},
		{AssetManagerNone, "none", true},
		{AssetManagerEmbed, "embed", true},
	}

	type config struct {
		AssetManager AssetManager
	}
	var cfg config

	for _, tc := range testCases {
		cfg.AssetManager = tc.input

		b, err := toml.Marshal(cfg)
		if tc.ok {
			if err != nil {
				t.Errorf("Marshal(%s) unexpected error: %s", tc.input.String(), err.Error())
				continue
			}
			actual := string(b)
			expect := fmt.Sprintf("asset_manager = %q\n", tc.expected)
			if actual != expect {
				t.Errorf("Marshal(%s): expected %q, got %q", tc.input.String(), expect, actual)
			}
		} else {
			if err == nil {
				t.Errorf("Marshal(%s) expected error not raised", tc.input.String())
			}

		}
	}
}

func Test_Var(t *testing.T) {
	var testCases = []struct {
		input    string
		expected AssetManager
		ok       bool
	}{
		{"NO_DECL", AssetManagerNone, true},
		{"", AssetManagerNone, true},
		{"None", AssetManagerNone, true},
		{"NONE", AssetManagerNone, true},
		{"Go-Bindata", AssetManagerGoBindata, true},
		{"Bindata", AssetManagerGoBindata, true},
		{"Rice", AssetManagerGoRice, true},
		{"Go.Rice", AssetManagerGoRice, true},
		{"Ri.ce", AssetManagerGoRice, false},
		{"Embed", AssetManagerEmbed, true},
		{"go:EMBED", AssetManagerEmbed, true},
		{"go_embed", AssetManagerEmbed, false},
	}

	for _, tc := range testCases {

		var args []string
		var am AssetManager

		if tc.input != "NO_DECL" {
			args = []string{"-asset-manager", tc.input}
		}

		fs := flag.NewFlagSet("ExampleValue", flag.ContinueOnError)
		fs.Var(&am, "asset-manager", "Asset manager")
		err := fs.Parse(args)

		if tc.ok {
			if err != nil {
				t.Errorf("Unexpected error for input %q: %s", tc.input, err.Error())
			} else if am != tc.expected {
				t.Errorf("Input %q: expected %q, found %q", tc.input, tc.expected, am)
			}

		} else {
			if err == nil {
				t.Errorf("Expected error for input %q: found %[2]d (%[2]s)", tc.input, am)
			}
		}
	}
}
