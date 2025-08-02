package main

import (
	source "go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tifye/flamingo/assert"
	"github.com/tifye/flamingo/compiler"
)

func main() {
	fset := source.NewFileSet()
	compiler.CompileDir("main", fset, ".", func(fi fs.FileInfo) (io.WriteCloser, error) {
		assert.Assert(!fi.IsDir(), "expected to be file")

		dir := filepath.Dir(fi.Name())
		name := strings.TrimSuffix(filepath.Base(fi.Name()), ".flamingo")
		file, err := os.OpenFile(filepath.Join(dir, name+"_flamingo.go"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		return file, nil
	})
}
