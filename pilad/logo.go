package main

import "log"

func logo(conn *Conn) {
	log.Println()
	log.Println("  .___.            _  _             _  _     ")
	log.Println(" /  _  \\    _ __  (_)| |  __ _   __| || |__  ")
	log.Println("|  |+|  |  | '_ \\ | || | / _` | / _` || '_ \\ ")
	log.Println("|  |-|  |  | |_) || || || (_| || (_| || |_) |")
	log.Println(" \\.___./   | .__/ |_||_| \\__,_| \\__,_||_.__/ ")
	log.Println("           |_|                               ")
	log.Println()
	log.Printf("Version:      %s", conn.Status.Version)
	log.Printf("Go Version:   %s", conn.Status.GoVersion)
	log.Printf("Host:         %s", conn.Status.Host)
	log.Printf("Port:         %d", conn.Config.Port())
	log.Printf("PID:          %d", conn.Status.PID)
	log.Printf("Started at:   %s", conn.Status.StartedAt)
	log.Println()

	if !conn.Config.NoDonate() {
		log.Println("If you want to support open source development of piladb")
		log.Println("please consider making a donation: https://www.paypal.me/oscillatingworks")
		log.Println("Thanks!")
		log.Println()
	}
}
