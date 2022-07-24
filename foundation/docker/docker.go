package docker

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"regexp"
)

// DockerContainer tracks information about the docker container started for tests.
type DockerContainer struct {
	ID   string
	Type string
	Host string
}

var (
	containerIpPortInsightRegexp, _ = regexp.Compile(`{"HostIp":"(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))","HostPort":"(\d{1,6})"}`)
)

const (
	REGEXP_HOST_IP_GROUP_INDEX   = 1
	REGEXP_HOST_PORT_GROUP_INDEX = 5
)

func RunContainer(imageDescriptor string, containerType string, containerPort string, args ...string) (*DockerContainer, error) {
	arg := []string{"run", "-P", "-d"}
	arg = append(arg, args...)
	arg = append(arg, imageDescriptor)

	cmd := exec.Command("docker", arg...)
	var output bytes.Buffer
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("Couldnt start container %s: %w", imageDescriptor, err)
	}

	id := output.String()[:12]

	tmpl := fmt.Sprintf("[{{range $k,$v := (index .NetworkSettings.Ports \"%s/tcp\")}}{{json $v}}{{end}}]", containerPort)
	cmd = exec.Command("docker", "inspect", "-f", tmpl, id)
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("Coudlnt inspect container %s: %w", id, err)
	}

	groups := containerIpPortInsightRegexp.FindStringSubmatch(output.String())

	runningContainer := DockerContainer{
		ID:   id,
		Type: containerType,
		Host: net.JoinHostPort(groups[REGEXP_HOST_IP_GROUP_INDEX], groups[REGEXP_HOST_PORT_GROUP_INDEX]),
	}

	fmt.Printf("- Container with image: %s\n- Started with cid: %s\n- With host: %s", imageDescriptor, runningContainer.ID, runningContainer.Host)

	return &runningContainer, nil
}

func RemoveContainer(id string) error {
	if err := exec.Command("docker", "stop", id).Run(); err != nil {
		return fmt.Errorf("Couldnt stop container: %w", err)
	}
	if err := exec.Command("docker", "rm", id, "-v").Run(); err != nil {
		return fmt.Errorf("Couldnt remove container: %w", err)
	}
	fmt.Println("Container stopped and rmed:", id)

	return nil
}
