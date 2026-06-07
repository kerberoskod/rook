package scanner

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type MavenParser struct{}

func (m *MavenParser) Name() string { return "maven" }

func (m *MavenParser) Glob() string { return "pom.xml" }

func (m *MavenParser) Update(path string, deps []Dependency) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := string(data)
	changed := false

	for _, d := range deps {
		if d.Latest == "unknown" || d.Latest == "" {
			continue
		}
		old := fmt.Sprintf(`<version>%s</version>`, d.Version)
		new := fmt.Sprintf(`<version>%s</version>`, d.Latest)
		if strings.Contains(content, old) {
			content = strings.ReplaceAll(content, old, new)
			changed = true
		}
	}

	if !changed {
		return nil
	}

	return os.WriteFile(path, []byte(content), 0644)
}

type pomXML struct {
	XMLName    xml.Name `xml:"project"`
	GroupID    string   `xml:"groupId"`
	ArtifactID string   `xml:"artifactId"`
	Version    string   `xml:"version"`
	Properties map[string]string
	Dependencies []struct {
		GroupID    string `xml:"groupId"`
		ArtifactID string `xml:"artifactId"`
		Version    string `xml:"version"`
	} `xml:"dependencies>dependency"`
}

func (m *MavenParser) Parse(path string) ([]Dependency, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pom pomXML
	if err := xml.Unmarshal(data, &pom); err != nil {
		return nil, err
	}

	var deps []Dependency
	for _, dep := range pom.Dependencies {
		if dep.Version == "" {
			continue
		}
		deps = append(deps, Dependency{
			Name:    fmt.Sprintf("%s:%s", dep.GroupID, dep.ArtifactID),
			Version: dep.Version,
		})
	}

	return deps, nil
}

func mavenLatest(name string) (string, error) {
	// name is groupId:artifactId
	parts := splitGroupArtifact(name)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid maven coordinate: %s", name)
	}

	group := url.QueryEscape(parts[0])
	artifact := url.QueryEscape(parts[1])
	url := fmt.Sprintf("https://search.maven.org/solrsearch/select?q=g:%s+AND+a:%s&rows=1&wt=json",
		group, artifact)

	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Response struct {
			Docs []struct {
				LatestVersion string `json:"latestVersion"`
			} `json:"docs"`
		} `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Response.Docs) == 0 {
		return "", fmt.Errorf("not found: %s", name)
	}
	return result.Response.Docs[0].LatestVersion, nil
}
