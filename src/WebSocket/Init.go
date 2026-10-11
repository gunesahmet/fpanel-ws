package WebSocket

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	wss "github.com/kataras/go-websocket"
)

func (s *WsServer) Start() *WsServer {
	ws := wss.New(wss.Config{})
	startedSysLogs = make(map[SysLogType][]string)
	startedServiceLogs = make(map[ServiceLogType][]string)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wsUpgrade := ws.Upgrade(w, r)
		if err := wsUpgrade.Err(); err != nil {
			s.Logger.Error(1004, "websocket error: %v", err)
			http.Error(w, fmt.Sprintf("websocket error: %v\n", err), http.StatusServiceUnavailable)
			return
		}
		s.clientHandler(wsUpgrade)
		wsUpgrade.Wait()
	})
	if _, err := os.Stat("./fws.crt"); err == nil {
		s.Logger.Info(1005, "wss server starting on %s", fmt.Sprintf("%s:%d", s.Host, s.Port))
		if err = http.ListenAndServeTLS(fmt.Sprintf("%s:%d", s.Host, s.Port), "./fws.crt", "./fws.key",
			nil); err != nil {
			s.Logger.Error(1005, "The process exited with https error: %v", err)
		}
	} else {
		s.Logger.Info(1005, fmt.Sprintf("ws server started on %s", fmt.Sprintf("%s:%d", s.Host, s.Port)))
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.Host, s.Port), nil)
		if err != nil {
			s.Logger.Error(1005, "The process exited with http error: %v", err)
			os.Exit(1)
		}
	}

	return s
}

func (s *WsServer) changeStatusChannel(cmd string, clientId string, status bool) bool {
	s.setChannelToClient(cmd, clientId)
	TopicProcessStatus.TopicMtx.Lock()
	defer TopicProcessStatus.TopicMtx.Unlock()
	currentStatus := TopicProcessStatus.Processes[cmd][clientId]
	TopicProcessStatus.Processes[cmd][clientId] = status
	return currentStatus
}

func NewWebSocket(config *WsConfig) *WsServer {
	TopicProcessStatus.TopicMtx = new(sync.RWMutex)
	TopicProcessStatus.Processes = make(map[string]map[string]bool)
	chans := make(map[string]chan interface{})
	chans["network"] = make(chan interface{})
	chans["sysinfo"] = make(chan interface{})
	chans["syslog"] = make(chan interface{})
	chans["process"] = make(chan interface{})
	return &WsServer{
		Host:     config.Host,
		Port:     config.Port,
		Channels: chans,
	}
}
