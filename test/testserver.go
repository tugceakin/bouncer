package main

import (
	// "fmt"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var tokens chan struct{}
var printchan chan struct{}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <port> <concurrency>", os.Args[0])
	}
	if _, err := strconv.Atoi(os.Args[1]); err != nil {
		log.Fatalf("Invalid port: %s (%s)\n", os.Args[1], err)
	}

	tokens = make(chan struct{})
	printchan = make(chan struct{})
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {

		t := <-tokens
		time.Sleep(200 * time.Millisecond)
		fmt.Fprintf(w, "Hi there, I love %s!\n", req.URL.Path)
		tokens <- t
	})
	go ConcurrencyCounter()
	go ConcurrencyPrinter()
	http.ListenAndServe(":"+os.Args[1], nil)
}

func ConcurrencyCounter() {
	var concurrency int
	for {
		select {
		case tokens <- struct{}{}:
			concurrency++
		case <-tokens:
			concurrency--
		case <-printchan:
			log.Println("Concurrency:", concurrency)
		}
	}

}

func ConcurrencyPrinter() {
	for {
		select {
		case <-time.After(1 * time.Second):
			printchan <- struct{}{}
		}
	}
}
