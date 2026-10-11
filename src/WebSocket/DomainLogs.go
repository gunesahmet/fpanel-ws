package WebSocket

import (
	"fmt"
	"io"
	"strings"

	"Defines"
	"Gorm"
	"Helpers"
	"github.com/nxadm/tail"
)

func (s *WsServer) domainLog(domainId string, clientId string, userType UserType) {
	var domains Gorm.Domains
	res := Defines.DbConn.DB.First(&domains, "id = ?", domainId)
	if res.RowsAffected <= 0 {
		Helpers.HandleError(nil, 9001, fmt.Sprintf("Domains not found"))
	} else {
		logPath := Env.PanelLogDir + "/domains/" + domains.Name + "/web/error.log"
		t, err := tail.TailFile(logPath, tail.Config{
			Follow:   true,
			ReOpen:   true,
			Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}, // <- line changed
		})
		if err != nil {
			s.Logger.Error(1007, "Tailing file has an error: %v", err)
			return
		}

		startedDomainLogs.ids = append(startedDomainLogs.ids, domainId)

		for line := range t.Lines {
			if !s.anyOneWatchedThisDomain(domainId) {
				s.Logger.Debug(1007, "Nobody watched this domain log")
				t.Stop()
				return
			}
			msg := make(map[string]interface{})
			msg["result"] = line.Text
			msg["cmd"] = LOGS
			msg["sub"] = DOMAIN_LOG
			s.Logger.Debug(1007, fmt.Sprintf("line.Text: %s", line.Text))
			s.SendToAllClients(msg, userType, LOGS, DOMAIN_LOG)
		}
	}
}

func (s *WsServer) anyOneWatchedThisDomain(domainId string) bool {
	for i := range Clients {
		if strings.Compare(Clients[i].WatchedDomainId, domainId) == 0 {
			return true
		}
	}
	removeStartedDomainFromList(domainId)
	return false
}

func removeStartedDomainFromList(domainId string) {
	if len(startedDomainLogs.ids) > 0 {
		for i, v := range startedDomainLogs.ids {
			if strings.Compare(v, domainId) == 0 {
				startedDomainLogs.ids = append(startedDomainLogs.ids[:i], startedDomainLogs.ids[i+1:]...)
				break
			}
		}
	}
}

func (s *WsServer) findStartedDomainLogs(domainId string) bool {
	status := false
	if len(startedDomainLogs.ids) > 0 {
		for i := range startedDomainLogs.ids {
			if strings.Compare(startedDomainLogs.ids[i], domainId) == 0 {
				status = true
				break
			}
		}
	}
	return status
}
