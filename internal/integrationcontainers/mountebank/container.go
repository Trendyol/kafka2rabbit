package mountebank

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
	port             int
}

const username = "sa"
const password = "Sa123456"

func NewContainer(image string) *Container {
	var exposedPort []string
	if isRunningOnOSX() {
		exposedPort = []string{
			"2525:2525",
			"9999:9999",
		}
	}

	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: exposedPort,
		Cmd:          []string{"mb", "--mock"},
		WaitingFor:   wait.ForLog("now taking orders"),
	}

	return &Container{
		containerRequest: req,
		port:             9999,
	}
}

func (c *Container) Run() {
	var err error
	c.container, err = testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: c.containerRequest,
		Started:          true,
	})
	Expect(err).NotTo(HaveOccurred())
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

func (c *Container) Port() int {
	return c.port
}

func (c *Container) AdminUri() string {
	return fmt.Sprintf("http://%s:2525", c.ip)
}

func (c *Container) Uri() string {
	return fmt.Sprintf("http://%s:%d", c.ip, c.port)
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
