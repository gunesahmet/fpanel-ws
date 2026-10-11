package TCPServer

import (
	"Helpers"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"strconv"
	"strings"
)

type Proc struct {
	Pid         int
	CommandName string
	CmdLine     string
	FullCommand string
	CPU         float64
	Mem         float64
	User        string
}

func getProcs() ([]Proc, error) {
	psProcs, err := process.Processes()
	if err != nil {
		//fmt.Println("widget.proc.err.gopsutil", err.Error())
	}

	procs := make([]Proc, len(psProcs))
	for i, psProc := range psProcs {
		pid := psProc.Pid
		command, err := psProc.Name()
		if err != nil {
			//sps := fmt.Sprintf("%v", psProc)
			//si := strconv.Itoa(i)
			//spid := fmt.Sprintf("%d", pid)
			//fmt.Println("widget.proc.err.getcmd", err.Error(), sps, si, spid)
		}
		cpu, err := psProc.CPUPercent()
		if err != nil {
			//sps := fmt.Sprintf("%v", psProc)
			//si := strconv.Itoa(i)
			//spid := fmt.Sprintf("%d", pid)
			//fmt.Println("widget.proc.err.cpupercent", err.Error(), sps, si, spid)
		}
		mem, err := psProc.MemoryPercent()
		if err != nil {
			//sps := fmt.Sprintf("%v", psProc)
			//si := strconv.Itoa(i)
			//spid := fmt.Sprintf("%d", pid)
			//fmt.Println("widget.proc.err.mempercent", err.Error(), sps, si, spid)
		}
		cmdLine, err := psProc.Cmdline()
		if err != nil {
			//sps := fmt.Sprintf("%v", psProc)
			//si := strconv.Itoa(i)
			//spid := fmt.Sprintf("%d", pid)
			//fmt.Println("widget.proc.err.cmdline", err.Error(), sps, si, spid)
		}
		uname, err := psProc.Username()
		if err != nil {
			//sps := fmt.Sprintf("%v", psProc)
			//si := strconv.Itoa(i)
			//spid := fmt.Sprintf("%d", pid)
			//fmt.Println("widget.proc.err.cmdline", err.Error(), sps, si, spid)
		}
		procs[i] = Proc{
			Pid:         int(pid),
			CommandName: command,
			CPU:         cpu,
			Mem:         float64(mem),
			FullCommand: command,
			CmdLine:     cmdLine,
			User:        uname,
		}
	}

	return procs, nil
}

func (s *GevServer) startProcess(serverId string, sockId string) {
	s.Puber.StartedPublishers.Mtx.RLock()
	defer s.Puber.StartedPublishers.Mtx.RUnlock()

	topic := serverId + "process"

	for {
		if !s.Puber.StartedPublishers.Pubers[topic] {
			break
		}

		procs, err := getProcs()
		if err != nil {
			fmt.Println("widget.proc.error.retrieve", err.Error())
			return
		}
		processes := make([]map[string]string, 0)
		for i := range procs {
			procData := make(map[string]string)
			//procs[i].CPU /= float64(1)
			if strings.Contains(procs[i].CommandName, "irq/") {
				continue
			}
			if strings.Contains(procs[i].CommandName, "scsi_eh") {
				continue
			}
			if strings.Contains(procs[i].CommandName, "kworker/") {
				continue
			}
			dangerCpu := "<span>"
			if procs[i].CPU > 50 {
				dangerCpu = "<span class='text-danger font-weight-bold'>"
			} else {
				if procs[i].Mem > 50 {
					dangerCpu = "<span class='text-danger font-weight-bold'>"
				} else {
					if procs[i].CPU > 15 {
						dangerCpu = "<span class='text-warning font-weight-bold'>"
					} else {
						if procs[i].Mem > 15 {
							dangerCpu = "<span class='text-warning font-weight-bold'>"
						}
					}
				}
			}
			dangerCpuEnd := "</span>"

			procData["pid"] = dangerCpu + strconv.Itoa(procs[i].Pid) + dangerCpuEnd
			procData["exe"] = dangerCpu + procs[i].CommandName + dangerCpuEnd
			procData["cmd"] = dangerCpu + procs[i].FullCommand + dangerCpuEnd
			procData["cpu"] = dangerCpu + fmt.Sprintf("%.2f", procs[i].CPU/2) + dangerCpuEnd
			procData["mem"] = dangerCpu + fmt.Sprintf("%.2f", procs[i].Mem) + dangerCpuEnd
			procData["cmd_line"] = "<span title='" + procs[i].CmdLine + "'>" + procs[i].CmdLine + dangerCpuEnd
			procData["user"] = dangerCpu + procs[i].User + dangerCpuEnd
			processes = append(processes, procData)
		}

		msg := make(map[string]interface{})
		msg["command"] = "process"
		msg["result"] = processes
		s.Puber.Publish("process", topic, sockId, msg)
		Helpers.Sleep(3000)
	}

}
