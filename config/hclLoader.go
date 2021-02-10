package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type HclLoader interface {
	LoadHclInPath(path string, incSubFold bool) (string, error)
}

type FileLoader struct {
}

func (f *FileLoader) LoadHclInPath(path string, incSubFold bool) (string, error) {
	var contents string

	if incSubFold {
		var files []string
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(filePath) != ".hcl" {
				return nil
			}
			files = append(files, filePath)
			return nil
		})
		if err != nil {
			panic(err)
		}

		for _, f := range files {
			data, err := ioutil.ReadFile(f)
			if err != nil {
				return "", err
			}
			contents = fmt.Sprintf("%s\n%s", contents, string(data))
		}
	} else {
		var files []os.FileInfo
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return "", err
		}

		for _, f := range files {
			if f.IsDir() {
				continue
			}
			if filepath.Ext(f.Name()) != ".hcl" {
				continue
			}

			data, err := ioutil.ReadFile(filepath.Join(path, f.Name()))
			if err != nil {
				return "", err
			}
			contents = fmt.Sprintf("%s\n%s", contents, string(data))
		}
	}
	return contents, nil
}

type FakeHclLoader struct {
	Hcl string
	Err error
}

func (f *FakeHclLoader) LoadHclInPath(path string, incSubFold bool) (string, error) {
	if f.Err != nil {
		return "", f.Err
	}
	return f.Hcl, nil
}
