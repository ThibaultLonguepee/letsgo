package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	pollDelay int = 500 // ms
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatalf("❌ Program was not specified.\nUsage: <files> | watcher <program> [args...]")
	}

	exe := os.Args[1]
	args := os.Args[2:]
	files, err := getFileList()

	var mods map[string]time.Time
	var cmd *exec.Cmd
	var done chan struct{}

	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		log.Fatalf("💥 Failed to open tty file.\n%v", err)
	}
	defer tty.Close()

	state, err := getTtyState(tty)
	handle(err)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigch
		_ = kill(cmd, done)
		_ = setTtyState(tty, state)
		os.Exit(0)
	}()

	for {
		changed := false

		nmods, err := getFileMods(files)
		handle(err)
		for f, t := range mods {
			if nmods[f].After(t) {
				changed = true
			}
		}

		if changed || len(mods) == 0 {
			mods = nmods
			err = kill(cmd, done)
			handle(err)
			err = setTtyState(tty, state)
			handle(err)
			cmd, done, err = start(exe, args...)
			handle(err)
		}

		time.Sleep(time.Duration(pollDelay) * time.Millisecond)
	}
}

func handle(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func getFileList() ([]string, error) {
	files := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			files = append(files, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("💥 Failed to read file list.\n%v", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("❌ No files provided.\n")
	}
	return files, nil
}

func getFileMods(files []string) (map[string]time.Time, error) {
	mods := make(map[string]time.Time, len(files))

	for _, f := range files {
		stat, err := os.Stat(f)
		if err != nil {
			return nil, fmt.Errorf("💥 Failed to access file '%v'.\n%v", f, err)
		}
		mods[f] = stat.ModTime()
	}
	return mods, nil
}

func kill(cmd *exec.Cmd, done chan struct{}) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}

	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("💥 Failed to interrupt the process\n%v", err)
	}

	select {
	case <-done:
		// exited cleanly
	case <-time.After(2 * time.Second):
		_ = cmd.Process.Kill()
		<-done
	}

	return nil
}

func start(exe string, args ...string) (*exec.Cmd, chan struct{}, error) {
	cmd := exec.Command(exe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("💥 Failed to start the process.\n%v", err)
	}

	done := make(chan struct{})

	go func() {
		_ = cmd.Wait()
		close(done)
	}()

	return cmd, done, nil
}

func getTtyState(tty *os.File) (string, error) {
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = tty
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("💥 Failed to save tty settings.\n%v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func setTtyState(tty *os.File, state string) error {
	if state == "" {
		return nil
	}
	cmd := exec.Command("stty", state)
	cmd.Stdin = tty
	return cmd.Run()
}
