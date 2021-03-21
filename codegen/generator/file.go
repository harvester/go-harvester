package generator

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
)

func DefaultSourceTree() string {
	paths := strings.Split(os.Getenv("GOPATH"), string(filepath.ListSeparator))
	if len(paths) > 0 && len(paths[0]) > 0 {
		return filepath.Join(paths[0], "src")
	}
	return "./"
}

func prepareDirs(dirs ...string) error {
	for _, dir := range dirs {
		if dir == "" {
			continue
		}

		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, file := range files {
			if strings.HasPrefix(file.Name(), "zz_generated") {
				if err := os.Remove(path.Join(dir, file.Name())); err != nil {
					return errors.Wrapf(err, "failed to delete %s", path.Join(dir, file.Name()))
				}
			}
		}
	}

	return nil
}

func Gofmt(workDir, pkg string) error {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		formatted, err := imports.Process(path, content, &imports.Options{
			Fragment:   true,
			Comments:   true,
			TabIndent:  true,
			TabWidth:   8,
			FormatOnly: true,
		})
		if err != nil {
			return err
		}
		f, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Write(formatted)
		return err
	}

	return filepath.Walk(filepath.Join(workDir, pkg), walkFn)
}
