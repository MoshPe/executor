package pkg

import (
	"os"
	"path"
	"strings"
)

type File struct {
	FileName string
	Path string
	IsLoaded bool
}

func(f *File) FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func(f *File) FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func(f *File) IsFileLoaded() bool{
	return f.IsLoaded
}