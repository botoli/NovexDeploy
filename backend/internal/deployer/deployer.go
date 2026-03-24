package deployer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Deployer struct {
	BaseDir   string
	DeployDir string
	UseDocker bool
}

func NewDeployer() *Deployer {
	cwd, _ := os.Getwd()
	return &Deployer{
		BaseDir:   filepath.Join(cwd, "workspace", "builds"),
		DeployDir: filepath.Join(cwd, "workspace", "deployments"),
		UseDocker: os.Getenv("USE_DOCKER_BUILDS") == "1",
	}
}

// executeCommand runs a command and returns output, or error
func (d *Deployer) executeCommand(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String(), fmt.Errorf("command failed: %s, stderr: %s", err, stderr.String())
	}
	return out.String(), nil
}

func (d *Deployer) executeShellCommand(dir, command string) (string, error) {
	if runtime.GOOS == "windows" {
		return d.executeCommand(dir, "cmd", "/C", command)
	}
	return d.executeCommand(dir, "sh", "-c", command)
}

// BuildProject performs the build step
func (d *Deployer) BuildProject(ctx context.Context, jobID string, repoURL string, branch string, framework string, buildCmd string) (string, error) {
	buildPath := filepath.Join(d.BaseDir, jobID)

	// 1. Clone
	if err := os.MkdirAll(buildPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create build dir: %v", err)
	}

	// Use Git to clone (assuming token is in repoURL or handling via ssh key is separate concern for now)
	log.Printf("Cloning %s (%s) into %s", repoURL, branch, buildPath)
	_, err := d.executeCommand(buildPath, "git", "clone", "--depth", "1", "--branch", branch, repoURL, ".")
	if err != nil {
		return "", fmt.Errorf("git clone failed: %v", err)
	}

	// 2. Build
	log.Printf("Building project %s with %s", jobID, framework)

	// Auto-detect framework if missing
	if framework == "" {
		detected, err := d.DetectFramework(buildPath)
		if err != nil {
			log.Printf("Framework detection error: %v, defaulting to static", err)
			detected = "static"
		}
		framework = detected
		log.Printf("Autodetected framework: %s", framework)

		// Also update buildCmd if missing
		if buildCmd == "" {
			buildCmd = d.GetBuildCommand(framework)
		}
	}

	var outputLog string
	if d.UseDocker {
		// Use Docker for isolation
		// Mount buildPath to /app in container
		// Determine image based on framework
		image := "node:18-alpine" // default
		if framework == "go" {
			image = "golang:1.21-alpine"
		} else if framework == "python" {
			image = "python:3.11-alpine"
		}

		// Run build in container
		// Example: docker run --rm -v /abs/path:/app -w /app node:18-alpine sh -c "npm install && npm run build"

		absPath, _ := filepath.Abs(buildPath)
		dockerArgs := []string{
			"run", "--rm",
			"-v", fmt.Sprintf("%s:/app", absPath),
			"-w", "/app",
			image,
			"sh", "-c", buildCmd,
		}

		out, err := d.executeCommand(buildPath, "docker", dockerArgs...)
		outputLog = out
		if err != nil {
			return outputLog, fmt.Errorf("docker build failed: %v", err)
		}
	} else {
		out, err := d.executeShellCommand(buildPath, buildCmd)
		outputLog = out
		if err != nil {
			return outputLog, fmt.Errorf("local build failed: %v", err)
		}
	}

	return outputLog, nil
}

// DeployArtifacts handles the deployment (copying to final location)
func (d *Deployer) DeployArtifacts(jobID string, outputDir string) (string, error) {
	// Assume buildPath is <BaseDir>/<jobID>
	buildPath := filepath.Join(d.BaseDir, jobID)
	artifactPath := filepath.Join(buildPath, outputDir)

	finalPath := filepath.Join(d.DeployDir, jobID)

	// Verify artifact path exists
	if _, err := os.Stat(artifactPath); os.IsNotExist(err) {
		return "", fmt.Errorf("output directory %s not found in %s", outputDir, buildPath)
	}

	// Move/Copy artifacts
	if err := os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil {
		return "", err
	}

	log.Printf("Deploying artifacts from %s to %s", artifactPath, finalPath)

	if err := os.RemoveAll(finalPath); err != nil {
		return "", fmt.Errorf("failed to cleanup target dir: %w", err)
	}
	if err := copyDir(artifactPath, finalPath); err != nil {
		return "", fmt.Errorf("deployment copy failed: %w", err)
	}

	return finalPath, nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		return copyFile(path, targetPath, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}
