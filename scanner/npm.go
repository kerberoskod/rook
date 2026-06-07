package scanner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type NPMParser struct{}

func (n *NPMParser) Name() string { return "npm" }

func (n *NPMParser) Glob() string { return "package.json" }

func (n *NPMParser) Parse(path string) ([]Dependency, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	var deps []Dependency
	for name, ver := range pkg.Dependencies {
		deps = append(deps, Dependency{
			Name:    name,
			Version: strings.TrimPrefix(ver, "^"),
		})
	}
	for name, ver := range pkg.DevDependencies {
		deps = append(deps, Dependency{
			Name:    name,
			Version: strings.TrimPrefix(ver, "^"),
		})
	}

	return deps, nil
}

func npmLatest(name string) (string, error) {
	url := fmt.Sprintf("https://registry.npmjs.org/%s/latest", name)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Version, nil
}
