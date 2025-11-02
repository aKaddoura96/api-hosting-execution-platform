package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type ExecutionRequest struct {
	Code       string                 `json:"code"`
	Runtime    string                 `json:"runtime"`
	Input      map[string]interface{} `json:"input,omitempty"`
	TimeoutSec int                    `json:"timeout_sec,omitempty"`
}

type ExecutionResult struct {
	Output     string                 `json:"output"`
	Error      string                 `json:"error,omitempty"`
	StatusCode int                    `json:"status_code"`
	Duration   int64                  `json:"duration_ms"`
	ExitCode   int                    `json:"exit_code"`
	Result     map[string]interface{} `json:"result,omitempty"`
}

type Executor struct {
	client *client.Client
}

func NewExecutor() (*Executor, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &Executor{client: cli}, nil
}

func (e *Executor) Close() error {
	return e.client.Close()
}

// Execute runs code in a containerized environment
func (e *Executor) Execute(req *ExecutionRequest) (*ExecutionResult, error) {
	startTime := time.Now()

	// Set default timeout
	if req.TimeoutSec == 0 {
		req.TimeoutSec = 30
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(req.TimeoutSec)*time.Second)
	defer cancel()

	// Get runtime configuration
	runtimeConfig, err := e.getRuntimeConfig(req.Runtime)
	if err != nil {
		return nil, err
	}

	// Ensure image exists
	if err := e.ensureImage(ctx, runtimeConfig.Image); err != nil {
		return nil, err
	}

	// Prepare code and input files
	tempDir, err := e.prepareCodeFiles(req, runtimeConfig)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	// Create and run container
	result, err := e.runContainer(ctx, runtimeConfig, tempDir, req.TimeoutSec)
	if err != nil {
		return nil, err
	}

	result.Duration = time.Since(startTime).Milliseconds()
	return result, nil
}

type runtimeConfig struct {
	Image      string
	Extension  string
	EntryPoint []string
	Command    []string
}

func (e *Executor) getRuntimeConfig(runtime string) (*runtimeConfig, error) {
	configs := map[string]*runtimeConfig{
		"python": {
			Image:      "python:3.11-slim",
			Extension:  ".py",
			EntryPoint: []string{"python"},
			Command:    []string{"/app/main.py"},
		},
		"nodejs": {
			Image:      "node:18-alpine",
			Extension:  ".js",
			EntryPoint: []string{"node"},
			Command:    []string{"/app/main.js"},
		},
		"go": {
			Image:      "golang:1.22-alpine",
			Extension:  ".go",
			EntryPoint: []string{"/bin/sh"},
			Command:    []string{"-c", "cd /app && go run main.go"},
		},
	}

	config, ok := configs[strings.ToLower(runtime)]
	if !ok {
		return nil, fmt.Errorf("unsupported runtime: %s", runtime)
	}

	return config, nil
}

func (e *Executor) ensureImage(ctx context.Context, imageName string) error {
	// Check if image exists
	_, _, err := e.client.ImageInspectWithRaw(ctx, imageName)
	if err == nil {
		return nil // Image exists
	}

	// Pull image
	reader, err := e.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close()

	// Wait for pull to complete
	_, err = io.Copy(io.Discard, reader)
	return err
}

func (e *Executor) prepareCodeFiles(req *ExecutionRequest, config *runtimeConfig) (string, error) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "api-exec-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Write code file
	codeFile := filepath.Join(tempDir, "main"+config.Extension)
	if err := os.WriteFile(codeFile, []byte(req.Code), 0644); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to write code file: %w", err)
	}

	// Write input as JSON if provided
	if req.Input != nil && len(req.Input) > 0 {
		inputJSON, err := json.Marshal(req.Input)
		if err != nil {
			os.RemoveAll(tempDir)
			return "", fmt.Errorf("failed to marshal input: %w", err)
		}

		inputFile := filepath.Join(tempDir, "input.json")
		if err := os.WriteFile(inputFile, inputJSON, 0644); err != nil {
			os.RemoveAll(tempDir)
			return "", fmt.Errorf("failed to write input file: %w", err)
		}
	}

	return tempDir, nil
}

func (e *Executor) runContainer(ctx context.Context, config *runtimeConfig, codePath string, timeoutSec int) (*ExecutionResult, error) {
	// Read code and input files
	codeContent, err := os.ReadFile(filepath.Join(codePath, "main"+config.Extension))
	if err != nil {
		return nil, fmt.Errorf("failed to read code file: %w", err)
	}

	// Create inline script based on runtime
	var cmd []string
	switch config.Extension {
	case ".py":
		cmd = []string{"python", "-c", string(codeContent)}
	case ".js":
		cmd = []string{"node", "-e", string(codeContent)}
	case ".go":
		// For Go, we still need file-based execution
		cmd = config.Command
	default:
		cmd = config.Command
	}

	// Create container configuration
	containerConfig := &container.Config{
		Image:      config.Image,
		Cmd:        cmd,
		WorkingDir: "/app",
		Env: []string{
			"PYTHONUNBUFFERED=1",
			"NODE_ENV=production",
		},
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
	}

	// Host configuration with resource limits
	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			Memory:     256 * 1024 * 1024, // 256MB
			NanoCPUs:   500000000,         // 0.5 CPU
			PidsLimit:  newInt64(50),      // Limit processes
		},
		NetworkMode: "none", // Disable network for security
		ReadonlyRootfs: false,
	}

	// Create container
	resp, err := e.client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	// Ensure container is removed after execution
	defer e.client.ContainerRemove(context.Background(), resp.ID, container.RemoveOptions{Force: true})

	// Start container
	if err := e.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	// Wait for container to finish or timeout
	statusCh, errCh := e.client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	
	var exitCode int64
	select {
	case err := <-errCh:
		if err != nil {
			return nil, fmt.Errorf("container wait error: %w", err)
		}
	case status := <-statusCh:
		exitCode = status.StatusCode
	case <-ctx.Done():
		// Timeout - kill the container
		timeout := 5
		e.client.ContainerStop(context.Background(), resp.ID, container.StopOptions{Timeout: &timeout})
		return &ExecutionResult{
			Error:      "Execution timeout exceeded",
			StatusCode: 408,
			ExitCode:   -1,
		}, nil
	}

	// Get container logs
	logs, err := e.getContainerLogs(resp.ID)
	if err != nil {
		return &ExecutionResult{
			Error:      fmt.Sprintf("Failed to retrieve logs: %v", err),
			StatusCode: 500,
			ExitCode:   int(exitCode),
		}, nil
	}

	// Parse result
	result := &ExecutionResult{
		Output:     logs,
		StatusCode: 200,
		ExitCode:   int(exitCode),
	}

	if exitCode != 0 {
		result.Error = "Code execution failed"
		result.StatusCode = 500
	}

	// Try to parse JSON output if it looks like JSON
	logs = strings.TrimSpace(logs)
	if strings.HasPrefix(logs, "{") && strings.HasSuffix(logs, "}") {
		var jsonResult map[string]interface{}
		if err := json.Unmarshal([]byte(logs), &jsonResult); err == nil {
			result.Result = jsonResult
		}
	}

	return result, nil
}

func (e *Executor) getContainerLogs(containerID string) (string, error) {
	ctx := context.Background()

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
		Tail:       "1000",
	}

	reader, err := e.client.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func newInt64(i int64) *int64 {
	return &i
}
