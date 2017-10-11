package main

import (
	"os/exec"
	"runtime"
)

func runcmd(command string) []byte {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/c"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}

	out, err := exec.Command(shell, flag, command).Output()
	if err != nil {
		panic(err)
	}

	return out
}
