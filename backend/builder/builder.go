package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"localVercel/models"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Builder struct {
	BuildsDir string
}

func NewBuilder() *Builder {
	return &Builder{
		BuildsDir: "./builds",
	}
}

func (b *Builder) DetectFramework(files []string) string {
	for _, file := range files {
		switch {
		case strings.Contains(file, "package.json"):
			// Проверяем зависимости для определения фронтенд фреймворка
			return b.detectNodeFramework(file)
		case strings.Contains(file, "go.mod"):
			return "go"
		case strings.Contains(file, "requirements.txt"):
			return "python"
		case strings.Contains(file, "Gemfile"):
			return "ruby"
		case strings.Contains(file, "Cargo.toml"):
			return "rust"
		case strings.Contains(file, "index.html"):
			return "static"
		}
	}
	return "static"
}

func (b *Builder) detectNodeFramework(packageJSONPath string) string {
	data, err := ioutil.ReadFile(packageJSONPath)
	if err != nil {
		return "node"
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
		Scripts         map[string]string `json:"scripts"`
	}

	if err := json.Unmarshal(data, &pkg); err != nil {
		return "node"
	}

	// Объединяем все зависимости
	allDeps := make(map[string]string)
	for k, v := range pkg.Dependencies {
		allDeps[k] = v
	}
	for k, v := range pkg.DevDependencies {
		allDeps[k] = v
	}

	// Определяем фреймворк
	if _, ok := allDeps["next"]; ok {
		return "nextjs"
	}
	if _, ok := allDeps["gatsby"]; ok {
		return "gatsby"
	}
	if _, ok := allDeps["vue"]; ok || strings.Contains(pkg.Scripts["build"], "vue") {
		return "vue"
	}
	if _, ok := allDeps["@angular/core"]; ok {
		return "angular"
	}
	if _, ok := allDeps["react"]; ok || strings.Contains(pkg.Scripts["build"], "react") {
		return "react"
	}
	if _, ok := allDeps["express"]; ok {
		return "express"
	}

	return "node"
}

func (b *Builder) GetBuildCommand(framework string) string {
	commands := map[string]string{
		"react":    "npm run build",
		"vue":      "npm run build",
		"angular":  "ng build --prod",
		"nextjs":   "npm run build",
		"gatsby":   "gatsby build",
		"node":     "npm ci && npm run build",
		"express":  "npm ci",
		"go":       "go build -o app",
		"python":   "pip install -r requirements.txt",
		"static":   "echo 'Static site - no build needed'",
		"ruby":     "bundle install",
		"rust":     "cargo build --release",
	}
	
	if cmd, ok := commands[framework]; ok {
		return cmd
	}
	return "echo 'No build command specified'"
}

func (b *Builder) GetOutputDir(framework string) string {
	dirs := map[string]string{
		"react":    "build",
		"vue":      "dist",
		"angular":  "dist",
		"nextjs":   "out",
		"gatsby":   "public",
		"node":     ".",
		"express":  ".",
		"go":       ".",
		"python":   ".",
		"static":   ".",
		"ruby":     ".",
		"rust":     "target/release",
	}
	
	if dir, ok := dirs[framework]; ok {
		return dir
	}
	return "."
}

func (b *Builder) BuildProject(config models.BuildConfig) (*models.BuildResult, error) {
	buildID := fmt.Sprintf("build_%d", time.Now().UnixNano())
	buildPath := filepath.Join(b.BuildsDir, buildID)
	
	// Создаем директорию для билда
	if err := os.MkdirAll(buildPath, 0755); err != nil {
		return nil, err
	}

	result := &models.BuildResult{
		ID:            buildID,
		ProjectID:     config.ProjectID,
		Status:        "building",
		CommitSHA:     config.CommitSHA,
		CommitMessage: config.CommitMessage,
		Branch:        config.Branch,
		StartedAt:     time.Now(),
	}

	// Клонируем репозиторий
	clonePath := filepath.Join(buildPath, "repo")
	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", config.Branch, config.ProjectID, clonePath)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		result.Status = "failed"
		result.Logs = stderr.String()
		result.CompletedAt = time.Now()
		result.Duration = int(time.Since(result.StartedAt).Seconds())
		return result, fmt.Errorf("git clone failed: %v", err)
	}

	// Определяем фреймворк если не указан
	if config.Framework == "" {
		files, _ := filepath.Glob(filepath.Join(clonePath, "*"))
		config.Framework = b.DetectFramework(files)
	}

	// Устанавливаем команду билда если не указана
	if config.BuildCommand == "" {
		config.BuildCommand = b.GetBuildCommand(config.Framework)
	}

	// Устанавливаем выходную директорию если не указана
	if config.OutputDir == "" {
		config.OutputDir = b.GetOutputDir(config.Framework)
	}

	// Логируем информацию о билде
	buildLog := fmt.Sprintf("Framework: %s\nBuild command: %s\nOutput dir: %s\n\n",
		config.Framework, config.BuildCommand, config.OutputDir)
	result.Logs = buildLog

	// Выполняем билд
	cmd = exec.Command("sh", "-c", config.BuildCommand)
	cmd.Dir = clonePath
	cmd.Env = os.Environ()
	
	// Добавляем переменные окружения
	for k, v := range config.Environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	stdout.Reset()
	stderr.Reset()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startTime := time.Now()
	err := cmd.Run()
	result.Duration = int(time.Since(startTime).Seconds())

	result.Logs += stdout.String()
	if stderr.Len() > 0 {
		result.Logs += "\nSTDERR:\n" + stderr.String()
	}

	if err != nil {
		result.Status = "failed"
		result.Logs += fmt.Sprintf("\nBuild failed: %v", err)
	} else {
		result.Status = "success"
		
		// Копируем выходные файлы в preview директорию
		outputPath := filepath.Join(buildPath, "output")
		if err := os.MkdirAll(outputPath, 0755); err == nil {
			sourcePath := filepath.Join(clonePath, config.OutputDir)
			
			// Копируем файлы
			cmd = exec.Command("cp", "-r", sourcePath+"/.", outputPath+"/")
			if err := cmd.Run(); err == nil {
				result.OutputPath = outputPath
				result.PreviewURL = fmt.Sprintf("/preview/%s", buildID)
			}
		}
	}

	result.CompletedAt = time.Now()

	return result, nil
}