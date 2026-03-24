package deployer

import (
	"encoding/json"
	"io/ioutil"
)

// DetectFramework reads files in the directory and returns the framework type
func (d *Deployer) DetectFramework(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var filenames []string
	for _, f := range files {
		filenames = append(filenames, f.Name())
	}

	for _, file := range filenames {
		switch {
		case file == "package.json":
			return d.detectNodeFramework(dir + "/package.json"), nil
		case file == "go.mod":
			return "go", nil
		case file == "requirements.txt":
			return "python", nil
		case file == "Gemfile":
			return "ruby", nil
		case file == "Cargo.toml":
			return "rust", nil
		case file == "index.html":
			return "static", nil
		}
	}
	return "static", nil
}

func (d *Deployer) detectNodeFramework(packageJSONPath string) string {
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

	allDeps := make(map[string]string)
	for k, v := range pkg.Dependencies {
		allDeps[k] = v
	}
	for k, v := range pkg.DevDependencies {
		allDeps[k] = v
	}

	if _, ok := allDeps["next"]; ok {
		return "nextjs"
	}
	if _, ok := allDeps["gatsby"]; ok {
		return "gatsby"
	}
	if _, ok := allDeps["vue"]; ok {
		return "vue"
	}
	if _, ok := allDeps["@angular/core"]; ok {
		return "angular"
	}
	if _, ok := allDeps["react"]; ok {
		return "react"
	}
	if _, ok := allDeps["express"]; ok {
		return "express"
	}
	return "node"
}

func (d *Deployer) GetBuildCommand(framework string) string {
	commands := map[string]string{
		"react":    "npm install && npm run build",
		"vue":      "npm install && npm run build",
		"angular":  "npm install && ng build --prod",
		"nextjs":   "npm install && npm run build",
		"gatsby":   "npm install && gatsby build",
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

func (d *Deployer) GetOutputDir(framework string) string {
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
