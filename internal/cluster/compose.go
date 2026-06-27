package cluster

import "path/filepath"

const project = "caravan"

func composeFile(dir string) string {
	return filepath.Join(dir, "docker-compose.yml")
}
