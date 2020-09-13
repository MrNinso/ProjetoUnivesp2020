package podman

import (
	"os"
	"os/exec"
)

type Container struct {
	Id string
	ImageID string

	State struct {
		ExitCode int8
		Running  bool
	}
}

type Termianl struct {
	Id string
	Cmd *exec.Cmd
	TTY *os.File
}

func (t Termianl) Kill() {
	_ = t.Cmd.Process.Kill()
	_, _ = t.Cmd.Process.Wait()
	_ = t.TTY.Close()
	_ = exec.Command("podman", "stop", t.Id).Run()
}

