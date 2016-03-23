package server

import (
	"errors"
	"fmt"
	"net"
)

var (
	StartFailed = errors.New("Server start failed")
)

func Run(addr string) error {
	Init()
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err.Error)
		return StartFailed
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
		}
		go ConnectionHandle(conn)

	}
	return nil
}
