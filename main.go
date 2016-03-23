package main

import (
	"flag"
	"fmt"
	"server"
)

func main() {
	ip := flag.String("ip", "127.0.0.1", "server addr")
	port := flag.String("port", "11211", "server port")
	flag.Parse()
	err := server.Run(*ip + ":" + *port)
	if err != nil {
		fmt.Println(err.Error())
	}
}
