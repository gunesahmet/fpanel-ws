package WebSocket

import (
	"context"
	"sync"

	"Crypt"
	"Environments"
	"Logger"
	"PubSub"
	"github.com/go-redis/redis"
	wss "github.com/kataras/go-websocket"
	_ "golang.org/x/net/websocket"
)

var Env *Environments.Env
var Redis *redis.Client
var ctx = context.Background()
var Clients = make(map[string]*ClientsStruct) // connected clients
var TopicProcessStatus struct {
	TopicMtx  *sync.RWMutex
	Processes map[string]map[string]bool
}

type UserType int64
type SysLogType int64
type ServiceLogType int64

const (
	NONE     UserType = -1
	ROOT     UserType = 0
	ADMIN    UserType = 1
	RESELLER UserType = 2
	CLIENT   UserType = 3
	DOMAIN   UserType = 4
	EMAIL    UserType = 5
	USER     UserType = 6
)
const (
	SYSINFO      = 0
	NETWORK      = 1
	LOGS         = 2
	PROCESS_LIST = 3

	// Sub Commands
	SYSTEM_LOG  = 0
	SERVICE_LOG = 1
	DOMAIN_LOG  = 2

	// Services
	APACHE ServiceLogType = 10
	NGINX  ServiceLogType = 11
	IIS    ServiceLogType = 12

	FTP ServiceLogType = 20

	MYSQL      ServiceLogType = 30
	POSTGRESQL ServiceLogType = 31
	MARIA      ServiceLogType = 33
	MSSQL      ServiceLogType = 32

	MAIL ServiceLogType = 40

	SYSTEM ServiceLogType = 50

	CRONTAB ServiceLogType = 60

	ANTIVIRUS ServiceLogType = 80

	DNS ServiceLogType = 90

	// Sys Log Types
	PANEL_ERR     SysLogType = 100
	PANEL_INFO    SysLogType = 101
	WEBSOCKET_ERR SysLogType = 102
	FE_ERR        SysLogType = 103
	FE_INFO       SysLogType = 104
	AGENT_ERR     SysLogType = 105
	AGENT_INFO    SysLogType = 106
)

var WebSockMtx = sync.RWMutex{}

type Message struct {
	Data    map[string]interface{} `json:"data,omitempty"`
	Command int                    `json:"cmd,omitempty"`
	Sub     int                    `json:"sub,omitempty"`
	State   bool                   `json:"state,omitempty"`
}
type SessionData struct {
	SessionID string   `json:"session_id"`
	Username  string   `json:"username"`
	UserID    string   `json:"user_id"`
	UserType  UserType `json:"user_type"`
}
type ClientsStruct struct {
	Id                 string // a unique id to identify the player by the frontend
	IPAddress          string
	LocalIp            string
	Socket             *wss.Connection // websocket connection of the player
	Mtx                *sync.RWMutex
	IsWriting          bool
	Session            SessionData
	IsRegistered       bool
	Network            bool
	SysInfo            bool
	ProcessList        bool
	ServiceLog         bool
	SystemLog          bool
	DomainLog          bool
	WatchedDomainId    string
	WatchedSysLogs     map[SysLogType]string
	WatchedServiceLogs map[ServiceLogType]string
}
type WsConfig struct {
	Host string
	Port int
}
type WsServer struct {
	serverConn   wss.Connection
	Host         string
	Port         int
	Channel      chan interface{}
	Key          string
	Iv           string
	Logger       Logger.Interface
	connectedSub *PubSub.FPanelSubscriber
	crypt        *Crypt.Crypter
	Channels     map[string]chan interface{}
}
type StartedLogs struct {
	ids []string
}

var (
	startedDomainLogs  StartedLogs
	startedSysLogs     map[SysLogType][]string
	startedServiceLogs map[ServiceLogType][]string
)

func (s *WsServer) getSysLogType(str interface{}) (bool, SysLogType) {
	var logt int
	switch str.(type) {
	case int:
		logt = str.(int)
	case float64:
		logt = int(str.(float64))
	default:
		s.Logger.Error(1007, "Incorrect sys log type requested")
		return false, 1000
	}
	logType := SysLogType(logt)

	return true, logType
}

func (s *WsServer) getServiceType(str interface{}) (bool, ServiceLogType) {
	var logt int
	switch str.(type) {
	case int:
		logt = str.(int)
	case float64:
		logt = int(str.(float64))
	default:
		s.Logger.Error(1007, "Incorrect service log type requested")
		return false, 1000
	}
	servicelogType := ServiceLogType(logt)

	return true, servicelogType
}
