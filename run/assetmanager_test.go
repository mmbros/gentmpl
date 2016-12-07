package run

import (
	"fmt"
	"testing"

	"github.com/naoina/toml"
)

func TestString(t *testing.T) {

	var testCases = []struct {
		assetMngr AssetManagerEnum
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
		expected AssetManagerEnum
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
		expected AssetManagerEnum
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
	}

	type config struct {
		AssetManager AssetManagerEnum
	}
	var cfg config
	var text string

	for i, tc := range testCases {
		// case 0: test toml without asset manager declaration
		if i > 0 {
			text = fmt.Sprintf("asset_manager = %q", tc.input)
		}
		t.Logf(text)

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
