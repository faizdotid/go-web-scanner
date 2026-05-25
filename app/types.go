package app

import "time"

// Exploit defines a scanning module loaded from config.json.
type Exploit struct {
	Type               string `json:"type"`
	Description        string `json:"description"`
	FilePath           string `json:"file_path"`
	ValidationCriteria string `json:"validation_criteria"`
	SaveAs             string `json:"save_as"`
	Response           string `json:"response"` // "body" or "header"
}

// Configuration holds global scanner settings.
type Configuration struct {
	Timeout        time.Duration `json:"-"`
	TimeoutSeconds int           `json:"timeout"`
	RequestHeaders []string      `json:"request_headers"`
}

// Config is the root configuration structure.
type Config struct {
	Exploits      []Exploit     `json:"exploits"`
	Configuration Configuration `json:"configuration"`
}
