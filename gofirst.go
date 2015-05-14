package main

import (
    "log"
    "os"
    "os/exec"
    "os/signal"
    "syscall"
    "time"
    "flag"
    "fmt"
)

func debugMsg(msg ...interface{}) {
    if *debug == true {
        log.Print(msg...)
    }
}

var debug = flag.Bool("debug", false, "Show debug messages")

// Waits for the command to finish execution and reaps any child other
// child process that comes along.
func waitAndReap(cmd *exec.Cmd, timeout int) {
    var status syscall.WaitStatus
    for {
        // Wait for any process, not just the command.
        pid, err := syscall.Wait4(-1, &status, 0, nil)
        if err != nil {
            // Stop waiting and return if there are no more children
            if err == syscall.ECHILD {
                break
            }
            log.Fatal("Error: ", err, "(",int(err.(syscall.Errno)), ")")
        }
        if cmd.Process.Pid == pid {
            debugMsg("Command completed with status: ", status)
            // Terminate any other children, as the main command has exited.
            syscall.Kill(-1, syscall.SIGTERM)
            // And add a kill safeguard.
            go func() {
                time.Sleep(time.Duration(timeout) * time.Second)
                debugMsg("Killing remaining children.")
                syscall.Kill(os.Getpid(), syscall.SIGUSR1)
            }()
        } else {
            debugMsg("Reaped child with pid ", pid, " status: ", status)
        }
    }
}

// Installs the signal handler, trapping SiGINT, SIGTERM and SIGHUP
// We're not bothering about SIGKILL as we can't catch it anyways.
func installSignalHandler(cmd *exec.Cmd) chan os.Signal {
    c := make(chan os.Signal, 2)
    signal.Notify(c,
        syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGQUIT)
    go signalHandler(cmd, c)
    return c
}

// Broadcast received signals to all processes, SIGUSR1 triggers a kill signal
func signalHandler(cmd *exec.Cmd, c chan os.Signal) {
    for sig := range c {
        switch sig {
            case syscall.SIGUSR1:
                // Kill everything and exit the handler
                debugMsg("Broadcasting SIGKILL.")
                syscall.Kill(-1, syscall.SIGKILL)
                return
            default:
                // Broadcast the signal to all processes
                debugMsg("Broadcasting ", sig)
                syscall.Kill(-1, sig.(syscall.Signal))
        }
    }
}

// Run the process, route it's input and output to standard channels
func startProcess(args []string) *exec.Cmd {
    debugMsg("Starting command: ", args)
    cmd := exec.Command(args[0], args[1:]...)
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
    log.SetPrefix("[gofirst] ")

    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage: %s [options] [--] command\nOptions:\n", os.Args[0])
        flag.PrintDefaults()
    }

    flag.Parse()
    args := flag.Args()

    if len(args) == 0 {
        log.Fatal("No command supplied. ", args)
    }

    // Start the process
    cmd := startProcess(args)

    // Trap and handle signals
    c := installSignalHandler(cmd)

    // Make sure to terminate the handler when we're done
    defer close(c)

    // Wait until we're done while reaping any process that comes along
    waitAndReap(cmd, 30)
}
