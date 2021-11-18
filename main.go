package hootail

import (
	"log"
)

var slogs []slog
var manager wsClientManager

func init() {
	manager = wsClientManager{
		broadcast:  make(chan logLine),
		register:   make(chan *wsClient),
		unregister: make(chan *wsClient),
		clients:    make(map[*wsClient]bool),
	}
}

// register monitor log file
func Tail(logName, logPath string) {

	if logName == "" || logPath == "" {
		log.Fatal("log name and path should not be empty")
		return
	}

	for _, sl := range slogs {
		if sl.LogName == logName {
			log.Fatalf("log name has been registered: %s", logName)
			return
		}
	}
	slogs = append(slogs, slog{logName, logPath})
}

// start monitor
func Serve(port int) {
	maxPort := 1<<16 - 1
	if port <= 0 || port > maxPort {
		log.Fatalf("port should be ranged in (0, %d]", maxPort)
		return
	}

	if len(slogs) < 1 {
		log.Fatalf("no log file registered")
		return
	}

	// start to monitor all registered log files
	go manager.monitorAllLogs(slogs)

	// register HTTP handler, start listen and serve
	go manager.listenAndServe(port)

	// start websocket manager
	go manager.start()
}
