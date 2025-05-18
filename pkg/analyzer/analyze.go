package analyzer

import (
	"go/parser"
	"go/token"
	"strings"

	"github.com/EmreKb/pgw/pkg/fshelper"
)

type Import struct {
	Path string `json:"path"`
}

type Package struct {
	Path     string     `json:"path"`
	Imports  []*Import  `json:"imports"`
	Packages []*Package `json:"packages"`
}

func Analyze(dir *fshelper.Dir) (*Package, error) {
	pkg, err := analyzeDir(dir)
	if err != nil {
		return nil, err
	}

	for _, subDir := range dir.Directories {
		subPkg, err := Analyze(subDir)
		if err != nil {
			return nil, err
		}

		pkg.Packages = append(pkg.Packages, subPkg)
	}

	return pkg, nil
}

func analyzeDir(dir *fshelper.Dir) (*Package, error) {
	pkg := &Package{
		Path:     dir.Path,
		Imports:  make([]*Import, 0),
		Packages: make([]*Package, 0),
	}

	for _, file := range dir.Files {
		fset := token.NewFileSet()
		astFile, err := parser.ParseFile(fset, file.FullPath, nil, 0)
		if err != nil {
			return nil, err
		}

		for _, importSpec := range astFile.Imports {
			exists := false
			for _, imp := range pkg.Imports {
				if imp.Path == importSpec.Path.Value {
					exists = true
					break
				}
			}
			if !exists {
				// Skip built-in packages
				if strings.HasPrefix(importSpec.Path.Value, "\"") && !strings.Contains(importSpec.Path.Value[1:len(importSpec.Path.Value)-1], ".") {
					continue
				}
				pkg.Imports = append(pkg.Imports, &Import{Path: importSpec.Path.Value})
			}
		}
	}

	return pkg, nil
}
