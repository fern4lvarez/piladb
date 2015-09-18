package main

import "log"

func logo(conn *Conn) {
	log.Println()
	log.Println("         d8b 888               888 888      ")
	log.Println("         Y8P 888               888 888      ")
	log.Println("             888               888 888      ")
	log.Println("88888b.  888 888  8888b.   .d88888 88888b.  ")
	log.Println("888 \"88b 888 888    \"88b  d88\" 888 888 \"88b ")
	log.Println("888  888 888 888 .d888888 888  888 888  888 ")
	log.Println("888 d88P 888 888 888  888 Y88b 888 888 d88P ")
	log.Println("88888P\"  888 888 \"Y888888  \"Y88888 88888P\"  ")
	log.Println("888")
	log.Println("888")
	log.Println("888")
	log.Println()
	log.Printf("Version: %s", conn.Status.Version)
	log.Println()
}
