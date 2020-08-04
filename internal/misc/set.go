package misc

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

type Set map[string]struct{}

func (set Set) Contains(key string) bool {
	_, ok := set[key]
	return ok
}

func (set *Set) UnmarshalYAML(node *yaml.Node) error {
	keys := make([]string, 0)
	if err := node.Decode(&keys); err != nil {
		return err
	}
	items := make(map[string]struct{}, len(keys))
	for _, element := range keys {
		items[element] = struct{}{}
	}
	*set = items
	return nil
}

func (set Set) String() string {
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	return fmt.Sprintf("[%s]", strings.Join(keys, " "))
}
