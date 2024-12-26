package tests

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

type TestCaseInstance struct {
	ID                     string
	Product                string
	TestCase               string
	Version                string
	DockerComposeDirectory string
	DockerComposePath      string
	ServiceName            string
	ServicePort            int
}

func (instance *TestCaseInstance) Up() error {
	_, err := instance.runDockerComposeCommand("up", "-d")
	return err
}

func (instance *TestCaseInstance) Down() error {
	_, err := instance.runDockerComposeCommand("down")
	return err
}

func (instance *TestCaseInstance) ServiceIsHealthy() (bool, error) {
	containerName := fmt.Sprintf("verdex-test-%s-%s-1", instance.ID, instance.ServiceName)
	cmd := exec.Command("docker", "inspect", "--format", "{{json .State.Health.Status }}", containerName)
	stdout, err := cmd.Output()

	if string(stdout) == "\"unhealthy\"\n" {
		return false, fmt.Errorf("failed to start Docker container (unhealthy)")
	}

	isHealthy := string(stdout) == "\"healthy\"\n"
	return isHealthy, err
}

func (instance *TestCaseInstance) ServiceOrigin() (string, error) {
	port := fmt.Sprintf("%d", instance.ServicePort)
	stdout, err := instance.runDockerComposeCommand("port", instance.ServiceName, port)
	if err != nil {
		return "", err
	}

	parts := strings.Split(strings.ReplaceAll(string(stdout), "\n", ""), ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid response from docker")
	}

	origin := fmt.Sprintf("http://localhost:%s", parts[1])
	return origin, nil
}

func (instance *TestCaseInstance) Destroy() {
	instance.Down()

	if _, err := os.Stat(instance.DockerComposeDirectory); !os.IsNotExist(err) {
		os.RemoveAll(instance.DockerComposeDirectory)
	}
}

func (instance *TestCaseInstance) runDockerComposeCommand(arg ...string) (string, error) {
	arg = append([]string{"-f", instance.DockerComposePath}, arg...)
	cmd := exec.Command("docker-compose", arg...)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	log.Debug().
		Str("command", "docker-compose "+strings.Join(arg, " ")).
		Str("product", instance.Product).
		Str("test_case", instance.TestCase).
		Msg("Execute command")

	err := cmd.Run()
	if err != nil {
		log.Debug().
			Str("stdout", outb.String()).
			Str("stderr", errb.String()).
			Str("command", "docker-compose "+strings.Join(arg, " ")).
			Msg("Failed to execute command command")
		return outb.String(), err
	}

	return outb.String(), nil
}
