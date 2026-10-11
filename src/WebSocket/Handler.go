package WebSocket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"PhpSerializer"
	"PubSub"
	"TCPServer"
	wss "github.com/kataras/go-websocket"
	"github.com/yvasiyarov/php_session_decoder/php_serialize"
)

func (s *WsServer) clientHandler(conn wss.Connection) {
	var clientStr ClientsStruct
	clientStr.Mtx = new(sync.RWMutex)
	clientStr.Id = conn.ID()
	clientStr.Socket = &conn
	Clients[conn.ID()] = &clientStr

	conn.On("chat", func(message string) {
		WebSockMtx.Lock()
		defer WebSockMtx.Unlock()
		s.Logger.Debug(1007, "Client connected start: %v", conn.ID())
		Clients[conn.ID()].IsRegistered = false
		Clients[conn.ID()].Session = SessionData{}
		Clients[conn.ID()].Network = false
		Clients[conn.ID()].DomainLog = false
		Clients[conn.ID()].SysInfo = false
		Clients[conn.ID()].ProcessList = false
		Clients[conn.ID()].ServiceLog = false
		Clients[conn.ID()].SystemLog = false
		Clients[conn.ID()].WatchedServiceLogs = make(map[ServiceLogType]string)
		Clients[conn.ID()].WatchedSysLogs = make(map[SysLogType]string)
	})

	conn.OnMessage(func(message []byte) {
		WebSockMtx.Lock()
		defer WebSockMtx.Unlock()
		if Clients[conn.ID()].IsRegistered {
			var msg Message
			err := json.Unmarshal(message, &msg)
			if err == nil {
				s.commandHandler(msg, conn.ID())
			} else {
				s.Logger.Error(1007, "Incorrect message from client")
				_ = conn.Disconnect()
			}
		} else {
			msg := make(map[string]interface{})
			err := json.Unmarshal(message, &msg)
			if err == nil {
				s.registerUser(msg, conn)
			} else {
				s.Logger.Error(1007, "Incorrect message from client")
				_ = conn.Disconnect()
			}
		}
	})

	conn.OnDisconnect(func() {
		s.Logger.Info(1007, "Websocket/Handler client has disconnect")
		WebSockMtx.Lock()
		defer WebSockMtx.Unlock()
		if Clients[conn.ID()].DomainLog {
		}
		Clients[conn.ID()] = nil
		delete(Clients, conn.ID())
	})
}

func (s *WsServer) commandHandler(msg Message, clientId string) {
	Clients[clientId].Mtx.Lock()
	defer Clients[clientId].Mtx.Unlock()
	userType := Clients[clientId].Session.UserType
	switch msg.Command {
	case NETWORK:
		if userType < RESELLER {
			Clients[clientId].Network = msg.State
		}
		break
	case SYSINFO:
		if userType < RESELLER {
			Clients[clientId].SysInfo = msg.State
		}
		break
	case PROCESS_LIST:
		if userType < RESELLER {
			Clients[clientId].ProcessList = msg.State
		}
		break
	case LOGS:
		switch msg.Sub {
		case DOMAIN_LOG:
			if _, ok := msg.Data["did"]; !ok {
				s.Logger.Error(1007, "Missing parameter")
				// (*(Clients[clientId].Socket)).Disconnect()
				return
			}
			domainId := msg.Data["did"].(string)
			s.Logger.Debug(9999, "Domain log requested by client. Domain ID: %s", domainId)
			Clients[clientId].DomainLog = msg.State
			userType = Clients[clientId].Session.UserType
			Clients[clientId].WatchedDomainId = ""
			if msg.State {
				Clients[clientId].WatchedDomainId = domainId
				if !s.findStartedDomainLogs(domainId) {
					s.Logger.Debug(9999, "Nobody requested before. Domain ID: %s", domainId)
					go s.domainLog(domainId, clientId, userType)
				}
			}
			break
		case SYSTEM_LOG:
			if _, ok := msg.Data["sid"]; !ok {
				s.Logger.Error(1007, "Missing parameter")
				// (*(Clients[clientId].Socket)).Disconnect()
				return
			}
			if msg.State {
				if _, ok := msg.Data["type"]; !ok {
					s.Logger.Error(1007, "Missing parameter")
					// (*(Clients[clientId].Socket)).Disconnect()
					return
				}
			}
			if userType < RESELLER {
				status, logType := s.getSysLogType(msg.Data["type"])
				if !status {
					return
				}
				serverId := msg.Data["sid"].(string)
				Clients[clientId].SystemLog = msg.State
				userType = Clients[clientId].Session.UserType
				for i := range Clients[clientId].WatchedSysLogs {
					Clients[clientId].WatchedSysLogs[i] = ""
				}
				if msg.State {
					Clients[clientId].WatchedSysLogs[logType] = serverId
					if !s.findStartedSysLogs(serverId, logType) {
						go s.sysLog(serverId, userType, logType)
					}
				}
			}
			break
		case SERVICE_LOG:
			if _, ok := msg.Data["sid"]; !ok {
				s.Logger.Error(1007, "Missing parameter")
				// (*(Clients[clientId].Socket)).Disconnect()
				return
			}
			if msg.State {
				if _, ok := msg.Data["type"]; !ok {
					s.Logger.Error(1007, "Missing parameter")
					// (*(Clients[clientId].Socket)).Disconnect()
					return
				}
			}
			if userType < RESELLER {
				status, logType := s.getServiceType(msg.Data["type"])
				if !status {
					return
				}
				serverId := msg.Data["sid"].(string)
				Clients[clientId].ServiceLog = msg.State
				for i := range Clients[clientId].WatchedServiceLogs {
					Clients[clientId].WatchedServiceLogs[i] = ""
				}
				if msg.State {
					Clients[clientId].WatchedServiceLogs[logType] = serverId
					userType = Clients[clientId].Session.UserType
					if !s.findStartedServiceLogs(serverId, logType) {
						go s.serviceLog(serverId, userType, logType)
					}
				}
			}
			break
		}
		break
	}
}

