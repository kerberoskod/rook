package scanner

import "strings"

func normalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimLeft(v, "^~>=< ")
	v = strings.TrimPrefix(v, "v")
	v = strings.Split(v, "+")[0]
	return strings.TrimSpace(v)
}

func splitGroupArtifact(name string) []string {
	return strings.Split(name, ":")
}
