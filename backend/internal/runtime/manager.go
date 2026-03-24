package runtime

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type ProcessInfo struct {
	ProjectID string
	PID       int
	Command   string
	Container string
}

type Manager struct {
	mu         sync.RWMutex
	containers map[string]string
}

func NewManager() *Manager {
	return &Manager{
		containers: make(map[string]string),
	}
}

func (m *Manager) Start(projectID, workDir, command string, env []string, port int, image string) (*ProcessInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.containers[projectID]; ok {
		return nil, fmt.Errorf("project runtime already running")
	}
	if command == "" {
		return nil, fmt.Errorf("start command is required")
	}
	if port <= 0 {
		return nil, fmt.Errorf("runtime_port must be greater than 0")
	}
	if strings.TrimSpace(image) == "" {
		image = "node:20-alpine"
	}
	containerName := "novex-app-" + sanitize(projectID)

	_, _ = runDocker("rm", "-f", containerName)

	args := []string{"run", "-d", "--name", containerName, "-w", "/app", "-v", workDir + ":/app", "-p", fmt.Sprintf("%d:%d", port, port)}
	for _, item := range env {
		args = append(args, "-e", item)
	}
	args = append(args, image, "sh", "-c", command)
	containerID, err := runDocker(args...)
	if err != nil {
		return nil, err
	}

	m.containers[projectID] = containerName

	return &ProcessInfo{
		ProjectID: projectID,
		PID:       0,
		Command:   command,
		Container: strings.TrimSpace(containerID),
	}, nil
}

func (m *Manager) Stop(projectID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	containerName, ok := m.containers[projectID]
	if !ok {
		return fmt.Errorf("project runtime is not running")
	}
	if _, err := runDocker("rm", "-f", containerName); err != nil {
		return err
	}
	delete(m.containers, projectID)
	return nil
}

func (m *Manager) Restart(projectID, workDir, command string, env []string, port int, image string) (*ProcessInfo, error) {
	_ = m.Stop(projectID)
	return m.Start(projectID, workDir, command, env, port, image)
}

func (m *Manager) ProvisionDatabase(projectID, dbName, dbUser, dbPassword string, dbPort int) (string, error) {
	containerName := "novex-db-" + sanitize(projectID)
	_, _ = runDocker("rm", "-f", containerName)
	_, err := runDocker(
		"run", "-d",
		"--name", containerName,
		"-e", "POSTGRES_DB="+dbName,
		"-e", "POSTGRES_USER="+dbUser,
		"-e", "POSTGRES_PASSWORD="+dbPassword,
		"-p", fmt.Sprintf("%d:5432", dbPort),
		"postgres:16-alpine",
	)
	if err != nil {
		return "", err
	}
	return containerName, nil
}

func (m *Manager) StopDatabase(projectID string) error {
	containerName := "novex-db-" + sanitize(projectID)
	_, err := runDocker("rm", "-f", containerName)
	return err
}

func (m *Manager) DatabaseState(projectID string) string {
	containerName := "novex-db-" + sanitize(projectID)
	out, err := runDocker("inspect", "-f", "{{.State.Status}}", containerName)
	if err != nil {
		return "not_found"
	}
	return strings.TrimSpace(out)
}

func runDocker(args ...string) (string, error) {
	cmd := exec.Command("docker", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("docker %s failed: %s", strings.Join(args, " "), strings.TrimSpace(stderr.String()))
	}
	return out.String(), nil
}

func sanitize(value string) string {
	value = strings.ToLower(value)
	value = strings.ReplaceAll(value, "_", "-")
	value = strings.ReplaceAll(value, "/", "-")
	value = strings.ReplaceAll(value, ".", "-")
	return value
}

func FindFreePort(start int) int {
	// deterministic fallback strategy to avoid net.Listen side effects in locked environments
	if start <= 0 {
		start = 15432
	}
	return start + int(len(strconv.Itoa(start)))
}
