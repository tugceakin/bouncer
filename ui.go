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

func getAllConfigs(w http.ResponseWriter, r *http.Request) {
	allConfigs := configStore.GetAllConfigs()
	j, err := json.Marshal(allConfigs)
	if err != nil {
		panic(err)
	}
	w.Write(j)
}

func removeConfiguration(w http.ResponseWriter, r *http.Request) {
	var configMap map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&configMap)
	if err != nil {
		panic(err)
	}
	config := configStore.GetConfig(configMap["Host"].(string), configMap["Path"].(string))
	configStore.RemoveConfig(config)
}

func updateConfiguration(w http.ResponseWriter, r *http.Request) {
	var configMap map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&configMap)
	if err != nil {
		panic(err)
	}

	config := configStore.GetConfig(configMap["host"].(string), configMap["path"].(string))
	config.TargetPath = configMap["targetPath"].(string)
	config.MaxConcurrentPerBackendServer = int(configMap["concurrency"].(float64))
	config.ReqPerSecond = int(configMap["reqPerSecond"].(float64))
	backendServerArr := make([]BackendServer, len(configMap["backendServers"].([]interface{})))

	for k, v := range configMap["backendServers"].([]interface{}) {
		backendServer := NewBackendServer(v.(map[string]interface{})["host"].(string))
		backendServerArr[k] = backendServer
	}
	config.BackendServers = backendServerArr

	configStore.UpdateConfig(config)
	config.Reload()
}

func addConfiguration(w http.ResponseWriter, r *http.Request) {
	var benchmarkMap map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&benchmarkMap)
	if err != nil {
		panic(err)
	}

	backendServerArr := make([]BackendServer, len(benchmarkMap["backendServers"].([]interface{})))

	for k, v := range benchmarkMap["backendServers"].([]interface{}) {
		backendServer := NewBackendServer(v.(map[string]interface{})["host"].(string))
		backendServerArr[k] = backendServer
	}

	concurrency, _ := strconv.Atoi(benchmarkMap["concurrency"].(string))
	reqPerSecond, _ := strconv.Atoi(benchmarkMap["reqPerSecond"].(string))
	config := NewConfig(benchmarkMap["host"].(string), backendServerArr, benchmarkMap["path"].(string), benchmarkMap["targetPath"].(string), concurrency, reqPerSecond)
	configStore.AddConfig(&config)
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
		} else { //Get host name
			log.Println(string(msg))
			log.Println("got message")
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
	nc := make(chan GlobalStatRecord)
	SubscribeConfigStats(defaultConfig, nc)
	// for conn := range connections {
	log.Println("in conn")
	go closeConnectionListener(conn, quit)
	for {
		select {
		//case astat := <-globalStatSink:
		case astat := <-nc:
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

			log.Println(len(connections))
			for conn := range connections {
				//go closeConnectionListener(conn, quit)
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
	//}

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
	http.HandleFunc("/addConfiguration", addConfiguration)
	http.HandleFunc("/removeConfiguration", removeConfiguration)
	http.HandleFunc("/updateConfiguration", updateConfiguration)
	http.HandleFunc("/getAllConfigs", getAllConfigs)
	http.ListenAndServe(listen, nil)
}
