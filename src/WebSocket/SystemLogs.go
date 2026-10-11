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

func (s *WsServer) sysLog(serverId string, userType UserType, syslogType SysLogType) {
	s.Logger.Debug(1007, "Requested syslog")
	var servers Gorm.Servers
	res := Defines.DbConn.DB.First(&servers, "id = ?", serverId)
	if res.RowsAffected <= 0 {
		Helpers.HandleError(nil, 9001, fmt.Sprintf("Server not found"))
	} else {
		logPath := s.getLogPath(syslogType)
		t, err := tail.TailFile(logPath, tail.Config{
			Follow:   true,
			ReOpen:   true,
			Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}, // <- line changed
		})
		if err != nil {
			s.Logger.Error(1007, "Tailing file has an error: %v", err)
			return
		}
		if len(startedSysLogs[syslogType]) == 0 {
			startedSysLogs[syslogType] = make([]string, 0)
		}

		startedSysLogs[syslogType] = append(startedSysLogs[syslogType], serverId)
		for line := range t.Lines {
			if !s.anyOneWatchedSysLog(serverId, syslogType) {
				s.Logger.Debug(1007, "Nobody watched this domain log")
				t.Stop()
				return
			}
			msg := make(map[string]interface{})
			msg["result"] = line.Text
			msg["cmd"] = LOGS
			msg["sub"] = SYSTEM_LOG
			s.Logger.Debug(1007, fmt.Sprintf("line.Text: %s", line.Text))
			s.SendToAllClients(msg, userType, LOGS, SYSTEM_LOG)
		}
	}
}

func (s *WsServer) getLogPath(logType SysLogType) string {
	logPath := Env.PanelLogDir
	switch logType {
	case PANEL_ERR:
		logPath += "/error.log"
		s.Logger.Debug(1007, "Requested PANEL_ERR")
		break
	case PANEL_INFO:
		logPath += "/success.log"
		s.Logger.Debug(1007, "Requested PANEL_INFO")
		break
	case WEBSOCKET_ERR:
		logPath += "/websocket/error.log"
		s.Logger.Debug(1007, "Requested WEBSOCKET_ERR")
		break
	case FE_ERR:
		logPath = Env.PanelHttpdDir + "/var/logs/system.error.log"
		s.Logger.Debug(1007, "Requested FE_ERR")
		break
	case FE_INFO:
		logPath = Env.PanelHttpdDir + "/var/logs/system.info.log"
		s.Logger.Debug(1007, "Requested FE_INFO")
		break
	case AGENT_ERR:
		logPath += "/agent/error.log"
		s.Logger.Debug(1007, "Requested AGENT_ERR")
		break
	case AGENT_INFO:
		logPath += "/agent/info.log"
		s.Logger.Debug(1007, "Requested AGENT_INFO")
		break
	}
	return logPath
}

func (s *WsServer) anyOneWatchedSysLog(serverId string, syslogType SysLogType) bool {
	for i := range Clients {
		if strings.Compare(Clients[i].WatchedSysLogs[syslogType], serverId) == 0 {
			return true
		}
	}
	removeStartedSyslogFromList(serverId, syslogType)
	return false
}

func removeStartedSyslogFromList(serverId string, syslogType SysLogType) {
	if len(startedSysLogs[syslogType]) > 0 {
		for i, v := range startedSysLogs[syslogType] {
			if strings.Compare(v, serverId) == 0 {
				startedSysLogs[syslogType] = append(startedSysLogs[syslogType][:i], startedSysLogs[syslogType][i+1:]...)
				break
			}
		}
	}
}

func (s *WsServer) findStartedSysLogs(serverId string, syslogType SysLogType) bool {
	status := false
	if len(startedSysLogs[syslogType]) > 0 {
		for i := range startedSysLogs[syslogType] {
			if strings.Compare(startedSysLogs[syslogType][i], serverId) == 0 {
				status = true
				break
			}
		}
	}
	return status
}
