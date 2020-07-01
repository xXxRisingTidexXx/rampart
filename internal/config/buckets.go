package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Buckets []float64

func (buckets *Buckets) UnmarshalYAML(node *yaml.Node) error {
	var floats []float64
	if err := node.Decode(&floats); err == nil {
		*buckets = floats
		return nil
	}
	var linearBuckets linearBuckets
	if err := node.Decode(&linearBuckets); err == nil {
		floats, err = linearBuckets.toFloats()
		if err != nil {
			return err
		}
		*buckets = floats
		return nil
	}
	var exponentialBuckets exponentialBuckets
	if err := node.Decode(&exponentialBuckets); err == nil {
		floats, err = exponentialBuckets.toFloats()
		if err != nil {
			return err
		}
		*buckets = floats
		return nil
	}
	return fmt.Errorf("config: buckets failed to unmarshal")
}
