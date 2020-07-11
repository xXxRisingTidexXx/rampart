package migrations

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"rampart/internal/homedir"
	"sort"
	"strconv"
	"strings"
)

func listVersions() ([]*version, error) {
	versionsDir := homedir.Resolve("internal/migrations/versions")
	fileInfos, err := ioutil.ReadDir(versionsDir)
	if err != nil {
		return nil, fmt.Errorf("migrations: failed to list the versions, %v", err)
	}
	versions := make([]*version, 0, len(fileInfos))
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			name := fileInfo.Name()
			ext := filepath.Ext(name)
			if ext != ".sql" {
				return nil, fmt.Errorf("migrations: got a non-sql file, %s", name)
			}
			id, err := strconv.ParseInt(strings.TrimSuffix(name, ext), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("migrations: failed to extract the id in %s, %v", name, err)
			}
			versions = append(versions, &version{id, filepath.Join(versionsDir, name)})
		}
	}
	sort.Slice(
		versions,
		func(i, j int) bool {
			return versions[i].id < versions[j].id
		},
	)
	return versions, nil
}
