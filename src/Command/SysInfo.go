package Command

import (
	"encoding/json"

	"Helpers"
	"WebSocket"
)

func ReadSystemInfo() {
	for {
		err, response := Helpers.Exec("/bin/sh", "-c", Env.PanelCoreDir+"/bin/adminmng --sysinfo")

		if err != nil {
			Logging.Error(1007, "Core command execute failed. Error: %v", err)
		} else {
			if len(response) > 0 {
				respData := make(map[string]interface{})
				_ = json.Unmarshal([]byte(response), &respData)
				msg := make(map[string]interface{})
				msg["cmd"] = WebSocket.SYSINFO
				msg["result"] = respData
				WebSock.SendToAllClients(msg, WebSocket.ROOT, WebSocket.SYSINFO, 0)
			}
		}
		Helpers.Sleep(5000)
	}
}
