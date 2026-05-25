package app

import "fmt"

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Reset  = "\033[0m"
)

func DisplayBanner() {
	fmt.Println(Cyan + `
   █▀▀ █▀█ ▄▄ █░█░█ █▀▀ █▄▄
   █▄█ █▄█ ░░ ▀▄▀▄▀ ██▄ █▄█` + Green + `

   █▀ █▀▀ ▄▀█ █▄░█ █▄░█ █▀▀ █▀█
   ▄█ █▄▄ █▀█ █░▀█ █░▀█ ██▄ █▀▄` + Yellow + `  v2.0.0
` + Reset)
}
