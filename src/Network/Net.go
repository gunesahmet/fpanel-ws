package Network

import (
	"encoding/json"
)

type NetIOCountersStat_j struct {
	Name        string `json:"name"`         // interface name
	BytesSent   uint64 `json:"bytes_sent"`   // number of bytes sent
	BytesRecv   uint64 `json:"bytes_recv"`   // number of bytes received
	SpeedSent   uint64 `json:"speed_sent"`   // number of bytes sent speed
	SpeedRecv   uint64 `json:"speed_Recv"`   // number of bytes recv speed
	PacketsSent uint64 `json:"packets_sent"` // number of packets sent
	PacketsRecv uint64 `json:"packets_recv"` // number of packets received
	Errin       uint64 `json:"errin"`        // total number of errors while receiving
	Errout      uint64 `json:"errout"`       // total number of errors while sending
	Dropin      uint64 `json:"dropin"`       // total number of incoming packets which were dropped
	Dropout     uint64 `json:"dropout"`      // total number of outgoing packets which were dropped (always 0 on OSX and BSD)
	Timestamp   int64  `json:"timestamp"`    //timestamp
}

func (n NetIOCountersStat_j) String() string {
	s, _ := json.Marshal(n)
	return string(s)
}
