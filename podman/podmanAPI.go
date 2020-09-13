package podman

import (
	"github.com/creack/pty"
	"github.com/google/uuid"
	"os/exec"
)

func StartTerminal(imageID string) (*Termianl, error) {
	id := uuid.New().String()

	c := exec.Command(
		"podman",
		"run",
		"--it", "--rm", "--name", id,
		imageID,
	)

	tty, err := pty.Start(c)

	return &Termianl {
		Id: id,
		Cmd: c,
		TTY: tty,
	}, err
}

func KillAllTerminals() error {
	return exec.Command("podman", "rm", "-a").Run()
}