package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gowtham2003/gotable/pkg/tui"
)

func main() {
	// Add flags
	cliMode := flag.Bool("cli", false, "Run in CLI mode")
	flag.Parse()

	if *cliMode {
		// Run in CLI mode
		if len(os.Args) < 3 {
			fmt.Println("Usage: program -cli <input_file> <output_file>")
			os.Exit(1)
		}
		// ... existing CLI mode code ...
		return
	}

	// Run TUI mode by default
	if err := tui.StartTUI(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
