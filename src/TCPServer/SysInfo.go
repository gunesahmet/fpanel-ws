package TCPServer

import (
	"Helpers"
	"encoding/json"
)

func (s *GevServer) startSysInfo(serverId string, sockId string) {
	s.Puber.StartedPublishers.Mtx.RLock()
	defer s.Puber.StartedPublishers.Mtx.RUnlock()
	topic := serverId + "sysinfo"

	for {
		if !s.Puber.StartedPublishers.Pubers[topic] {
			break
		}

		_, response := Helpers.Exec("/bin/sh", "-c", s.env.PanelCoreDir+"/bin/adminmng --sysinfo")
		if len(response) > 0 {
			respData := make(map[string]interface{})
			_ = json.Unmarshal([]byte(response), &respData)
			msg := make(map[string]interface{})
			msg["command"] = "sysinfo"
			msg["result"] = respData
			s.Puber.Publish("sysinfo", topic, sockId, msg)
		}
		Helpers.Sleep(5000)
	}
}
