package scanner

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type MavenParser struct{}

func (m *MavenParser) Name() string { return "maven" }

func (m *MavenParser) Glob() string { return "pom.xml" }

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

	groupPath := url.PathEscape(stringsReplaceAll(parts[0], ".", "/"))
	artifact := url.PathEscape(parts[1])
	url := fmt.Sprintf("https://search.maven.org/solrsearch/select?q=g:%s+AND+a:%s&rows=1&wt=json",
		groupPath, artifact)

	resp, err := http.Get(url)
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
