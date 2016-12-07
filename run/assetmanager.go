package run

import (
	"fmt"
	"strings"
)

// AssetManager enum type definition
type AssetManagerEnum uint8

// AssetManager possible values
const (
	AssetManagerNone AssetManagerEnum = iota
	AssetManagerGoBindata
	AssetManagerGoRice
)

// IsNone returns true if AssetManager is None
func (am AssetManagerEnum) IsNone() bool { return am == AssetManagerNone }

// IsGoBindata returns true if AssetManager is GoBindata
func (am AssetManagerEnum) IsGoBindata() bool { return am == AssetManagerGoBindata }

// IsGoRice returns true if AssetManager is GoRice
func (am AssetManagerEnum) IsGoRice() bool { return am == AssetManagerGoRice }

// ParseAssetManager converts a string to an AssetManager value.
func ParseAssetManager(s string) (AssetManagerEnum, error) {
	switch strings.ToLower(s) {
	case "go-bindata", "bindata":
		return AssetManagerGoBindata, nil
	case "go.rice", "rice":
		return AssetManagerGoRice, nil
	case "none", "":
		return AssetManagerNone, nil
	}
	return AssetManagerNone, fmt.Errorf("Invalid asset manager: %q", s)
}

// String return the string representation of an AssetManager value.
func (am AssetManagerEnum) String() string {
	var repr = [...]string{"none", "go-bindata", "go.rice"}
	if am < 0 || am >= AssetManagerEnum(len(repr)) {
		return fmt.Sprint("AssetManager(", int(am), ")")
	}
	return repr[am]
}

// UnmarshalTOML implements the toml.Unmarshaler interface.
func (am *AssetManagerEnum) UnmarshalTOML(data []byte) (err error) {
	s := string(data)
	// s must be a string with left and right double quotes
	L := len(s)
	if L < 2 || s[0] != '"' || s[L-1] != '"' {
		return fmt.Errorf("Not a valid TOML string: %q", s)
	}
	*am, err = ParseAssetManager(s[1 : L-1])
	return
}
