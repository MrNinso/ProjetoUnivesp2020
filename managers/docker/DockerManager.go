package docker

import (
	"ProjetoUnivesp2020/objets"
	"github.com/creack/pty"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"os/exec"
)

func StartTerminal(imageName string) (*objets.Termianl, error) {
	id := uuid.New().String()

	c := exec.Command(
		"docker",
		"run",
		"-it", "--rm", "--name", id,
		imageName+":1.0",
	)

	tty, err := pty.Start(c)

	return &objets.Termianl{
		Id:  id,
		Cmd: c,
		TTY: tty,
	}, err
}

func RemoveImage(dockerImageName string) error {
	return exec.Command("docker", "rmi", "-f", dockerImageName).Run()
}

func BuildImage(imageName, dockerFile string) error { //TODO TESTAR PRIMEIRO
	tmpDir, err := ioutil.TempDir("", imageName)

	if err != nil {
		return err
	}

	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	tmpfile, err := ioutil.TempFile(tmpDir, "Dockerfile")

	if err != nil {
		return err
	}

	if _, err := tmpfile.Write([]byte(dockerFile)); err != nil {
		return err
	}

	return exec.Command("docker", "build", "--file", tmpfile.Name(), "--tag", imageName+":1.0", tmpDir).Run()
}

func KillAllTerminals() error {
	return exec.Command("sh", "-c", "'docker rm $(docker container ls -aq)'").Run()
}
