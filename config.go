package main

type Config struct {
	Id                            int
	Host                          string
	Path                          string
	TargetPath                    string
	ReqPerSecond                  int
	MaxConcurrentPerBackendServer int
	BackendServers                []BackendServer
	NextBackendServer             chan BackendServer
	done                          chan struct{}
}

type BackendServer struct {
	Id       int
	Host     string
	ConfigId int
}

func (config *Config) BackendServerBootstrapRoutine() {
	for i := 0; i < config.MaxConcurrentPerBackendServer; i++ {
		for _, next := range config.BackendServers {
			config.NextBackendServer <- next
		}
	}
}

func (config *Config) Destroy() {
	config.done <- struct{}{}
}

func NewConfig(hostname string, backendServers []BackendServer, path string, targetPath string, concurrency int, reqPerSecond int) Config {
	var config Config
	config.BackendServers = backendServers
	config.NextBackendServer = make(chan BackendServer, concurrency)
	config.done = make(chan struct{})
	config.MaxConcurrentPerBackendServer = concurrency
	config.ReqPerSecond = reqPerSecond
	config.Path = path
	config.TargetPath = targetPath
	go config.BackendServerBootstrapRoutine()
	return config
}

func NewBackendServer(url string) BackendServer {
	var backendServer BackendServer
	backendServer.Host = url
	return backendServer
}
