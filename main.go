package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"go-web-scanner/app"
	"os"
	"sync"
	"time"
)

type options struct {
	configPath string
	listFile   string
	exploitIdx int
	workers    int
	outputDir  string
}

func parseFlags() options {
	var opts options
	flag.StringVar(&opts.configPath, "config", "./files/config.json", "Path to config.json")
	flag.StringVar(&opts.listFile, "list", "", "File containing target URLs")
	flag.IntVar(&opts.exploitIdx, "exploit", 0, "Exploit index to run (1-based, 0 = list available exploits)")
	flag.IntVar(&opts.workers, "workers", 20, "Number of concurrent workers")
	flag.StringVar(&opts.outputDir, "output", "results", "Output directory for results")
	flag.Parse()
	return opts
}

func loadTargets(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open target file: %w", err)
	}
	defer f.Close()

	var targets []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if u := app.NormalizeURL(scanner.Text()); u != "" {
			targets = append(targets, u)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read target file: %w", err)
	}
	return targets, nil
}

func buildTasks(targets, paths []string) []string {
	tasks := make([]string, 0, len(targets)*len(paths))
	for _, target := range targets {
		for _, path := range paths {
			tasks = append(tasks, target+path)
		}
	}
	return tasks
}

func run(opts options) error {
	app.DisplayBanner()

	cfg, err := app.LoadConfiguration(opts.configPath)
	if err != nil {
		return err
	}

	// List exploits if none selected.
	if opts.exploitIdx == 0 {
		fmt.Println("Available exploits:")
		for i, e := range cfg.Exploits {
			app.ColorPrint(app.Green, "  %d. %s\n", i+1, e.Description)
		}
		fmt.Println("\nUse -exploit <number> to select one.")
		return nil
	}

	if opts.exploitIdx < 1 || opts.exploitIdx > len(cfg.Exploits) {
		return fmt.Errorf("invalid exploit index: %d (available: 1-%d)", opts.exploitIdx, len(cfg.Exploits))
	}

	exploit := cfg.Exploits[opts.exploitIdx-1]
	fmt.Printf("Selected: %s\n\n", exploit.Description)

	if opts.listFile == "" {
		return fmt.Errorf("-list flag is required")
	}

	targets, err := loadTargets(opts.listFile)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return fmt.Errorf("no valid targets found")
	}

	paths, err := app.LoadPaths(exploit.FilePath)
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		return fmt.Errorf("no paths loaded from %s", exploit.FilePath)
	}

	tasks := buildTasks(targets, paths)
	fmt.Printf("Scanning %d tasks with %d workers\n\n", len(tasks), opts.workers)

	writer, err := app.NewResultWriter(opts.outputDir, exploit.SaveAs)
	if err != nil {
		return err
	}
	defer writer.Close()

	scanner, err := app.NewScanner(cfg.Configuration, exploit, writer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(len(tasks))*cfg.Configuration.Timeout+30*time.Second)
	defer cancel()

	jobs := make(chan string, opts.workers)
	var wg sync.WaitGroup

	// Start workers.
	for i := 0; i < opts.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range jobs {
				if err := scanner.Scan(ctx, task); err != nil {
					app.ColorPrint(app.Red, "[-] %s error: %v\n", task, err)
				}
			}
		}()
	}

	// Dispatch tasks.
	for _, task := range tasks {
		jobs <- task
	}
	close(jobs)

	wg.Wait()
	fmt.Printf("\nDone. Results saved to %s/%s\n", opts.outputDir, exploit.SaveAs)
	return nil
}

func main() {
	opts := parseFlags()
	if err := run(opts); err != nil {
		app.ColorPrint(app.Red, "[-] %v\n", err)
		os.Exit(1)
	}
}
