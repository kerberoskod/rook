package scanner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type PipParser struct{}

func (p *PipParser) Name() string { return "pip" }

func (p *PipParser) Glob() string { return "requirements.txt" }

func (p *PipParser) Update(path string, deps []Dependency) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, d := range deps {
		if d.Latest == "unknown" || d.Latest == "" {
			continue
		}
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if !strings.HasPrefix(trimmed, d.Name) {
				continue
			}
			var sep string
			for _, candidate := range []string{"==", ">=", "<=", "!=", "~="} {
				if strings.Contains(trimmed, d.Name+candidate) {
					sep = candidate
					break
				}
			}
			if sep == "" {
				continue
			}
			prefix := d.Name + sep
			idx := strings.Index(trimmed, prefix) + len(prefix)
			rest := trimmed[idx:]
			endIdx := strings.IndexAny(rest, " ;#\t")
			if endIdx < 0 {
				endIdx = len(rest)
			}
			before := line[:strings.Index(line, trimmed)+idx]
			after := rest[endIdx:]
			lines[i] = before + d.Latest + after
		}
	}

	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}

func (p *PipParser) Parse(path string) ([]Dependency, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var deps []Dependency
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "-") {
			continue
		}

		name := line
		version := "*"
		for _, sep := range []string{"==", ">=", "<=", "!=", "~="} {
			if idx := strings.Index(line, sep); idx >= 0 {
				name = strings.TrimSpace(line[:idx])
				version = strings.TrimSpace(line[idx+len(sep):])
				break
			}
		}
		if idx := strings.Index(name, "["); idx >= 0 {
			name = name[:idx]
		}

		if name != "" {
			deps = append(deps, Dependency{
				Name:    name,
				Version: version,
			})
		}
	}

	return deps, scanner.Err()
}

func pipLatest(name string) (string, error) {
	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", name)
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Info.Version, nil
}
