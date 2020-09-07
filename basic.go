package rrr

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Basic(r Root) {
	// register
	errors := r.Register()
	if errors != nil {
		for _, err := range errors {
			log.Printf("error during root register: %v\r\n", err)
		}
		panic("cant proceed, have errors during root register")
	}

	// resolve
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	ctx, close := context.WithCancel(context.Background())
	go func() {
		<-termChan
		log.Print("got term signal, starting close")
		close()
	}()
	err := r.Resolve(ctx)
	if err != nil {
		log.Printf("cant register root: %v", err)
	}

	// release
	err = r.Release()
	if err != nil {
		log.Fatalf("err during root releasing: %v", err)
	}
	log.Print("app finish successful")
}
