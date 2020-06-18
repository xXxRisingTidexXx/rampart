package migrations

import (
	"fmt"
	"io/ioutil"
)

type version struct {
	id   int64
	path string
}

func (version *version) load() (string, error) {
	bytes, err := ioutil.ReadFile(version.path)
	if err != nil {
		return "", fmt.Errorf("migrations: version %d failed to load, %v", version.id, err)
	}
	return string(bytes), nil
}

func (version *version) String() string {
	return fmt.Sprintf("{%d %s}", version.id, version.path)
}
