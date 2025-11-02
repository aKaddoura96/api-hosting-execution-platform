package container

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type Manager struct {
	client *client.Client
}

func NewManager() (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &Manager{client: cli}, nil
}

func (m *Manager) Close() error {
	return m.client.Close()
}

// DeployAPI creates and starts a container for an API
func (m *Manager) DeployAPI(apiID, runtime, codePath string) (string, error) {
	ctx := context.Background()

	// Build Docker image based on runtime
	imageName := m.getImageForRuntime(runtime)
	
	// Ensure image exists
	if err := m.ensureImage(ctx, imageName); err != nil {
		return "", err
	}

	// Create container configuration
	containerName := fmt.Sprintf("api-%s", apiID)
	
	// Prepare volume mount for code
	absCodePath, err := filepath.Abs(codePath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	hostConfig := &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/app/code:ro", filepath.Dir(absCodePath)),
		},
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
		Resources: container.Resources{
			Memory:   512 * 1024 * 1024, // 512MB
			NanoCPUs: 1000000000,         // 1 CPU
		},
	}

	// Expose port for API
	exposedPorts := nat.PortSet{
		"8080/tcp": struct{}{},
	}

	config := &container.Config{
		Image:        imageName,
		ExposedPorts: exposedPorts,
		WorkingDir:   "/app",
		Env: []string{
			fmt.Sprintf("API_ID=%s", apiID),
			"PORT=8080",
		},
	}

	// Create container
	resp, err := m.client.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Start container
	if err := m.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	log.Printf("Container started: %s (ID: %s)", containerName, resp.ID)
	return resp.ID, nil
}

// StopAPI stops and removes a container
func (m *Manager) StopAPI(containerID string) error {
	ctx := context.Background()

	// Stop container
	timeout := 10 // seconds
	if err := m.client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	// Remove container
	if err := m.client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	log.Printf("Container stopped and removed: %s", containerID)
	return nil
}

// GetContainerStatus checks if container is running
func (m *Manager) GetContainerStatus(containerID string) (string, error) {
	ctx := context.Background()

	info, err := m.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", err
	}

	return info.State.Status, nil
}

// GetContainerLogs retrieves logs from a container
func (m *Manager) GetContainerLogs(containerID string) (string, error) {
	ctx := context.Background()

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "100",
	}

	reader, err := m.client.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(logs), nil
}

func (m *Manager) getImageForRuntime(runtime string) string {
	switch runtime {
	case "python":
		return "python:3.11-slim"
	case "nodejs":
		return "node:18-alpine"
	case "go":
		return "golang:1.21-alpine"
	default:
		return "python:3.11-slim"
	}
}

func (m *Manager) ensureImage(ctx context.Context, imageName string) error {
	// Check if image exists
	_, _, err := m.client.ImageInspectWithRaw(ctx, imageName)
	if err == nil {
		return nil // Image exists
	}

	// Pull image
	log.Printf("Pulling Docker image: %s", imageName)
	reader, err := m.client.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close()

	// Wait for pull to complete
	_, err = io.Copy(os.Stdout, reader)
	return err
}
