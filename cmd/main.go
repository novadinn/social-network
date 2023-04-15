package main

import (
	"log"

	"github.com/novadinn/social-network/handler"
)

var h *handler.Handler

func main() {
	h = handler.New()

	if err := h.Run(); err != nil {
		log.Fatal(err)
	}
}
