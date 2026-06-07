package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNPMParser(t *testing.T) {
	dir := t.TempDir()
	pkg := `{
		"dependencies": {
			"react": "^19.0.0",
			"axios": "^1.7.0"
		},
		"devDependencies": {
			"typescript": "^5.5.0"
		}
	}`
	path := filepath.Join(dir, "package.json")
	require.NoError(t, os.WriteFile(path, []byte(pkg), 0644))

	p := &NPMParser{}
	deps, err := p.Parse(path)
	require.NoError(t, err)
	assert.Len(t, deps, 3)

	names := make(map[string]string)
	for _, d := range deps {
		names[d.Name] = d.Version
	}
	assert.Equal(t, "19.0.0", names["react"])
	assert.Equal(t, "1.7.0", names["axios"])
	assert.Equal(t, "5.5.0", names["typescript"])
}

func TestPipParser(t *testing.T) {
	dir := t.TempDir()
	req := `fastapi>=0.115.0
uvicorn[standard]==0.34.0
httpx==0.28.0
# comment
`
	path := filepath.Join(dir, "requirements.txt")
	require.NoError(t, os.WriteFile(path, []byte(req), 0644))

	p := &PipParser{}
	deps, err := p.Parse(path)
	require.NoError(t, err)
	assert.Len(t, deps, 3)

	names := make(map[string]string)
	for _, d := range deps {
		names[d.Name] = d.Version
	}
	assert.Equal(t, "0.34.0", names["uvicorn"])
	assert.Equal(t, "0.28.0", names["httpx"])
}

func TestGoModParser(t *testing.T) {
	gomod := `module example.com/test

go 1.23

require (
	github.com/spf13/cobra v1.8.0
	github.com/stretchr/testify v1.9.0
)
`
	dir := t.TempDir()
	path := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(path, []byte(gomod), 0644))

	p := &GoModParser{}
	deps, err := p.Parse(path)
	require.NoError(t, err)
	assert.Len(t, deps, 2)
}

func TestCargoParser(t *testing.T) {
	dir := t.TempDir()
	toml := `[package]
name = "test"
version = "0.1.0"

[dependencies]
serde = "1.0.0"
tokio = { version = "1.36", features = ["full"] }
reqwest = { version = "0.12", features = ["json"] }
`
	path := filepath.Join(dir, "Cargo.toml")
	require.NoError(t, os.WriteFile(path, []byte(toml), 0644))

	p := &CargoParser{}
	deps, err := p.Parse(path)
	require.NoError(t, err)
	assert.Len(t, deps, 3)

	names := make(map[string]string)
	for _, d := range deps {
		names[d.Name] = d.Version
	}
	assert.Equal(t, "1.0.0", names["serde"])
	assert.Equal(t, "1.36", names["tokio"])
	assert.Equal(t, "0.12", names["reqwest"])
}

func TestPubspecParser(t *testing.T) {
	dir := t.TempDir()
	yaml := `name: test
description: A test project

dependencies:
  flutter:
    sdk: flutter
  http: ^0.13.0
  provider: ^6.1.0

dev_dependencies:
  flutter_test:
    sdk: flutter
`
	path := filepath.Join(dir, "pubspec.yaml")
	require.NoError(t, os.WriteFile(path, []byte(yaml), 0644))

	p := &PubspecParser{}
	deps, err := p.Parse(path)
	require.NoError(t, err)
	assert.Len(t, deps, 2)

	names := make(map[string]string)
	for _, d := range deps {
		names[d.Name] = d.Version
	}
	assert.Equal(t, "0.13.0", names["http"])
	assert.Equal(t, "6.1.0", names["provider"])
}

func TestNormalizeVersion(t *testing.T) {
	assert.Equal(t, "19.0.0", normalizeVersion("^19.0.0"))
	assert.Equal(t, "1.7.0", normalizeVersion("~1.7.0"))
	assert.Equal(t, "1.0.0", normalizeVersion(">=1.0.0"))
	assert.Equal(t, "2.0.0", normalizeVersion("v2.0.0"))
	assert.Equal(t, "3.0.0", normalizeVersion("  v3.0.0  "))
}
