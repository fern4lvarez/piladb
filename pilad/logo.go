package main

import (
	"fmt"
	"log/slog"
)

func logo(conn *Conn) {
	// Print ASCII art to stdout (not structured logs)
	fmt.Println()
	fmt.Println("         d8b 888               888 888      ")
	fmt.Println("         Y8P 888               888 888      ")
	fmt.Println("             888               888 888      ")
	fmt.Println("88888b.  888 888  8888b.   .d88888 88888b.  ")
	fmt.Println("888 \"88b 888 888    \"88b  d88\" 888 888 \"88b ")
	fmt.Println("888  888 888 888 .d888888 888  888 888  888 ")
	fmt.Println("888 d88P 888 888 888  888 Y88b 888 888 d88P ")
	fmt.Println("88888P\"  888 888 \"Y888888  \"Y88888 88888P\"  ")
	fmt.Println("888")
	fmt.Println("888")
	fmt.Println("888")
	fmt.Println()

	slog.Info("starting pilad",
		"version", conn.Status.Version,
		"go_version", conn.Status.GoVersion,
		"host", conn.Status.Host,
		"port", conn.Config.Port(),
		"pid", conn.Status.PID,
	)
}
