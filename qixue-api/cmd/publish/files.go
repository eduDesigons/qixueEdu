package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// listLocalFiles returns relative (slash) path → absolute path, sorted by rel.
func listLocalFiles(distDir string) (map[string]string, error) {
	absRoot, err := filepath.Abs(distDir)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string)
	err = filepath.WalkDir(absRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(absRoot, path)
		if err != nil {
			return err
		}
		rel = strings.ReplaceAll(rel, "\\", "/")
		out[rel] = path
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func sortedRels(files map[string]string) []string {
	rels := make([]string, 0, len(files))
	for rel := range files {
		rels = append(rels, rel)
	}
	sort.Strings(rels)
	return rels
}
