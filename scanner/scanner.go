package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	for _, d := range deps {
		groups[d.File] = append(groups[d.File], d)
	}

	for file, deps := range groups {
		rel := file
		if !filepath.IsAbs(file) {
			rel = filepath.Join(root, file)
		}
		if err := updateFile(rel, deps); err != nil {
			return fmt.Errorf("failed to update %s: %w", rel, err)
		}
	}

	return nil
}

func updateFile(path string, deps []Dependency) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(data)
	for _, d := range deps {
		old := fmt.Sprintf(`"%s": "%s"`, d.Name, d.Version)
		new := fmt.Sprintf(`"%s": "%s"`, d.Name, d.Latest)
		content = strings.ReplaceAll(content, old, new)
	}

	return os.WriteFile(path, []byte(content), 0644)
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
