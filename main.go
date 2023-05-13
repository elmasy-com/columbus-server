package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/elmasy-com/columbus-sdk/db"
	"github.com/elmasy-com/columbus-server/config"
	"github.com/elmasy-com/columbus-server/server"
)

var (
	Version string
	Commit  string
)

// Returns "up-to-date" and code 0, if no update required.
// Returns the latest release version (eg.: "v0.9.1") and code 1, if update available.
// Returns the error string and code 2, if error happened.
func checkUpdate() {

	resp, err := http.Get("https://api.github.com/repos/elmasy-com/columbus-server/releases/latest")
	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP failed: %s\n", err)
		os.Exit(2)
	}
	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read body: %s\n", err)
		os.Exit(2)
	}

	var v struct {
		TagName string `json:"tag_name"`
	}

	err = json.Unmarshal(out, &v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal: %s\n", err)
		os.Exit(2)
	}
	if v.TagName == "" {
		fmt.Fprint(os.Stderr, "Failed to unmarshal: TagName is empty\n")
		os.Exit(2)
	}

	if Version < v.TagName {
		fmt.Printf("%s", v.TagName)
		os.Exit(1)
	}

	fmt.Printf("up-to-date")
	os.Exit(0)
}

func main() {

	path := flag.String("config", "", "Path to the config file.")
	version := flag.Bool("version", false, "Print version informations.")
	check := flag.Bool("check", false, "Check for updates.")
	flag.Parse()

	if *version {
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Git Commit: %s\n", Commit)
		os.Exit(0)
	}

	if *check {
		// checkUpdate will exit inside the function
		checkUpdate()
	}

	if *path == "" {
		fmt.Fprintf(os.Stderr, "Path to the config file is missing!\n")
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Parsing config file...\n")
	if err := config.Parse(*path); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config file: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Connecting to MongoDB...\n")
	if err := db.Connect(config.MongoURI); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to MongoDB: %s\n", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	if config.EnableStatAPI {
		fmt.Printf("Starting UpdateStat...\n")
		go server.UpdateStat()
	}

	fmt.Printf("Starting HTTP server...\n")
	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("HTTP server stopped!\n")
	}
}
