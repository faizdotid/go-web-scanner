package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// LoadConfiguration reads and parses the JSON configuration file.
func LoadConfiguration(path string) (Config, error) {
	var cfg Config

	body, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}

	if err := json.Unmarshal(body, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config: %w", err)
	}

	cfg.Configuration.Timeout = time.Duration(cfg.Configuration.TimeoutSeconds) * time.Second
	return cfg, nil
}

// LoadPaths reads a wordlist file and returns non-empty lines.
func LoadPaths(path string) ([]string, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read wordlist: %w", err)
	}

	lines := strings.Split(strings.ReplaceAll(string(body), "\r\n", "\n"), "\n")
	paths := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			paths = append(paths, line)
		}
	}
	return paths, nil
}

// NormalizeURL ensures the URL has a scheme and trims trailing slashes.
func NormalizeURL(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimSuffix(raw, "/")
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	return "http://" + raw
}

// ResultWriter handles buffered, thread-safe writing to result files.
type ResultWriter struct {
	mu     sync.Mutex
	writer *bufio.Writer
	file   *os.File
}

// NewResultWriter creates a ResultWriter for the given output path.
func NewResultWriter(outputDir, filename string) (*ResultWriter, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("create output dir: %w", err)
	}

	path := filepath.Join(outputDir, filename)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open result file: %w", err)
	}

	return &ResultWriter{
		file:   f,
		writer: bufio.NewWriter(f),
	}, nil
}

// Write appends a line to the result file.
func (rw *ResultWriter) Write(line string) error {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	if _, err := fmt.Fprintln(rw.writer, line); err != nil {
		return err
	}
	return nil
}

// Flush ensures all buffered data is written to disk.
func (rw *ResultWriter) Flush() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	return rw.writer.Flush()
}

// Close flushes and closes the underlying file.
func (rw *ResultWriter) Close() error {
	if err := rw.Flush(); err != nil {
		_ = rw.file.Close()
		return err
	}
	return rw.file.Close()
}
