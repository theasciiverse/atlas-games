// Wilson's Revenge on its own, without the Atlas Games menu.
package main

import (
	"fmt"
	"os"

	"atlas.games/internal/wilson"
	tea "github.com/charmbracelet/bubbletea"
)

var Version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version":
			fmt.Printf("wilsons-revenge v%s\n", Version)
			return
		case "-h", "--help", "help":
			fmt.Println("Wilson's Revenge: a three-lane terminal runner.")
			fmt.Println("\nUsage:")
			fmt.Println("  wilsons-revenge      Play")
			fmt.Println("  wilsons-revenge -v   Show version")
			fmt.Println("  wilsons-revenge -h   Show this help")
			return
		}
	}

	p := tea.NewProgram(wilson.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running game: %v\n", err)
		os.Exit(1)
	}
}