func (s *WsServer) registerUser(msg map[string]interface{}, conn wss.Connection) {
	for key, value := range msg {
		key, _ := s.crypt.Decode(key)
		if strings.Compare(key, "ENCODED_SESS") == 0 {
			svalue, status := s.crypt.Decode(fmt.Sprintf("%v", value))
			if status {
				var sessionData SessionData
				err := json.Unmarshal([]byte(svalue), &sessionData)
				if err != nil {
					s.Logger.Error(10001, fmt.Sprintf("User validation is failed: %v", svalue))
					_ = conn.Disconnect()
					return
				}
				sess, err := Redis.Get(ctx, fmt.Sprintf("PHPREDIS_SESSION:%s", sessionData.SessionID)).Result()
				if err != nil {
					s.Logger.Error(10002, fmt.Sprintf("Getting session data fail %v, %v", err.Error(), sessionData.SessionID))
					return
				}
				decodedSess, err := PhpSerializer.Deserialize(sess)
				if err != nil {
					s.Logger.Error(10002, fmt.Sprintf("Decode data fail %v, %v", err, sess))
					return
				}
				identityMap := decodedSess["identity"].(php_serialize.PhpArray)
				userType, _ := strconv.ParseInt(identityMap["type"].(string), 10, 64)

				sessionData.UserType = UserType(userType)

				Clients[conn.ID()].IsRegistered = true
				Clients[conn.ID()].Session = sessionData
				Clients[conn.ID()].Network = false
				Clients[conn.ID()].DomainLog = false
				Clients[conn.ID()].SysInfo = false
				Clients[conn.ID()].ProcessList = false
				Clients[conn.ID()].ServiceLog = false
				Clients[conn.ID()].SystemLog = false
				Clients[conn.ID()].WatchedServiceLogs = make(map[ServiceLogType]string)
				Clients[conn.ID()].WatchedSysLogs = make(map[SysLogType]string)
				s.Logger.Debug(1007, "User registered success.")
				return
			} else {
				s.Logger.Error(1007, "User data decryption is failed.")
				_ = conn.Disconnect()
			}
		}
	}
	_ = (*Clients[conn.ID()].Socket).Disconnect()
	s.Logger.Error(1007, "User can not register. Session is incorrect. %v", msg)
}

func (s *WsServer) connectServer(ip string, key string, iv string) *TCPServer.TCPClient {
	connectedClient, err := TCPServer.NewClient().
		Servers(ip, 8890).
		Secrets(key, iv).
		Connect()
	if err != nil {
		s.Logger.Error(10001, "Cloud Client Connect Error: %v", err)
	}
	return connectedClient
}

func (s *WsServer) setChannelToClient(cmd string, clientId string) {
	// Clients[clientId].Mtx.Lock()
	// defer Clients[clientId].Mtx.Unlock()
	TopicProcessStatus.TopicMtx.Lock()
	defer TopicProcessStatus.TopicMtx.Unlock()
	if _, ok := TopicProcessStatus.Processes[cmd]; !ok {
		TopicProcessStatus.Processes[cmd] = make(map[string]bool)
	}
	if _, ok := TopicProcessStatus.Processes[cmd][clientId]; !ok {
		TopicProcessStatus.Processes[cmd][clientId] = false
	}
}

func (s *WsServer) unSetChannelToClient(cmd string, clientId string) {
	// Clients[clientId].Mtx.Lock()
	// defer Clients[clientId].Mtx.Unlock()
	TopicProcessStatus.TopicMtx.Lock()
	defer TopicProcessStatus.TopicMtx.Unlock()
	if _, ok := TopicProcessStatus.Processes[cmd]; ok {
		delete(TopicProcessStatus.Processes, cmd)
	}
}

func (s *WsServer) AddSubscriber(pub *PubSub.FPanelSubscriber) *WsServer {
	s.connectedSub = pub
	return s
}
