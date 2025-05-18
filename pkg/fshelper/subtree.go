package fshelper

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var ErrNotADir = errors.New("not a directory")

var IgnoreDirs = []string{".git", ".idea", ".vscode", ".DS_Store", ".env", ".env.local", ".env.development.local", ".env.test.local", ".env.production.local"}

type File struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
}

type Dir struct {
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	FullPath    string  `json:"full_path"`
	Directories []*Dir  `json:"directories"`
	Files       []*File `json:"files"`
}

func SubTree(path string) (*Dir, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		return nil, ErrNotADir
	}

	baseName := filepath.Base(absPath)
	dir := &Dir{
		Name:        baseName,
		Path:        baseName,
		FullPath:    absPath,
		Directories: make([]*Dir, 0),
		Files:       make([]*File, 0),
	}

	files, err := os.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			if slices.Contains(IgnoreDirs, file.Name()) {
				continue
			}
			childDir, err := SubTree(filepath.Join(absPath, file.Name()))
			if err != nil {
				return nil, err
			}
			// Update child paths to be relative to base
			childDir.Path = filepath.Join(baseName, childDir.Path)
			dir.Directories = append(dir.Directories, childDir)

		} else {
			if !strings.HasSuffix(file.Name(), ".go") {
				continue
			}

			dir.Files = append(dir.Files, &File{
				Name:     file.Name(),
				Path:     filepath.Join(baseName, file.Name()),
				FullPath: filepath.Join(absPath, file.Name()),
			})
		}
	}

	return dir, nil
}
