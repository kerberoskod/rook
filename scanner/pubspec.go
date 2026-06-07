package scanner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type PubspecParser struct{}

func (p *PubspecParser) Name() string { return "pubspec" }

func (p *PubspecParser) Glob() string { return "pubspec.yaml" }

func (p *PubspecParser) Parse(path string) ([]Dependency, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var deps []Dependency
	scanner := bufio.NewScanner(f)
	inDeps := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "dependencies:" {
			inDeps = true
			continue
		}
		if inDeps && line != "" && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			if strings.HasPrefix(line, "dev_dependencies:") || strings.HasPrefix(line, "environment:") {
				break
			}
			inDeps = false
			continue
		}
		if !inDeps {
			continue
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if strings.Contains(trimmed, ":") && !strings.HasPrefix(trimmed, "- ") {
			parts := strings.SplitN(trimmed, ":", 2)
			name := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			if val == "" || val == "^" {
				continue
			}
			if name == "sdk" || name == "path" {
				continue
			}
			val = strings.Trim(val, `"'`)
			if strings.HasPrefix(val, "^") {
				val = val[1:]
			}

			if name != "" && val != "" && !strings.Contains(name, " ") {
				deps = append(deps, Dependency{
					Name:    name,
					Version: val,
				})
			}
		}
	}

	return deps, scanner.Err()
}

func pubspecLatest(name string) (string, error) {
	url := fmt.Sprintf("https://pub.dev/api/packages/%s", name)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Latest struct {
			Version string `json:"version"`
		} `json:"latest"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Latest.Version, nil
}
