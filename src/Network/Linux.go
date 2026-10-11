//go:build linux
// +build linux

package Network

import (
	"github.com/shirou/gopsutil/net"
	"time"
)

func NetIOCountersStat_f() ([]NetIOCountersStat_j, error) {
	network, _ := net.IOCounters(true)
	time.Sleep(time.Second * 1)
	networknew, _ := net.IOCounters(true)
	timestamp := time.Now().Unix()

	var ret []NetIOCountersStat_j
	//	fmt.Println(mem)
	var info NetIOCountersStat_j
	for i, v := range networknew {
		info.Name = v.Name
		info.BytesSent = v.BytesSent
		info.BytesRecv = v.BytesRecv
		info.SpeedSent = v.BytesSent - network[i].BytesSent
		info.SpeedRecv = v.BytesRecv - network[i].BytesRecv
		info.PacketsSent = v.PacketsSent
		info.PacketsRecv = v.PacketsRecv
		info.Errin = v.Errin
		info.Errout = v.Errout
		info.Dropin = v.Dropin
		info.Dropout = v.Dropout
		info.Timestamp = timestamp
		ret = append(ret, info)
	}
	return ret, nil

}
