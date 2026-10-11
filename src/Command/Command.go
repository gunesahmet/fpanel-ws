package Command

import (
	"Environments"
	"Logger"
	"WebSocket"
)

var Env *Environments.Env
var Logging Logger.Interface
var WebSock *WebSocket.WsServer
