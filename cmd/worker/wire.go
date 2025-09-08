package main

import (
	"log"
	"time"
)

func run() int {
	log.Println("worker ready (stub)")
	// TODO: implement expiration scanner
	time.Sleep(1 * time.Second)
	return 0
}
