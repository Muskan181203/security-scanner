package services

import (
	"os"
	"os/exec"
)

func CloneRepo(repoURL string) error {

	os.RemoveAll("./repos/project")

	cmd := exec.Command(
		"git",
		"clone",
		repoURL,
		"./repos/project",
	)

	return cmd.Run()
}
