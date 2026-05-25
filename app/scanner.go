package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
)

// Scanner handles HTTP probing for a specific exploit.
type Scanner struct {
	client        *http.Client
	exploit       Exploit
	configuration Configuration
	regex         *regexp.Regexp
	writer        *ResultWriter
}

// NewScanner creates a Scanner with a pre-compiled validation regex.
func NewScanner(cfg Configuration, exploit Exploit, writer *ResultWriter) (*Scanner, error) {
	re, err := regexp.Compile(exploit.ValidationCriteria)
	if err != nil {
		return nil, fmt.Errorf("compile validation regex: %w", err)
	}

	return &Scanner{
		client: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		exploit:       exploit,
		configuration: cfg,
		regex:         re,
		writer:        writer,
	}, nil
}

// Scan performs a single HTTP request and validates the response.
func (s *Scanner) Scan(ctx context.Context, target string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	headers := s.configuration.RequestHeaders
	if len(headers) > 0 {
		req.Header.Set("User-Agent", headers[rand.Intn(len(headers))])
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	var data string
	switch s.exploit.Response {
	case "body":
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read body: %w", err)
		}
		data = string(body)
	case "header":
		data = resp.Header.Get("Content-Type")
	default:
		// Fallback to body if response type is unrecognized.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read body: %w", err)
		}
		data = string(body)
	}

	if s.regex.MatchString(data) {
		ColorPrint(Green, "%s %s\n", target, s.exploit.Description)
		if s.writer != nil {
			if err := s.writer.Write(target); err != nil {
				fmt.Fprintf(os.Stderr, "[-] failed to write result: %v\n", err)
			}
		}
	} else {
		ColorPrint(Yellow, "%s %s\n", target, "Not Found")
	}
	return nil
}

// ColorPrint formats and prints colored text to stdout.
func ColorPrint(color string, format string, a ...interface{}) {
	fmt.Printf(color+format+Reset, a...)
}
