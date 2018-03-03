package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	nats "github.com/nats-io/go-nats"
)

func processMessage(ch chan *nats.Msg) {
	fmt.Println("Subscribing to channel")
	for msg := range ch {
		fmt.Println(string(msg.Data))
	}
}

func main() {
	natsURL := "nats://localhost:4222"
	nc, _ := nats.Connect(natsURL)
	ch := make(chan *nats.Msg, 64)
	sub, err := nc.ChanSubscribe("foo", ch)
	if err != nil {
		fmt.Printf("Error subscribing to channel")
	}
	go processMessage(ch)

	// Go signal notification works by sending `os.Signal`
	// values on a channel. We'll create a channel to
	// receive these notifications (we'll also make one to
	// notify us when the program can exit).
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// This goroutine executes a blocking receive for
	// signals. When it gets one it'll print it out
	// and then notify the program that it can finish.
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	// The program will wait here until it gets the
	// expected signal (as indicated by the goroutine
	// above sending a value on `done`) and then exit.
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

	sub.Unsubscribe()
}
