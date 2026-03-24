package deployer

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Deployer struct {
	BaseDir      string
	DeployDir    string
	UseDocker    bool // Use Docker for builds? (Security)
}

func NewDeployer() *Deployer {
	cwd, _ := os.Getwd()
	return &Deployer{
		BaseDir:   filepath.Join(cwd, "workspace", "builds"),
		DeployDir: filepath.Join(cwd, "workspace", "deployments"),
		UseDocker: true,
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
		// Fallback to local execution (INSECURE)
		out, err := d.executeCommand(buildPath, "sh", "-c", buildCmd)
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

	// Simple copy (recursive)
	// In production, upload to S3 here
	log.Printf("Deploying artifacts from %s to %s", artifactPath, finalPath)
	
	// Use 'cp -r' for simplicity
	_, err := d.executeCommand(buildPath, "cp", "-r", outputDir, finalPath)
	if err != nil {
		return "", fmt.Errorf("deployment copy failed: %v", err)
	}

	return finalPath, nil
}
