package misc

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

func NewSet(items []string) Set {
	set := make(Set, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}

type Set map[string]struct{}

func (s Set) Has(key string) bool {
	_, ok := s[key]
	return ok
}

func (s *Set) UnmarshalYAML(node *yaml.Node) error {
	items := make([]string, 0)
	if err := node.Decode(&items); err != nil {
		return err
	}
	*s = NewSet(items)
	return nil
}

func (s Set) String() string {
	keys := make([]string, 0, len(s))
	for key := range s {
		keys = append(keys, key)
	}
	return fmt.Sprintf("[%s]", strings.Join(keys, " "))
}
