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

func TestNormalizeVersion(t *testing.T) {
	assert.Equal(t, "19.0.0", normalizeVersion("^19.0.0"))
	assert.Equal(t, "1.7.0", normalizeVersion("~1.7.0"))
	assert.Equal(t, "1.0.0", normalizeVersion(">=1.0.0"))
	assert.Equal(t, "2.0.0", normalizeVersion("v2.0.0"))
	assert.Equal(t, "3.0.0", normalizeVersion("  v3.0.0  "))
}
