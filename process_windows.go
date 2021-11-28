package main

import (
	"os/exec"
	"strconv"
)

type WindowsProcess struct {
	cmd *exec.Cmd
}

func NewProcessKillable(cmd *exec.Cmd) ProcessKillable {
	return &WindowsProcess{cmd: cmd}
}

// Kill kills the process and its children.
//
// In windows we have to execute a specific command to kill the process and its children.
func (w *WindowsProcess) Kill() error {
	kill := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(w.cmd.Process.Pid))
	return kill.Run()
}
