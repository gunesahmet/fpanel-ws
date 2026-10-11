package WebSocket

import (
	"encoding/json"
)

func (s *WsServer) SendToAllClients(msg interface{}, userType UserType, commandType int, subCommand int) {
	for clientId, _ := range Clients {
		s.SendMessageToClient(clientId, msg, userType, commandType, subCommand)
	}
}

func checkCommand(clientId string, userType UserType, commandType int, subCommand int) bool {
	if commandType == SYSINFO {
		if userType < RESELLER {
			return Clients[clientId].SysInfo
		}
	}
	if commandType == NETWORK {
		if userType < RESELLER {
			return Clients[clientId].Network
		}
	}
	if commandType == PROCESS_LIST {
		if userType < RESELLER {
			return Clients[clientId].ProcessList
		}
	}
	if commandType == LOGS {
		if subCommand == SYSTEM_LOG {
			if userType < RESELLER {
				return Clients[clientId].SystemLog
			}
		}
		if subCommand == SERVICE_LOG {
			if userType < RESELLER {
				return Clients[clientId].ServiceLog
			}
		}
		if subCommand == DOMAIN_LOG {
			return Clients[clientId].DomainLog
		}
	}
	return false
}

func (s *WsServer) SendMessageToClient(clientId string, command interface{}, userType UserType, commandType int, subCommand int) {
	WebSockMtx.RLock()
	defer WebSockMtx.RUnlock()
	if Clients[clientId] != nil {
		if !checkCommand(clientId, userType, commandType, subCommand) {
			return
		}
		c := *(Clients[clientId].Socket)
		data, _ := json.Marshal(command)
		err := c.EmitMessage(data)
		if err != nil {
			s.Logger.Error(1008, "Error send: %v", err)
			err := c.Disconnect()
			if err != nil {
				s.Logger.Error(1009, "Client socket already closed")
			}
			delete(Clients, clientId)
		}
	}
}
