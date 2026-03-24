package handlers

import "strings"

const (
	ProjectTypeBackend  = "backend"
	ProjectTypeTelegram = "telegram"
)

func normalizeProjectType(value string) string {
	v := strings.TrimSpace(strings.ToLower(value))
	switch v {
	case "", "service", "backend":
		return ProjectTypeBackend
	case "telegram", "telegram_bot":
		return ProjectTypeTelegram
	default:
		return v
	}
}

func isAllowedProjectType(value string) bool {
	return value == ProjectTypeBackend || value == ProjectTypeTelegram
}

func isFrontendDeployConfig(buildCmd, startCmd, outputDir string) bool {
	b := strings.ToLower(strings.TrimSpace(buildCmd))
	s := strings.ToLower(strings.TrimSpace(startCmd))
	o := strings.ToLower(strings.TrimSpace(outputDir))
	if strings.Contains(b, "vite") || strings.Contains(b, "next export") || strings.Contains(b, "react-scripts build") {
		return true
	}
	if (o == "dist" || o == "build" || o == "out") && (s == "" || strings.Contains(s, "serve") || strings.Contains(s, "nginx")) {
		return true
	}
	return false
}
