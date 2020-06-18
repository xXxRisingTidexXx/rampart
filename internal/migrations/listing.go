package migrations

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"rampart/internal/homedir"
)

func listVersions() ([]*version, error) {
	versionsDir := homedir.Resolve("internal/migrations/versions")
	fileInfos, err := ioutil.ReadDir(versionsDir)
	if err != nil {
		return nil, fmt.Errorf("migrations: failed to list the versions, %v", err)
	}
	for _, fileInfo := range fileInfos {
		log.Debug(fileInfo.Name())
		log.Debug(fileInfo.IsDir())
		log.Debug(fileInfo.Sys())
	}
	return nil, nil
}
