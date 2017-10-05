package main

import (
	l4g "github.com/alecthomas/log4go"
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
		l4g.Crash(err)
	}
	srv.server = listner
	return srv
}
