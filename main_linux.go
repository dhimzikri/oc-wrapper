//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
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

	// 2. Check opencode-core next to wrapper (for bundled install)
	if exePath, err := os.Executable(); err == nil {
		corePath := filepath.Join(filepath.Dir(exePath), "opencode-core")
		if _, err := os.Stat(corePath); err == nil {
			return corePath
		}
	}

	// 3. Check common install locations
	home, _ := os.UserHomeDir()
	locations := []string{
		filepath.Join(home, ".bun", "bin", "opencode"),
		filepath.Join(home, ".local", "bin", "opencode"),
		"/usr/local/bin/opencode",
		"/usr/bin/opencode",
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return loc
		}
	}

	// 4. Try PATH
	if path, err := exec.LookPath("opencode"); err == nil {
		return path
	}

	return ""
}

func main() {
	// Ignore SIGINT in wrapper - child handles its own signals
	signal.Ignore(syscall.SIGINT)

	// Always reset terminal on exit
	defer resetTerminal()

	opencodePath := findOpencode()
	if opencodePath == "" {
		fmt.Fprintln(os.Stderr, "Error: opencode not found")
		fmt.Fprintln(os.Stderr, "Set OPENCODE_PATH env var or install opencode first")
		os.Exit(1)
	}

	cmd := exec.Command(opencodePath, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Create child in new process group - isolate from terminal signals
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
