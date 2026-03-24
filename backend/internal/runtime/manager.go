package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

type ProcessInfo struct {
	ProjectID string
	PID       int
	Command   string
}

type Manager struct {
	mu        sync.RWMutex
	processes map[string]*exec.Cmd
}

func NewManager() *Manager {
	return &Manager{
		processes: make(map[string]*exec.Cmd),
	}
}

func (m *Manager) Start(projectID, workDir, command string, env []string) (*ProcessInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if existing, ok := m.processes[projectID]; ok && existing.Process != nil {
		return nil, fmt.Errorf("project runtime already running")
	}

	if command == "" {
		return nil, fmt.Errorf("start command is required")
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Dir = filepath.Clean(workDir)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	m.processes[projectID] = cmd

	return &ProcessInfo{
		ProjectID: projectID,
		PID:       cmd.Process.Pid,
		Command:   command,
	}, nil
}

func (m *Manager) Stop(projectID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cmd, ok := m.processes[projectID]
	if !ok || cmd.Process == nil {
		return fmt.Errorf("project runtime is not running")
	}
	if err := cmd.Process.Kill(); err != nil {
		return err
	}
	delete(m.processes, projectID)
	return nil
}

func (m *Manager) Restart(projectID, workDir, command string, env []string) (*ProcessInfo, error) {
	_ = m.Stop(projectID)
	return m.Start(projectID, workDir, command, env)
}
