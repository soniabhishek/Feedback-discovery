package main

import (
	l4g "github.com/alecthomas/log4go"
)

func main() {
	srv := InitServer()
	// accept connection on port

	// run loop forever (or until ctrl-c)
	for {
		conn, err := srv.server.Accept()
		if err != nil {
			l4g.Error(err)
		}
		go handleClient(conn)
	}
}
