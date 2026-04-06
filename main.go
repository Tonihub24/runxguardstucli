package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// ---------------- Structures ----------------
type FileCheck struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

type Baseline struct {
	SystemName string      `json:"system_name"`
	Files      []FileCheck `json:"files"`
	Processes  []string    `json:"processes"`
}

// ---------------- Global Variables ----------------
var baselineFile string
var logFile string

// ---------------- Helper Functions ----------------
func getTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func logMessage(level string, msg string) {
	timestamp := getTimestamp()
	logLine := fmt.Sprintf("[%s] [%s] %s", timestamp, level, msg)

	// Print to terminal
	fmt.Println(logLine)

	// Write to log file
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		fmt.Fprintln(f, logLine)
	}
}

func calculateFileHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}

func getCurrentShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "bash"
	}
	parts := strings.Split(shell, "/")
	return parts[len(parts)-1]
}

// ---------------- Banner ----------------
// ---------------- Banner ----------------
func printBanner() {
	logMessage("INFO", "====================================")
	logMessage("INFO", "   🛡️ RuntimeGuard CLI 🛡️          ")
	logMessage("INFO", "   Author: Antonio Kione            ")
	logMessage("INFO", "====================================")
}

// ---------------- Baseline Handling ----------------
func initBaseline() {
	logMessage("INFO", "Initializing baseline...")

	criticalFiles := []string{"/etc/passwd", "/etc/group"}
	files := []FileCheck{}

	for _, f := range criticalFiles {
		hash, err := calculateFileHash(f)
		if err != nil {
			logMessage("WARN", fmt.Sprintf("Cannot read %s: %v", f, err))
			continue
		}
		files = append(files, FileCheck{Path: f, Hash: hash})
	}

	currentShell := getCurrentShell()
	allowedProcesses := []string{"systemd", currentShell, "runtimeguard"}

	baseline := Baseline{
		SystemName: "StudentSystem",
		Files:      files,
		Processes:  allowedProcesses,
	}

	data, _ := json.MarshalIndent(baseline, "", "  ")
	if err := os.WriteFile(baselineFile, data, 0644); err != nil {
		logMessage("ERROR", fmt.Sprintf("Failed to write baseline: %v", err))
		return
	}

	logMessage("INFO", fmt.Sprintf("Baseline created at %s", baselineFile))
}

func checkBaseline() {
	data, err := os.ReadFile(baselineFile)
	if err != nil {
		logMessage("ERROR", "Baseline not found. Run init first.")
		return
	}

	var baseline Baseline
	if err := json.Unmarshal(data, &baseline); err != nil {
		logMessage("ERROR", fmt.Sprintf("Failed to parse baseline: %v", err))
		return
	}

	logMessage("INFO", "System Name: "+baseline.SystemName)
	logMessage("INFO", "Files in baseline:")
	for _, f := range baseline.Files {
		logMessage("INFO", " - "+f.Path)
	}
	logMessage("INFO", "Allowed processes:")
	for _, p := range baseline.Processes {
		logMessage("INFO", " - "+p)
	}
}

