package mailhog

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	startupTimeout = 1 * time.Minute
	smtpPort       = "1025"
	apiPort        = "8025"
)

type Container struct {
	testcontainers.Container
	Host    string
	Port    int
	ApiPort int
}

// RunContainer is the entrypoint to the module
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*Container, error) {
	env := make(map[string]string)
	env["MH_SMTP_BIND_ADDR"] = "0.0.0.0:" + smtpPort
	env["MH_API_BIND_ADDR"] = "0.0.0.0:" + apiPort

	_, err := nat.NewPort("", smtpPort)
	if err != nil {
		return nil, fmt.Errorf("failed to build smtp port: %w", err)
	}
	_, err = nat.NewPort("", apiPort)
	if err != nil {
		return nil, fmt.Errorf("failed to build api port: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Image:        "mailhog/mailhog:latest",
		Env:          env,
		ExposedPorts: []string{smtpPort, apiPort},
		WaitingFor: wait.
			ForListeningPort(smtpPort).
			WithStartupTimeout(startupTimeout),
	}
	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}
	for _, opt := range opts {
		opt.Customize(&genericContainerReq)
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %v", err)
	}
	realSMTPPort, err := container.MappedPort(ctx, smtpPort)
	if err != nil {
		return nil, fmt.Errorf("failed to get exposed container smtp port: %v", err)
	}
	realApiPort, err := container.MappedPort(ctx, apiPort)
	if err != nil {
		return nil, fmt.Errorf("failed to get exposed container api port: %v", err)
	}

	return &Container{
		Container: container,
		Host:      host,
		Port:      realSMTPPort.Int(),
		ApiPort:   realApiPort.Int(),
	}, nil
}

func (c *Container) ApiURI() string {
	return fmt.Sprintf("http://%s:%d/api/v2", c.Host, c.ApiPort)
}
