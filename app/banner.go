package app

import "fmt"

func DisplayBanner() {
	fmt.Print(Red + `
	█▀▀ █▀█ ▄▄ █░█░█ █▀▀ █▄▄` + Reset + Green + `
	█▄█ █▄█ ░░ ▀▄▀▄▀ ██▄ █▄█` + Reset + Yellow + `

	█▀ █▀▀ ▄▀█ █▄░█ █▄░█ █▀▀ █▀█` + Reset + Green + `
	▄█ █▄▄ █▀█ █░▀█ █░▀█ ██▄ █▀▄ V 1.0.0

`)
}
