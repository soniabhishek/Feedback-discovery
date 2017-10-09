package main

import (
	"log"
)

func main() {
	InitConfig()
	srv := InitServer()
	// accept connection on port

	// run loop forever (or until ctrl-c)
	for {
		conn, err := srv.server.Accept()
		if err != nil {
			log.Println(err)
		}
		go handleClient(conn)
	}
}
