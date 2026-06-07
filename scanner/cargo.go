package scanner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type CargoParser struct{}

func (c *CargoParser) Name() string { return "cargo" }

func (c *CargoParser) Glob() string { return "Cargo.toml" }

func (c *CargoParser) Parse(path string) ([]Dependency, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var deps []Dependency
	scanner := bufio.NewScanner(f)
	inDeps := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "[dependencies]" {
			inDeps = true
			continue
		}
		if strings.HasPrefix(line, "[") && line != "[dependencies]" {
			inDeps = false
			continue
		}
		if !inDeps || line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		val = strings.Trim(val, `"`)

	if strings.HasPrefix(val, "{") {
		inner := strings.Trim(val, "{} ")
		for _, field := range strings.Split(inner, ",") {
			field = strings.TrimSpace(field)
			if strings.HasPrefix(field, "version") {
				kv := strings.SplitN(field, "=", 2)
				if len(kv) == 2 {
					val = strings.Trim(kv[1], ` "'`)
				}
			}
		}
	}

		if name != "" && val != "" {
			deps = append(deps, Dependency{
				Name:    name,
				Version: val,
			})
		}
	}

	return deps, scanner.Err()
}

func cargoLatest(name string) (string, error) {
	url := fmt.Sprintf("https://crates.io/api/v1/crates/%s", name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "rook-cli/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Crate struct {
			MaxVersion string `json:"max_version"`
		} `json:"crate"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Crate.MaxVersion, nil
}
