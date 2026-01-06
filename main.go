package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func resetTerminal() {
	// Disable mouse tracking modes
	fmt.Print("\x1b[?1000l") // Disable mouse click tracking
	fmt.Print("\x1b[?1002l") // Disable mouse drag tracking
	fmt.Print("\x1b[?1003l") // Disable all mouse tracking
	fmt.Print("\x1b[?1006l") // Disable SGR mouse mode
	fmt.Print("\x1b[?1015l") // Disable urxvt mouse mode
	fmt.Print("\x1b[?25h")   // Show cursor
	fmt.Print("\x1b[0m")     // Reset text attributes
}

func findOpencode() string {
	// 1. Check OPENCODE_PATH env var
	if path := os.Getenv("OPENCODE_PATH"); path != "" {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// 2. Check opencode-core.exe next to wrapper (for bundled install)
	if exePath, err := os.Executable(); err == nil {
		corePath := filepath.Join(filepath.Dir(exePath), "opencode-core.exe")
		if _, err := os.Stat(corePath); err == nil {
			return corePath
		}
	}

	// 3. Check common install locations
	home, _ := os.UserHomeDir()
	locations := []string{
		filepath.Join(home, ".bun", "bin", "opencode.exe"),
		filepath.Join(home, ".local", "bin", "opencode.exe"),
		filepath.Join(home, "AppData", "Local", "Programs", "opencode", "opencode.exe"),
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return loc
		}
	}

	// 4. Try PATH
	if path, err := exec.LookPath("opencode.exe"); err == nil {
		return path
	}

	return ""
}

func main() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleCtrlHandler := kernel32.NewProc("SetConsoleCtrlHandler")

	// Block Ctrl+C in wrapper process - child won't receive it either
	procSetConsoleCtrlHandler.Call(0, 1)
	defer procSetConsoleCtrlHandler.Call(0, 0)

	// Always reset terminal on exit
	defer resetTerminal()

	opencodePath := findOpencode()
	if opencodePath == "" {
		fmt.Fprintln(os.Stderr, "Error: opencode.exe not found")
		fmt.Fprintln(os.Stderr, "Set OPENCODE_PATH env var or install opencode first")
		os.Exit(1)
	}

	cmd := exec.Command(opencodePath, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Create child in new process group - isolate from console Ctrl+C
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
