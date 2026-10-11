package TCPServer

import (
	"Network"
	"strings"
)

func (s *GevServer) startNetwork(serverId string, sockId string) {
	s.Puber.StartedPublishers.Mtx.RLock()
	defer s.Puber.StartedPublishers.Mtx.RUnlock()

	collectedNetworks := make(map[string]interface{}, 2)
	topic := serverId + "network"

	for {
		if !s.Puber.StartedPublishers.Pubers[topic] {
			break
		}

		network, _ := Network.NetIOCountersStat_f()
		var totalBytesRecv uint64
		var totalBytesSent uint64
		for _, v := range network {
			if strings.Compare(v.Name, "lo") == 0 {
				continue
			}
			totalBytesRecv += v.SpeedRecv
			totalBytesSent += v.SpeedSent
		}

		msg := make(map[string]interface{})
		msg["command"] = "network-status"

		collectedNetworks["upload"] = totalBytesSent * 8
		collectedNetworks["download"] = totalBytesRecv * 8
		msg["result"] = collectedNetworks
		s.Puber.Publish("network", topic, sockId, msg)
	}
}
