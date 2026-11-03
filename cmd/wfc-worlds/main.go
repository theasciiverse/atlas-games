// WFC Worlds: the Atlas land and city generators as a demo that fills the
// window and keeps going, land then city, until a key is pressed.
package main

import (
	"fmt"
	"os"

	"atlas.games/internal/worlds"
	tea "github.com/charmbracelet/bubbletea"
)

var Version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version":
			fmt.Printf("wfc-worlds v%s\n", Version)
			return
		case "-h", "--help", "help":
			fmt.Println("WFC Worlds: landscapes and cities growing by Wave Function Collapse.")
			fmt.Println("\nUsage:")
			fmt.Println("  wfc-worlds      Run (space: next world, q or Esc: quit)")
			fmt.Println("  wfc-worlds -v   Show version")
			return
		}
	}

	p := tea.NewProgram(worlds.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running demo: %v\n", err)
		os.Exit(1)
	}
}
