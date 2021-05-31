// +build !windows

package proxy

import (
	"os/exec"
	"syscall"
)

func getDetachedCmd(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}
	return cmd
}
