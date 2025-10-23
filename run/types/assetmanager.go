// Package types implements types used in the gentmpl utility
package types

import (
	"fmt"
	"strings"
)

// AssetManager is the enum type of the AssetManager possible values
type AssetManager uint8

// AssetManager possible values
const (
	AssetManagerNone      AssetManager = iota
	AssetManagerGoBindata              // unsupported
	AssetManagerGoRice                 // unsupported
	AssetManagerEmbed
)

// string representation of AssetManager
var reprAssetManager = [...]string{"none", "go-bindata", "go.rice", "embed"}

// IsNone returns true if AssetManager is None
func (am AssetManager) IsNone() bool { return am == AssetManagerNone }

// IsGoBindata returns true if AssetManager is GoBindata
func (am AssetManager) IsGoBindata() bool { return am == AssetManagerGoBindata }

// IsGoRice returns true if AssetManager is GoRice
func (am AssetManager) IsGoRice() bool { return am == AssetManagerGoRice }

// IsEmbed returns true if AssetManager is Embed
func (am AssetManager) IsEmbed() bool { return am == AssetManagerEmbed }

// ParseAssetManager converts a string to an AssetManager value.
func ParseAssetManager(s string) (AssetManager, error) {
	switch strings.ToLower(s) {
	case "go-embed", "go:embed", "embed":
		return AssetManagerEmbed, nil
	case "go-bindata", "bindata":
		return AssetManagerGoBindata, nil
	case "go.rice", "rice":
		return AssetManagerGoRice, nil
	case "none", "":
		return AssetManagerNone, nil
	}
	return AssetManagerNone, fmt.Errorf("invalid asset manager: %q", s)
}

// String return the string representation of an AssetManager value.
func (am AssetManager) String() string {
	if am >= AssetManager(len(reprAssetManager)) {
		return fmt.Sprint("AssetManager(", int(am), ")")
	}
	return reprAssetManager[am]
}

// Set function for implementing flag.Value interface.
func (am *AssetManager) Set(s string) error {
	if u, err := ParseAssetManager(s); err != nil {
		return err
	} else {
		*am = u
	}
	return nil
}

// UnmarshalText implements the toml.Unmarshaler interface.
func (am *AssetManager) UnmarshalText(data []byte) (err error) {
	// ref: https://godoc.org/github.com/naoina/toml#Unmarshaler
	s := string(data)
	// // s must be a string with left and right double quotes
	// L := len(s)
	// if L < 2 || s[0] != '"' || s[L-1] != '"' {
	// 	return fmt.Errorf("not a valid TOML string: %q", s)
	// }
	// *am, err = ParseAssetManager(s[1 : L-1])

	*am, err = ParseAssetManager(s)
	return
}

// MarshalTOML implements the toml.Marshaler interface.
func (am AssetManager) MarshalTOML() ([]byte, error) {
	// ref: https://godoc.org/github.com/naoina/toml#Marshaler
	if am >= AssetManager(len(reprAssetManager)) {
		return nil, fmt.Errorf("invalid asset manager: %d", am)
	}
	s := fmt.Sprintf("%q", am.String())
	return []byte(s), nil
}
