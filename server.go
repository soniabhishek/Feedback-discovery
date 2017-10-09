package main

import (
	"log"
	"net"
)

const (
	PORT = ":3333"
)

type Server struct {
	server net.Listener
}

func InitServer() *Server {
	srv := &Server{}
	listner, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal(err)
	}
	srv.server = listner
	return srv
}
