package minio

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	rootUser       = "root_user"
	port           = "9000"
	startupTimeout = 1 * time.Minute
)

type Container struct {
	testcontainers.Container
	Host            string
	Port            int
	AccessKeyID     string
	SecretAccessKey string
}

// RunContainer is the entrypoint to the module
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*Container, error) {
	env := make(map[string]string)
	env["MINIO_ROOT_USER"] = rootUser
	env["MINIO_ROOT_PASSWORD"] = generatePassword()

	_, err := nat.NewPort("", port)
	if err != nil {
		return nil, fmt.Errorf("failed to build port: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Image:        "minio/minio:latest",
		Env:          env,
		Cmd:          []string{"server", "/data"},
		ExposedPorts: []string{port + "/tcp"},
		WaitingFor: wait.
			ForExec([]string{"mc", "ready", "local"}).
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
	realPort, err := container.MappedPort(ctx, port)
	if err != nil {
		return nil, fmt.Errorf("failed to get exposed container port: %v", err)
	}

	return &Container{
		Container:       container,
		Host:            host,
		Port:            realPort.Int(),
		AccessKeyID:     genericContainerReq.Env["MINIO_ROOT_USER"],
		SecretAccessKey: genericContainerReq.Env["MINIO_ROOT_PASSWORD"],
	}, nil
}

type Customizer struct {
	AccessKeyID     string
	SecretAccessKey string
}

func (c Customizer) Customize(req *testcontainers.GenericContainerRequest) {
	req.Env["MINIO_ROOT_USER"] = c.AccessKeyID
	req.Env["MINIO_ROOT_PASSWORD"] = c.SecretAccessKey
}

func WithAuthData(accessKeyID, secretAccessKey string) testcontainers.ContainerCustomizer {
	return Customizer{
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
}

func generatePassword() string {
	const (
		letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
		n       = 16
	)

	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret)
}
