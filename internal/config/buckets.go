package config

import (
	"gopkg.in/yaml.v3"
)

type Buckets []float64

func (buckets *Buckets) UnmarshalYAML(node *yaml.Node) error {
	panic("implement me")
}
