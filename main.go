package main

import (
	"bufio"
	"fmt"
	"go-web-scanner/app"
	"os"

	// "runtime"
	"strconv"
	"strings"
	"sync"
)

func checkResultsFolder() {
	_, err := os.Stat("results")
	if os.IsNotExist(err) {
		os.Mkdir("results", 0755)
	}
}

func readUserChoice(inputReader *bufio.Reader) int {
	fmt.Print("\nChoice ? ")
	userChoice, _ := inputReader.ReadString('\n')
	userChoiceInt, err := strconv.Atoi(strings.TrimSpace(userChoice))
	if err != nil {
		panic(err)
	}
	return userChoiceInt
}

func readListFile(inputReader *bufio.Reader) []string {
	fmt.Print("List ? ")
	listFileName, _ := inputReader.ReadString('\n')
	listFileName = strings.TrimSpace(listFileName)
	fileContent, err := os.ReadFile(listFileName)
	if err != nil {
		panic(err)
	}
	urlList := strings.Split(string(fileContent), "\n")
	return urlList
}

func readThreads(inputReader *bufio.Reader) int {
	fmt.Print("Threads ? ")
	threadInputStr, _ := inputReader.ReadString('\n')
	threadInput, err := strconv.Atoi(strings.TrimSpace(threadInputStr))
	if err != nil {
		panic(err)
	}
	return threadInput
}

func main() {
	checkResultsFolder()

	app.DisplayBanner()
	inputReader := bufio.NewReader(os.Stdin)
	Config, err := app.LoadConfiguration()
	if err != nil {
		panic(err)
	}

	for index, config := range Config.Exploits {
		app.ColorPrint(app.Green, "%d. Exploit Description: %s\n", index+1, config.Description)
	}

	userChoiceInt := readUserChoice(inputReader)
	userChoiceConfig := Config.Exploits[userChoiceInt-1]

	urlList := readListFile(inputReader)

	threadInput := readThreads(inputReader)

	paths, err := app.LoadFile(userChoiceConfig)
	if err != nil {
		panic(err)
	}

	threadChannel := make(chan struct{}, threadInput)

	var tasks = make([]string, 0, len(urlList)*len(paths))
	for _, url := range urlList {
		for path := range paths {
			// tasks = append(tasks, app.FilterUrl(url)+paths[path])
			go func(url string, path string) {
				tasks = append(tasks, app.FilterUrl(url)+path)
			}(url, paths[path])
		}
	}

	fmt.Printf("\nScanning %d tasks with %d threads\n\n", len(tasks), threadInput)

	scanner := app.NewScanner(Config.Configuration, userChoiceConfig)

	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(task string) {
			threadChannel <- struct{}{}
			scanner.Scan(task)
			<-threadChannel
			wg.Done()
		}(task)
	}
	wg.Wait()
	// runtime.GOMAXPROCS(threadInput)
	// var wg sync.WaitGroup
	// for _, task := range tasks {
	// 	wg.Add(1)
	// 	go func(task string) {
	// 		scanner.Scan(task)
	// 		wg.Done()
	// 	}(task)
	// }
	// wg.Wait()
}
