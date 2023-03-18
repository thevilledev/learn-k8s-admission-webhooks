package main

import "github.com/thevilledev/learn-admission-controllers/pkg/server"

func main() {
	s := server.New()
	s.ListenAddr(":443")
}
