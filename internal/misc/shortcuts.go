package misc

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		panic("misc: failed to instantiate the root folder")
	}
	rootDir = filepath.Dir(filepath.Dir(filepath.Dir(filePath)))
}

var rootDir = ""

func ResolvePath(path string) string {
	return filepath.Join(rootDir, path)
}

func GetEnv(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("misc: failed to find the env %s", key)
}
