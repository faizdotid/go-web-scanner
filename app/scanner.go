package app

import (
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	// "strings"
	"time"
)

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Reset  = "\033[0m"
)

type Scanner struct {
	client        *http.Client
	exploit       Exploit
	configuration Configuration
}

func NewScanner(c Configuration, e Exploit) *Scanner {
	return &Scanner{
		configuration: c,
		exploit:       e,
		client: &http.Client{
			Timeout: time.Duration(c.Timeout) * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func ColorPrint(color string, format string, a ...interface{}) {
	fmt.Printf(color+format+Reset, a...)
}

func (s *Scanner) Scan(url string) {
	defer func() {
		if r := recover(); r != nil {
			ColorPrint(Red, "Error: %v\n", r)
		}
	}()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", s.configuration.RequestHeaders[rand.Intn(len(s.configuration.RequestHeaders))])
	resp, err := s.client.Do(req)
	if err != nil {
		panic(err)
	}
	var Data string
	if s.exploit.Response == "body" {
		byteData, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		Data = string(byteData)
	}
	if s.exploit.Response == "header" {
		headerData := resp.Header.Get("Content-Type")
		Data = headerData
	}
	defer resp.Body.Close()

	match, err := regexp.MatchString(s.exploit.ValidationCriteria, Data)
	if err != nil {
		panic(err)
	}
	if match {
		ColorPrint(Green, "%s %s\n", url, s.exploit.Description)
		WriteResult(fmt.Sprintf("results/%s", s.exploit.SaveAs), url)
	} else {
		ColorPrint(Yellow, "%s %s\n", url, "Not Found")
	}
}
