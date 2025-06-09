package core

import (
	"path/filepath"
	"runtime"
)

func GetPath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return filepath.Dir(file)
}
