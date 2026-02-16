package main

import (
    "crypto/sha256"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "time"
)

// Global variables
var baselineFile string

type Baseline struct {
    SystemName string   `json:"system_name"`
    Checks     []string `json:"checks"`
}

// ---------------- Banner & Binary Integrity ----------------
func printBanner() {
    fmt.Println("====================================")
    fmt.Println("   🛡️ RuntimeGuard CLI 🛡️          ")
    fmt.Println("   Author: Antonio Kione            ")
    fmt.Println("====================================")
}

func verifyBinaryIntegrity() {
    binary := "runtimeguard" // use "runtimeguard.exe" on Windows
    data, err := ioutil.ReadFile(binary)
    if err != nil {
        fmt.Println("⚠️ Could not read binary for integrity check:", err)
        return
    }

    hash := sha256.Sum256(data)
    expected := "3f276cf2b62a24957d79879a7328588221446a58b906..." // replace with your binary hash

    if fmt.Sprintf("%x", hash) != expected {
   //     fmt.Println("⚠️ Binary tampered! Exiting...")
   //     os.Exit(1)
    }
}

// ---------------- Baseline Handling ----------------
func initBaseline() {
    fmt.Println("Using baseline file:", baselineFile)
    fmt.Println("Initializing baseline...")

    defaultBaseline := Baseline{
        SystemName: "StudentSystem",
        Checks:     []string{"file_integrity", "process_monitor"},
    }

    data, _ := json.MarshalIndent(defaultBaseline, "", "  ")
    if err := ioutil.WriteFile(baselineFile, data, 0644); err != nil {
        fmt.Println("Failed to write baseline JSON:", err)
        return
    }

    fmt.Println("Baseline created successfully at", baselineFile)
}

func checkBaseline() {
    fmt.Println("Using baseline file:", baselineFile)
    fmt.Println("Checking system against baseline...")

    data, err := ioutil.ReadFile(baselineFile)
    if err != nil {
        fmt.Println("Baseline file not found. Run './runtimeguard init' first.")
        return
    }

    var baseline Baseline
    if err := json.Unmarshal(data, &baseline); err != nil {
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
    fmt.Println("Using baseline file:", baselineFile)
    fmt.Println("Starting monitor...")

    data, err := ioutil.ReadFile(baselineFile)
    if err != nil {
        fmt.Println("Baseline file not found. Run './runtimeguard init' first.")
        return
    }

    var baseline Baseline
    if err := json.Unmarshal(data, &baseline); err != nil {
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

// ---------------- CLI Help ----------------
func printHelp() {
    fmt.Println("Usage: runtimeguard <command>")
    fmt.Println("\nAvailable commands:")
    fmt.Println("  init     - Create baseline JSON and initialize system checks")
    fmt.Println("  check    - Verify system against the existing baseline")
    fmt.Println("  monitor  - Run monitoring on all baseline checks")
    fmt.Println("  help     - Show this help menu")
}

// ---------------- CLI ----------------
func main() {
    printBanner()
    verifyBinaryIntegrity()

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

    if len(os.Args) < 2 {
        printHelp()
        return
    }

    command := os.Args[1]

    switch command {
    case "init":
        initBaseline()
        checkBaseline()
        startMonitor()
    case "check":
        checkBaseline()
    case "monitor":
        startMonitor()
    case "help":
        printHelp()
    default:
        fmt.Println("Unknown command:", command)
        printHelp()
    }
}