// ---------------- Continuous Monitoring ----------------
func startMonitor() {
	data, err := os.ReadFile(baselineFile)
	if err != nil {
		logMessage("ERROR", "Baseline not found. Run init first.")
		select {} // block forever so systemd doesn't think it exited
	}

	var baseline Baseline
	if err := json.Unmarshal(data, &baseline); err != nil {
		logMessage("ERROR", fmt.Sprintf("Failed to parse baseline: %v", err))
		select {} // block forever
	}

	logMessage("INFO", "Starting continuous monitoring...")

	for {
		func() {
			// --- File Integrity ---
			logMessage("INFO", "---- File Integrity Check ----")
			for _, f := range baseline.Files {
				hash, err := calculateFileHash(f.Path)
				if err != nil {
					logMessage("WARN", fmt.Sprintf("%s cannot be read", f.Path))
					continue
				}
				if hash != f.Hash {
					logMessage("WARN", fmt.Sprintf("%s tampered!", f.Path))
				} else {
					logMessage("INFO", fmt.Sprintf("%s OK", f.Path))
				}
			}

			// --- Process Check ---
			logMessage("INFO", "---- Process Check ----")
			out, _ := exec.Command("ps", "-eo", "comm").Output()
			running := strings.Split(string(out), "\n")
			for _, p := range baseline.Processes {
				found := false
				for _, r := range running {
					if strings.Contains(r, p) {
						found = true
						break
					}
				}
				if !found {
					logMessage("WARN", fmt.Sprintf("Missing process: %s", p))
				} else {
					logMessage("INFO", fmt.Sprintf("Process running: %s", p))
				}
			}

			// --- Port Check ---
			logMessage("INFO", "---- Port Check ----")
			conns, _ := net.Connections("inet")
			for _, c := range conns {
				if c.Status == "LISTEN" {
					logMessage("INFO", fmt.Sprintf("Port listening: %d", c.Laddr.Port))
				}
			}
		}()

		time.Sleep(5 * time.Second)
	}
}

// ---------------- Extra Functions ----------------
func listProcesses() {
	procs, _ := process.Processes()
	logMessage("INFO", "PID    PROCESS")
	for _, p := range procs {
		name, _ := p.Name()
		logMessage("INFO", fmt.Sprintf("%d    %s", p.Pid, name))
	}
}

func watchDirectory(path string) {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	eventTimes := make(map[string]time.Time)
	debounce := 1 * time.Second

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// Skip temporary files
				if strings.HasPrefix(filepath.Base(event.Name), "#") {
					continue
				}

				last, exists := eventTimes[event.Name]
				if !exists || time.Since(last) > debounce {
					logMessage("INFO", fmt.Sprintf("File event: %s %s", event.Op, event.Name))
					eventTimes[event.Name] = time.Now()
				}

			case err := <-watcher.Errors:
				logMessage("ERROR", fmt.Sprintf("Watcher error: %v", err))
			}
		}
	}()

	watcher.Add(path)
	logMessage("INFO", "Watching directory: "+path)
	select {}
}

// ---------------- CLI ----------------
func printHelp() {
	logMessage("INFO", "Usage: runtimeguard <command>")
	logMessage("INFO", "  init        - Create baseline JSON and initialize system")
	logMessage("INFO", "  check       - Verify system against the existing baseline")
	logMessage("INFO", "  monitor     - Run continuous monitoring on baseline files")
	logMessage("INFO", "  processes   - List running processes")
	logMessage("INFO", "  ports       - Show listening ports")
	logMessage("INFO", "  watch <dir> - Watch a directory for changes")
	logMessage("INFO", "  help        - Show this help menu")
}

// ---------------- Main ----------------
func main() {
	printBanner()

	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".runtimeguard")
	os.MkdirAll(configDir, 0755)
	baselineFile = filepath.Join(configDir, "baseline.json")
	logFile = filepath.Join(configDir, "runtimeguard.log")

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]
	switch command {
	case "init":
		initBaseline()
		checkBaseline()
	case "check":
		checkBaseline()
	case "monitor":
		startMonitor()
	case "processes":
		listProcesses()
	case "ports":
		conns, _ := net.Connections("inet")
		logMessage("INFO", "LISTENING PORTS")
		for _, c := range conns {
			if c.Status == "LISTEN" {
				logMessage("INFO", fmt.Sprintf("Port: %d", c.Laddr.Port))
			}
		}
	case "watch":
		if len(os.Args) < 3 {
			fmt.Println("Usage: runtimeguard watch <directory>")
			return
		}
		watchDirectory(os.Args[2])
	case "help":
		printHelp()
	default:
		fmt.Println("Unknown command:", command)
		printHelp()
	}
}
