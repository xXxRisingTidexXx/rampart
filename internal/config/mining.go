package config

import (
	"fmt"
)

type Mining struct {
	DSNParams            map[string]string `yaml:"dsnParams"`
	DomriaPrimaryMiner   *DomriaMiner      `yaml:"domriaPrimaryMiner"`
	DomriaSecondaryMiner *DomriaMiner      `yaml:"domriaSecondaryMiner"`
}

func (mining *Mining) String() string {
	return fmt.Sprintf(
		"{%v %v %v}",
		mining.DSNParams,
		mining.DomriaPrimaryMiner,
		mining.DomriaSecondaryMiner,
	)
}
