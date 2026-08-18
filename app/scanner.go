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
	"strings"
)

const maxRedirects = 3

// Scanner handles HTTP probing for a specific exploit.
type Scanner struct {
	client        *http.Client
	exploit       Exploit
	configuration Configuration
	regex         *regexp.Regexp
	writer        *ResultWriter
}

// NewScanner creates a Scanner with a pre-compiled validation regex.
func NewScanner(cfg Configuration, exploit Exploit, writer *ResultWriter, insecure bool) (*Scanner, error) {
	re, err := regexp.Compile(exploit.ValidationCriteria)
	if err != nil {
		return nil, fmt.Errorf("compile validation regex: %w", err)
	}

	return &Scanner{
		client: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= maxRedirects {
					return http.ErrUseLastResponse
				}
				if len(via) > 0 && req.URL.Host != via[0].URL.Host {
					return http.ErrUseLastResponse
				}
				return nil
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

	for k, v := range s.configuration.Headers {
		req.Header.Set(k, v)
	}

	if len(s.configuration.UserAgents) > 0 {
		req.Header.Set("User-Agent", s.configuration.UserAgents[rand.Intn(len(s.configuration.UserAgents))])
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	var data string
	switch s.exploit.Response {
	case "body":
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read body: %w", err)
		}
		data = string(body)
	case "header":
		data = formatHeaders(resp.Header)
		_, _ = io.Copy(io.Discard, resp.Body)
	default:
		// Fallback to body if response type is unrecognized.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read body: %w", err)
		}
		data = string(body)
	}

	if s.regex.MatchString(data) {
		ColorPrint(Green, "%s [%d] %s\n", target, resp.StatusCode, s.exploit.Description)
		if s.writer != nil {
			if err := s.writer.Write(target); err != nil {
				fmt.Fprintf(os.Stderr, "[-] failed to write result: %v\n", err)
			}
		}
	} else {
		ColorPrint(Yellow, "%s [%d] %s\n", target, resp.StatusCode, "Not Found")
	}
	return nil
}

func formatHeaders(h http.Header) string {
	var b strings.Builder
	for k, v := range h {
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(strings.Join(v, ", "))
		b.WriteString("\n")
	}
	return b.String()
}

// ColorPrint formats and prints colored text to stdout.
func ColorPrint(color string, format string, a ...interface{}) {
	fmt.Printf(color+format+Reset, a...)
}
