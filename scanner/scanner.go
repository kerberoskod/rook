package scanner

import (
	"fmt"
	"os"
	"path/filepath"
)

type Dependency struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Latest   string `json:"latest,omitempty"`
	Outdated bool   `json:"outdated,omitempty"`
	File     string `json:"file,omitempty"`
	Manager  string `json:"manager"`
}

type Scanner struct {
	parsers []Parser
}

type Parser interface {
	Name() string
	Glob() string
	Parse(path string) ([]Dependency, error)
	Update(path string, deps []Dependency) error
}

func New() *Scanner {
	return &Scanner{
		parsers: []Parser{
			&NPMParser{},
			&MavenParser{},
			&GoModParser{},
			&PipParser{},
			&CargoParser{},
			&PubspecParser{},
		},
	}
}

func (s *Scanner) Scan(root string) ([]Dependency, error) {
	var all []Dependency

	for _, p := range s.parsers {
		pattern := filepath.Join(root, p.Glob())
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		for _, match := range matches {
			deps, err := p.Parse(match)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to parse %s: %v\n", match, err)
				continue
			}
			for i := range deps {
				deps[i].Manager = p.Name()
				deps[i].File = match
			}
			all = append(all, deps...)
		}
	}

	return all, nil
}

func (s *Scanner) CheckUpdates(deps []Dependency) ([]Dependency, error) {
	result := make([]Dependency, len(deps))

	for i, d := range deps {
		result[i] = d
		latest, err := fetchLatestVersion(d)
		if err != nil {
			result[i].Latest = "unknown"
			result[i].Outdated = false
			continue
		}
		result[i].Latest = latest
		result[i].Outdated = normalizeVersion(d.Version) != normalizeVersion(latest)
	}

	return result, nil
}

func (s *Scanner) ApplyUpdates(root string, deps []Dependency) error {
	groups := make(map[string][]Dependency)
	fileManager := make(map[string]string)

	for _, d := range deps {
		groups[d.File] = append(groups[d.File], d)
	}

	parserByGlob := make(map[string]Parser)
	for _, p := range s.parsers {
		parserByGlob[p.Glob()] = p
	}

	for file, fileDeps := range groups {
		rel := file
		if !filepath.IsAbs(file) {
			rel = filepath.Join(root, file)
		}

		base := filepath.Base(rel)
		parser, ok := parserByGlob[base]
		if !ok {
			// fallback: find parser by suffix match
			fileName := filepath.Base(rel)
			for _, p := range s.parsers {
				if matched, _ := filepath.Match(p.Glob(), fileName); matched {
					parser = p
					ok = true
					break
				}
			}
			if !ok {
				return fmt.Errorf("no parser found for %s", rel)
			}
		}
		fileManager[rel] = parser.Name()
		if err := parser.Update(rel, fileDeps); err != nil {
			return fmt.Errorf("failed to update %s: %w", rel, err)
		}
	}

	return nil
}

func fetchLatestVersion(d Dependency) (string, error) {
	switch d.Manager {
	case "npm":
		return npmLatest(d.Name)
	case "maven":
		return mavenLatest(d.Name)
	case "go":
		return goLatest(d.Name)
	case "pip":
		return pipLatest(d.Name)
	case "cargo":
		return cargoLatest(d.Name)
	case "pubspec":
		return pubspecLatest(d.Name)
	default:
		return d.Version, nil
	}
}
