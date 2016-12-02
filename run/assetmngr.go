package run

import (
	"fmt"
	"strings"
)

type AssetMngr uint8

const (
	AssetMngrNone AssetMngr = iota
	AssetMngrGoBindata
)

func parseAssetMngr(s string) (AssetMngr, error) {
	switch strings.ToLower(s) {
	case "go-bindata":
		return AssetMngrGoBindata, nil
	case "none", "":
		return AssetMngrNone, nil
	}
	return AssetMngrNone, fmt.Errorf("Invalid asset_manager value: %q", s)
}
