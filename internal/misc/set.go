package misc

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

type Set map[string]struct{}

func (s Set) Contains(key string) bool {
	_, ok := s[key]
	return ok
}

func (s *Set) UnmarshalYAML(node *yaml.Node) error {
	keys := make([]string, 0)
	if err := node.Decode(&keys); err != nil {
		return err
	}
	items := make(map[string]struct{}, len(keys))
	for _, element := range keys {
		items[element] = struct{}{}
	}
	*s = items
	return nil
}

func (s Set) String() string {
	keys := make([]string, 0, len(s))
	for key := range s {
		keys = append(keys, key)
	}
	return fmt.Sprintf("[%s]", strings.Join(keys, " "))
}
