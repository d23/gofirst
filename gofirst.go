package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

// Waits for the command to finish execution and reaps any child other
// child process that comes along.
func waitAndReap(cmd *exec.Cmd) {
	var status syscall.WaitStatus
	for {
		// Wait for any process, not just the command.
		pid, err := syscall.Wait4(-1, &status, 0, nil)
		// If no more children exist this should return an error.
		if err != nil {
			log.Fatal(err)
			break
		}
		if cmd.Process.Pid == pid {
			log.Print("Command completed with status ", status)
		} else {
			log.Print("Reaped child with pid ", pid, " status : ", status)
		}
	}
}

// Installs the signal handler, trapping SiGINT, SIGTERM and SIGHUP
// We're not bothering about SIGKILL as we can catch it anyways.
func installSignalHandler(cmd *exec.Cmd, timeout int) chan os.Signal {
	c := make(chan os.Signal, 2)
	signal.Notify(c,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1)
	go signalHandler(cmd, c, timeout)
	return c
}

// Waits for signals to arrive and transmits them to the running command
// In case of SIGTERM it will broadcast is to any process it can, and
// finally send SIGKILL if we have not yet exited before then.
func signalHandler(cmd *exec.Cmd, c chan os.Signal, timeout int) {
	for sig := range c {
		switch sig {
		case syscall.SIGINT, syscall.SIGHUP:
			// Just pass it along
			log.Print("Received ", sig)
			cmd.Process.Signal(sig)
		case syscall.SIGTERM:
			// Send signal for brutal kill after timeout
			go func() {
				time.Sleep(time.Duration(timeout) * time.Second)
				syscall.Kill(os.Getpid(), syscall.SIGUSR1)
			}()
			// Send SIGTERM to all children for graceful shutdown
			syscall.Kill(-1, syscall.SIGTERM)
		case syscall.SIGUSR1:
			// Kill everything and exit the handler
			syscall.Kill(-1, syscall.SIGKILL)
			return
		}
	}
}

// Run the process, route it's input and output to standard channels
func startProcess() *exec.Cmd {
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

    return cmd
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("No command supplied.")
	}

	// Start the process
	cmd := startProcess()

	// Trap and handle signals
	c := installSignalHandler(cmd, 30)

	// Make sure to terminate the handler when we're done
	defer close(c)

	// Wait until we're done while reaping and process that comes along
	waitAndReap(cmd)

}
