package main

import (
	// "fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <port>", os.Args[0])
	}
	if _, err := strconv.Atoi(os.Args[1]); err != nil {
		log.Fatalf("Invalid port: %s (%s)\n", os.Args[1], err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// println("--->", "sleeping")
		time.Sleep(50 * time.Millisecond)
		// println("--->", os.Args[1], req.URL.String())
		// fmt.Fprintf(w, "Hi there, I love %s!\n", req.URL.Path)
	})
	http.ListenAndServe(":"+os.Args[1], nil)
}
