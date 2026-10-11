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

func (s *WsServer) serviceLog(serverId string, userType UserType, servicelogType ServiceLogType) {
	var servers Gorm.Servers
	res := Defines.DbConn.DB.First(&servers, "id = ?", serverId)
	if res.RowsAffected <= 0 {
		Helpers.HandleError(nil, 9001, fmt.Sprintf("Server not found"))
	} else {
		logPath := s.getServiceLogPath(servicelogType)
		t, err := tail.TailFile(logPath, tail.Config{
			Follow:   true,
			ReOpen:   true,
			Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}, // <- line changed
		})
		if err != nil {
			s.Logger.Error(1007, "Tailing file has an error: %v", err)
			return
		}
		if len(startedServiceLogs[servicelogType]) == 0 {
			startedServiceLogs[servicelogType] = make([]string, 0)
		}
		startedServiceLogs[servicelogType] = append(startedServiceLogs[servicelogType], serverId)

		for line := range t.Lines {
			if !s.anyOneWatchedServiceLog(serverId, servicelogType) {
				s.Logger.Debug(1007, "Nobody watched this service log")
				t.Stop()
				return
			}
			msg := make(map[string]interface{})
			msg["result"] = line.Text
			msg["cmd"] = LOGS
			msg["sub"] = SERVICE_LOG
			s.Logger.Debug(1007, fmt.Sprintf("line.Text: %s", line.Text))
			s.SendToAllClients(msg, userType, LOGS, SERVICE_LOG)
		}
	}
}

func (s *WsServer) getServiceLogPath(logType ServiceLogType) string {
	logPath := Env.GeneralLogPath
	switch logType {
	case APACHE:
		if strings.Compare(Env.SystemType, "centos") == 0 {
			logPath += "/httpd/error.log"
		} else {
			logPath += "/apache2/error.log"
		}
		s.Logger.Debug(1007, "Requested APACHE")
		break
	case NGINX:
		if strings.Compare(Env.SystemType, "centos") == 0 {
			logPath += "/nginx/error.log" // TODO
		} else {
			logPath += "/nginx/error.log"
		}
		s.Logger.Debug(1007, "Requested NGINX")
		break
	case IIS:
		logPath += "/error.log" // TODO
		s.Logger.Debug(1007, "Requested IIS")
		break

	case FTP:
		if strings.Compare(Env.SystemOS, "linux") == 0 {
			logPath += "/proftpd/proftpd.log" // TODO
		} else {
			logPath += "/proftpd/proftpd.log" // TODO
		}
		s.Logger.Debug(1007, "Requested PROFTP")
		break

	case MYSQL | MARIA:
		if strings.Compare(Env.SystemType, "centos") == 0 {
			logPath += "/mysql/error.log" // TODO
		} else {
			logPath += "/mysql/error.log" // TODO
		}
		s.Logger.Debug(1007, "Requested MYSQL")
		break
	case POSTGRESQL:
		if strings.Compare(Env.SystemType, "centos") == 0 {
			logPath += "/postgresql/postgresql-12-main.log" // TODO
		} else {
			logPath += "/postgresql/postgresql-12-main.log" // TODO
		}
		s.Logger.Debug(1007, "Requested POSTGRESQL")
		break
	case MSSQL:
		logPath += "/error.log"
		s.Logger.Debug(1007, "Requested MSSQL")
		break

	case MAIL:
		if strings.Compare(Env.SystemOS, "linux") == 0 {
			logPath += "/mail/mail.log" // TODO
		} else {
			logPath += "/mysql/error.log" // TODO
		}
		s.Logger.Debug(1007, "Requested QMAIL")
		break

	case SYSTEM:
		logPath += "/error.log"
		s.Logger.Debug(1007, "Requested DMESG|WINDOWS")
		break

	case CRONTAB:
		logPath += "/error.log"
		s.Logger.Debug(1007, "Requested CRONTAB")
		break

	case ANTIVIRUS:
		logPath += "/error.log"
		s.Logger.Debug(1007, "Requested ANTIVIRUS")
		break

	case DNS:
		logPath += "/error.log"
		s.Logger.Debug(1007, "Requested DNS")
		break
	}
	return logPath
}

func (s *WsServer) anyOneWatchedServiceLog(serverId string, logType ServiceLogType) bool {
	for i := range Clients {
		if strings.Compare(Clients[i].WatchedServiceLogs[logType], serverId) == 0 {
			return true
		}
	}
	removeStartedServicelogFromList(serverId, logType)
	return false
}

func removeStartedServicelogFromList(serverId string, logType ServiceLogType) {
	if len(startedServiceLogs[logType]) > 0 {
		for i, v := range startedServiceLogs[logType] {
			if strings.Compare(v, serverId) == 0 {
				startedServiceLogs[logType] = append(startedServiceLogs[logType][:i], startedServiceLogs[logType][i+1:]...)
				break
			}
		}
	}
}

func (s *WsServer) findStartedServiceLogs(serverId string, logType ServiceLogType) bool {
	status := false
	if len(startedServiceLogs[logType]) > 0 {
		for i := range startedServiceLogs[logType] {
			if strings.Compare(startedServiceLogs[logType][i], serverId) == 0 {
				status = true
				break
			}
		}
	}
	return status
}
