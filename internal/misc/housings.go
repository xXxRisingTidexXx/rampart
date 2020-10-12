package misc

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Housing int

const (
	PrimaryHousing Housing = iota
	SecondaryHousing
)

var housingViews = map[Housing]string{
	PrimaryHousing:   "primary",
	SecondaryHousing: "secondary",
}

func (housing Housing) String() string {
	if view, ok := housingViews[housing]; ok {
		return view
	}
	return "undefined"
}

var viewHousings = map[string]Housing{
	housingViews[PrimaryHousing]:   PrimaryHousing,
	housingViews[SecondaryHousing]: SecondaryHousing,
}

func (housing *Housing) UnmarshalYAML(node *yaml.Node) error {
	view := ""
	if err := node.Decode(&view); err != nil {
		return err
	}
	newHousing, ok := viewHousings[view]
	if !ok {
		return fmt.Errorf("config: housing %s is undefined", view)
	}
	*housing = newHousing
	return nil
}
