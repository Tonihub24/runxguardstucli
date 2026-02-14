package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// global variables
var baselineFile string

type Baseline struct {
	SystemName string   `json:"system_name"`
	Checks     []string `json:"checks"`
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not determine home directory:", err)
		return
	}

	configDir := filepath.Join(homeDir, ".runtimeguard")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0755); err != nil {
			fmt.Println("Failed to create config directory:", err)
			return
		}
	}

	baselineFile = filepath.Join(configDir, "baseline.json")
	fmt.Println("Using baseline file:", baselineFile)

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		initBaseline()
	case "check":
		checkBaseline()
	case "monitor":
		startMonitor()
	default:
		fmt.Println("Unknown command:", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: runtimeguard <init|check|monitor>")
}

// --- Implemented functions ---

func initBaseline() {
	fmt.Println("Initializing baseline...")

	defaultBaseline := Baseline{
		SystemName: "StudentSystem",
		Checks:     []string{"file_integrity", "process_monitor"},
	}

	file, err := os.Create(baselineFile)
	if err != nil {
		fmt.Println("Failed to create baseline file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(defaultBaseline); err != nil {
		fmt.Println("Failed to write baseline JSON:", err)
		return
	}

	fmt.Println("Baseline created successfully at", baselineFile)
}

func checkBaseline() {
	fmt.Println("Checking system against baseline...")

	file, err := os.Open(baselineFile)
	if err != nil {
		fmt.Println("Baseline file not found. Run './runtimeguard init' first.")
		return
	}
	defer file.Close()

	var baseline Baseline
	if err := json.NewDecoder(file).Decode(&baseline); err != nil {
		fmt.Println("Failed to parse baseline JSON:", err)
		return
	}

	fmt.Println("System Name:", baseline.SystemName)
	fmt.Println("Checks in baseline:")
	for _, check := range baseline.Checks {
		fmt.Println(" -", check)
	}
}

func startMonitor() {
	fmt.Println("Starting monitor...")

	file, err := os.Open(baselineFile)
	if err != nil {
		fmt.Println("Baseline file not found. Run './runtimeguard init' first.")
		return
	}
	defer file.Close()

	var baseline Baseline
	if err := json.NewDecoder(file).Decode(&baseline); err != nil {
		fmt.Println("Failed to parse baseline JSON:", err)
		return
	}

	for _, check := range baseline.Checks {
		fmt.Printf("Checking %s ... ", check)
		time.Sleep(1 * time.Second) // simulate work
		fmt.Println("OK")
	}

	fmt.Println("Monitoring complete.")
}
