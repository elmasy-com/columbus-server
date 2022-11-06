package main

import (
	"fmt"
	"os"

	"github.com/elmasy-com/columbus-server/blacklist"
	"github.com/elmasy-com/columbus-server/config"
	"github.com/elmasy-com/columbus-server/db"
	"github.com/elmasy-com/columbus-server/server"
)

var (
	Version string
	Commit  string
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Path to config file is missing!\n")
		fmt.Printf("Usage: %s <path-to-config>\n", os.Args[0])
		os.Exit(1)
	}

	if os.Args[1] == "version" {
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Git Commit: %s\n", Commit)
		os.Exit(0)
	}

	fmt.Printf("Parsing config file...\n")
	if err := config.Parse(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config file: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Connecting to MongoDB...\n")
	if err := db.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to MongoDB: %s\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	fmt.Printf("Initializing database...\n")
	if err := db.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updating database...\n")
	if err := db.Update(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update database: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Initializing blacklist...\n")
	blacklist.Init()

	fmt.Printf("Starting HTTP server...\n")
	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("HTTP server stopped!\n")
	}
}
