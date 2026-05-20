package main

import (
	"log"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/server"
)

func main() {
	server, err := server.New()
	if err != nil {
		log.Fatal(err)
	}
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
