package models

type APIResponse struct {
	OK        bool        `json:"ok"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
}