package config

import (
	"fmt"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"gopkg.in/yaml.v3"
)

type Housing string

func (housing *Housing) UnmarshalYAML(node *yaml.Node) error {
	s := ""
	if err := node.Decode(&s); err != nil {
		return err
	}
	if s != misc.HousingPrimary && s != misc.HousingSecondary {
		return fmt.Errorf("config: housing %s is undefined", s)
	}
	*housing = Housing(s)
	return nil
}
