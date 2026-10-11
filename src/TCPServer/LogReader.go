package TCPServer

import (
	"Helpers"
	"fmt"

	//"github.com/gunesahmet/pub0sub/ops"
	//ops2 "github.com/itzmeanjan/pub0sub/ops"
	"github.com/nxadm/tail"
	"io"
)

func (s *GevServer) startDomainLog(domainId string, domainName string, sockId string) {
	fmt.Println(fmt.Sprintf("sockId: %v", sockId))
	s.Puber.StartedPublishers.Mtx.RLock()
	defer s.Puber.StartedPublishers.Mtx.RUnlock()
	logPath := s.env.PanelLogDir + "/domains/" + domainName + "/web/error.log"
	if Helpers.FileExists(logPath) {
		t, err := tail.TailFile(logPath, tail.Config{
			Follow:   true,
			Location: &tail.SeekInfo{Offset: -10, Whence: io.SeekEnd}, // <- line changed
		})
		if err != nil {
			s.Logger.Error(1007, "Fail tailing has error: %v", err)
			return
		}
		topic := domainId + "domain-log"

		//domainId + "domain-log" + client.Id
		for line := range t.Lines {
			//s.Logger.Debug(1007, "Domain log publish status: %v", StartedPublishers.Pubers[topic])
			if !s.Puber.StartedPublishers.Pubers[topic] {
				break
			}
			msg := make(map[string]interface{})
			msg["domain_id"] = domainId
			msg["domain_name"] = domainName
			msg["line"] = Helpers.Trim(line.Text)
			//s.Logger.Debug(1007, "line.Text: %s", line.Text)
			status := s.Puber.Publish("domain-log", topic, sockId, msg)
			if !status {
				s.Logger.Error(1007, "DomainLog publisher error: %v", err)
			}
		}
	} else {
		s.Logger.Error(1007, "file doesnt exists: %s", logPath)
	}
}
