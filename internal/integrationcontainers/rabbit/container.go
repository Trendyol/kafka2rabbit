package rabbit

import (
	"context"
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"runtime"
)

type Container struct {
	container        testcontainers.Container
	containerRequest testcontainers.ContainerRequest
	address          string
	ip               string
}

const username = "sa"
const password = "Sa123456"

func NewContainer(image string) *Container {
	var exposedPort []string
	if isRunningOnOSX() {
		exposedPort = []string{
			"5672:5672",
			"15672:15672",
		}
	}

	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: exposedPort,
		Env: map[string]string{
			"RABBITMQ_DEFAULT_USER": username,
			"RABBITMQ_DEFAULT_PASS": password,
		},
		WaitingFor: wait.ForLog("started"),
	}

	return &Container{
		containerRequest: req,
	}
}

func (c *Container) Run() {
	var err error
	c.container, err = testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: c.containerRequest,
		Started:          true,
	})
	c.ip, err = c.container.Host(context.Background())
	Expect(err).NotTo(HaveOccurred())
	if isRunningOnOSX() {
		c.ip = "127.0.0.1"
	}
}

func (c *Container) Address() string {
	return fmt.Sprintf("amqp://%s:%s@%s:5672", username, password, c.ip)
}

func (c *Container) Ip() string {
	return c.ip
}

func (c *Container) Stop() error {
	return c.container.Terminate(context.Background())
}

func (c *Container) ForceRemoveAndPrune() {
	if c.container != nil {
		err := c.container.Terminate(context.Background())
		Expect(err).NotTo(HaveOccurred())
	}
}

func isRunningOnOSX() bool {
	return runtime.GOOS == "darwin"
}
