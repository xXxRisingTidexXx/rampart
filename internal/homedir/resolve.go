package homedir

import (
	"path/filepath"
	"runtime"
)

func init() {
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("homedir: failed to instantiate the project folder")
	}
	homeDir = filepath.Dir(filepath.Dir(filepath.Dir(filePath)))
}

var homeDir = ""

func Resolve(path string) string {
	return filepath.Join(homeDir, path)
}
