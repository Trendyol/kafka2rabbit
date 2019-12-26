package kafka

import (
	"context"
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

func NewContainer(image string) *Container {
	advertisedHost := "172.17.0.1"
	if isRunningOnOSX() {
		advertisedHost = "127.0.0.1"
	}

	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{"2181:2181", "9092:9092"},
		Env: map[string]string{
			"ADVERTISED_HOST": advertisedHost,
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
	return c.ip + ":9092"
}

func (c *Container) Ip() string {
	return c.ip
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
