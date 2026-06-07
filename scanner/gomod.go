package scanner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type GoModParser struct{}

func (g *GoModParser) Name() string { return "go" }

func (g *GoModParser) Glob() string { return "go.mod" }

func (g *GoModParser) Update(path string, deps []Dependency) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(data)
	for _, d := range deps {
		if d.Latest == "unknown" || d.Latest == "" {
			continue
		}
		old := d.Name + " " + d.Version
		new := d.Name + " " + d.Latest
		content = strings.ReplaceAll(content, old, new)
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func (g *GoModParser) Parse(path string) ([]Dependency, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var deps []Dependency
	scanner := bufio.NewScanner(f)
	inRequire := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "require (" {
			inRequire = true
			continue
		}
		if inRequire && line == ")" {
			inRequire = false
			continue
		}
		if strings.HasPrefix(line, "require ") {
			rest := strings.TrimPrefix(line, "require ")
			parts := strings.Fields(rest)
			if len(parts) >= 2 {
				deps = append(deps, Dependency{
					Name:    parts[0],
					Version: parts[1],
				})
			}
			continue
		}
		if inRequire {
			parts := strings.Fields(line)
			if len(parts) >= 2 && !strings.HasPrefix(parts[0], "//") {
				deps = append(deps, Dependency{
					Name:    parts[0],
					Version: parts[1],
				})
			}
		}
	}

	return deps, scanner.Err()
}

func goLatest(name string) (string, error) {
	url := fmt.Sprintf("https://proxy.golang.org/%s/@latest", name)
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Version string `json:"Version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Version, nil
}
