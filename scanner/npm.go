package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type NPMParser struct{}

func (n *NPMParser) Name() string { return "npm" }

func (n *NPMParser) Glob() string { return "package.json" }

func (n *NPMParser) Update(path string, deps []Dependency) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(data)
	for _, d := range deps {
		if d.Latest == "unknown" || d.Latest == "" {
			continue
		}
		for _, prefix := range []string{"^", "~", ">=", "<=", ">", "<", "=", ""} {
			old := fmt.Sprintf(`"%s": "%s%s"`, d.Name, prefix, d.Version)
			new := fmt.Sprintf(`"%s": "%s%s"`, d.Name, prefix, d.Latest)
			if strings.Contains(content, old) {
				content = strings.ReplaceAll(content, old, new)
				break
			}
		}
	}

	return os.WriteFile(path, []byte(content), 0644)
}

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
			Version: strings.TrimLeft(ver, "^~>=< "),
		})
	}
	for name, ver := range pkg.DevDependencies {
		deps = append(deps, Dependency{
			Name:    name,
			Version: strings.TrimLeft(ver, "^~>=< "),
		})
	}

	return deps, nil
}

func npmLatest(name string) (string, error) {
	url := fmt.Sprintf("https://registry.npmjs.org/%s/latest", name)
	resp, err := httpClient.Get(url)
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
