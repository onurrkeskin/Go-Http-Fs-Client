package testinghelpers

import "gitlab.com/onurkeskin/go-http-fs-client/foundation/docker"

func StartFS() (*docker.DockerContainer, error) {
	image := "file-server-amd64:1.0"
	port := "8081"
	args := []string{}
	return docker.RunContainer(image, "testing", port, args...)
}

func StopDB(c *docker.DockerContainer) {
	docker.RemoveContainer(c.ID)
}
