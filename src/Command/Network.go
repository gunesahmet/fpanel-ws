package Command

import (
	"strings"

	"Network"
	"WebSocket"
)

func NetworkCounters() {
	collectedNetworks := make(map[string]interface{}, 2)

	for {
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
		msg["cmd"] = WebSocket.NETWORK

		collectedNetworks["upload"] = totalBytesSent * 8
		collectedNetworks["download"] = totalBytesRecv * 8
		msg["result"] = collectedNetworks
		WebSock.SendToAllClients(msg, WebSocket.ROOT, WebSocket.NETWORK, 0)
	}
}
