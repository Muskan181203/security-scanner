package services

import (
	"fmt"

	"os/exec"
	"time"
)

func CloneRepo(repoURL string) (string, error) {

	repoPath := fmt.Sprintf(
		"./repos/scan-%d",
		time.Now().UnixNano(),
	)

	cmd := exec.Command(
		"git",
		"clone",
		repoURL,
		repoPath,
	)

	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return repoPath, nil
}
