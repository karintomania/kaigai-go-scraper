package common

import (
	"log/slog"
	"os/exec"
)

type RunGitCommandFunc func([]string) (string, error)

func RunGitCommand(options []string) (string, error) {
	cmdStr := "git"
	slog.Info("running command", "cmd", cmdStr, "options", options)

	cmd := exec.Command(cmdStr, options...)
	cmd.Dir = GetEnv("blog_repo")

	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("error running command", "cmd", cmdStr, "options", options, "output", string(output), "err", err)
		return "", err
	}

	slog.Info("succeeded running command", "output", cmdStr, "options", options)
	return string(output), nil
}
