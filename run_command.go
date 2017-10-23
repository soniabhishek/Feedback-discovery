package main

import (
	"os/exec"
	"log"
	"runtime"
)

func runcmd(command string) (res []byte) {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/c"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}

	res, err := exec.Command(shell, flag, command).Output()
	if err != nil {
	log.Println("err ", err) 
		return
	}

	return
}
