package util

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

type Set struct {
	items map[string]struct{}
}

func (set *Set) Contains(key string) bool {
	_, ok := set.items[key]
	return ok
}

func (set *Set) UnmarshalYAML(node *yaml.Node) error {
	var elements []string
	if err := node.Decode(&elements); err != nil {
		return err
	}
	set.items = make(map[string]struct{}, len(elements))
	for _, element := range elements {
		set.items[element] = struct{}{}
	}
	return nil
}

func (set *Set) String() string {
	keys := make([]string, 0, len(set.items))
	for key := range set.items {
		keys = append(keys, key)
	}
	return fmt.Sprintf("[%s]", strings.Join(keys, " "))
}
