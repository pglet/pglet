// +build windows

package main

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestCommandLoop(t *testing.T) {
	pipeName := "111"
	pc, _ := newPipeImpl(pipeName)
	defer pc.close()

	go func() {
		for {
			// command echo loop
			cmd := <-pc.commands
			t.Log(cmd)

			pc.writeResult(fmt.Sprintf("%s - OK", strings.TrimSpace(strings.Split(cmd, "\n")[0])))
		}
	}()

	out, err := runPowerShell(10*time.Second, "test-command-loop.ps1", pipeName)
	t.Logf("combined out:\n%s\n", out)

	if err != nil {
		t.Fatalf("test-command-loop.ps1 failed with %s\n", err)
	}
}

func TestEventLoop(t *testing.T) {
	pipeName := "222"
	pc, _ := newPipeImpl(pipeName)
	defer pc.close()

	go func() {
		for i := 0; i < 10000; i++ {
			pc.emitEvent(fmt.Sprintf("event %d", i))
			time.Sleep(time.Millisecond * 10)
		}
	}()

	out, err := runPowerShell(10*time.Second, "test-event-loop.ps1", pipeName)
	t.Logf("combined out:\n%s\n", out)

	if err != nil {
		t.Fatalf("test-event-loop.ps1 failed with %s\n", err)
	}
}

func runPowerShell(timeout time.Duration, script string, args ...string) (out string, err error) {
	// current tests directory
	_, filename, _, _ := runtime.Caller(0)
	psPath := filepath.Join(filepath.Dir(filename), "..", "..", "tests", script)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	combinedArgs := append([]string{psPath}, args...)
	cmd := exec.CommandContext(ctx, "powershell.exe", combinedArgs...)
	outBytes, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		err = ctx.Err()
		return
	}

	if err != nil {
		return
	}
	out = string(outBytes)
	return
}
