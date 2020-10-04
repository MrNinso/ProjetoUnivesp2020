package docker

import (
	"ProjetoUnivesp2020/objets"
	"github.com/creack/pty"
	"github.com/google/uuid"
	"os/exec"
)

func StartTerminal(imageID string) (*objets.Termianl, error) {
	id := uuid.New().String()

	c := exec.Command(
		"docker",
		"run",
		"-it", "--rm", "--name", id,
		imageID,
	)

	tty, err := pty.Start(c)

	return &objets.Termianl{
		Id:  id,
		Cmd: c,
		TTY: tty,
	}, err
}

//TODO CRIAR SCRIPT PARA LIMPAR OS CONTAINERS
//func KillAllTerminals() error {
//	return exec.Command("docker", "rm", "$(docker container ls -aq)").Run()
//}
