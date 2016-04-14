package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func parseBenchmarkingForm(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var benchmarkMap map[string]interface{}
	err := decoder.Decode(&benchmarkMap)
	if err != nil {
		panic(err)
	}
	log.Println(benchmarkMap["benchmarkInput"])
}

func UIServer() {

	// Sets up the handlers and listen on port 8080
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	http.Handle("/", http.FileServer(http.Dir("./templates/")))

	// Default to :8080 if not defined via environmental variable.
	var listen string = os.Getenv("LISTEN")

	if listen == "" {
		listen = ":8080"
	}

	log.Println("listening on", listen)
	http.HandleFunc("/startBenchmarking", parseBenchmarkingForm)
	http.ListenAndServe(listen, nil)
}
