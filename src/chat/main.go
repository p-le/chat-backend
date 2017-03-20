package main

import (
	"log"
	"net/http"
)

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func main() {
	r := newRoom()
	http.Handle("/chat", r)
	go r.run()

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
