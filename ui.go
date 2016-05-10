package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
)

var connections map[*websocket.Conn]bool
var statChan chan GlobalStatRecord

//For later. Will be used to parse config inputs.
func parseBenchmarkingForm(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var benchmarkMap map[string]interface{}
	err := decoder.Decode(&benchmarkMap)
	if err != nil {
		panic(err)
	}
	log.Println(benchmarkMap["configId"])
}

func closeConnectionListener(conn *websocket.Conn, quit chan *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(connections, conn)
			conn.Close()
			return
		}
		if string(msg) == "quit" {
			quit <- conn
		}
	}
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	log.Println("Succesfully upgraded connection")
	connections[conn] = true

	quit := make(chan *websocket.Conn)
	//go func() {
	for {
		select {
		case astat := <-globalStatSink:
			newStatsMap := make(map[string]interface{})
			statusCountsMap := make(map[string]int)

			newStatsMap["totalReq"] = strconv.FormatInt(astat.TotalRequests, 10)
			newStatsMap["avgRespTime"] = strconv.Itoa(int(astat.AverageResponseTime))
			newStatsMap["endTime"] = astat.EndTime.Format("15:04:05")

			//Maps that have integer keys cannot be marshalled. Create new map with string keys.
			for k, v := range astat.ResponseStatusCounts {
				statusCountsMap[strconv.Itoa(k)] = v
			}
			newStatsMap["statusCount"] = statusCountsMap

			jsonStr, _ := json.Marshal(newStatsMap)

			for conn := range connections {
				go closeConnectionListener(conn, quit)
				if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonStr)); err != nil {
					delete(connections, conn)
					conn.Close()
				}
			}
		case socketConnection := <-quit: //Put an empty struct?
			delete(connections, socketConnection)
			socketConnection.Close()
			return
		}
	}
	//}()
}

func UIServer() {
	connections = make(map[*websocket.Conn]bool)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	http.Handle("/", http.FileServer(http.Dir("./templates/")))
	http.HandleFunc("/ws", socketHandler)

	// Default to :8080 if not defined via environmental variable.
	var listen string = os.Getenv("LISTEN")

	if listen == "" {
		listen = ":8080"
	}

	log.Println("listening on", listen)
	http.HandleFunc("/startBenchmarking", parseBenchmarkingForm)
	http.ListenAndServe(listen, nil)
}
