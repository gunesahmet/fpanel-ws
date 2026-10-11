package PubSub

import (
	"Crypt"
	"Helpers"

	//"Crypt"
	"Environments"
	"Logger"
	"WebSockHub"
	"fmt"
	"time"
)

type FPanelHub struct {
	Host        string
	Port        string
	Key         string
	Iv          string
	service     *WebSockHub.Service
	Env         *Environments.Env
	Crypt       *Crypt.Crypter
	Logger      Logger.Interface
	ServerInfos map[string]WebSockHub.ServerInfos
	netchan     chan interface{}
	procchan    chan interface{}
	sysinfochan chan interface{}
	logchan     chan interface{}
}

func NewHub(config *FPanelHub) *FPanelHub {
	sInfos := make(map[string]WebSockHub.ServerInfos)
	return &FPanelHub{
		Host:        config.Host,
		Port:        config.Port,
		Key:         config.Key,
		Iv:          config.Iv,
		ServerInfos: sInfos,
	}
}
func (s *FPanelHub) SetLogger(logg Logger.Interface) *FPanelHub {
	s.Logger = logg
	return s
}
func (s *FPanelHub) SetEnv(e *Environments.Env) *FPanelHub {
	s.Env = e
	crypter, _ := Crypt.NewCrypter([]byte(e.CliUsername), []byte(e.CliPassword))
	s.Crypt = crypter
	return s
}
func (s *FPanelHub) SetServerSecret(addr string, key string, iv string, osType string, serverId string) {
	serverIdx, _ := s.Crypt.Encode(key + iv)
	s.ServerInfos[serverIdx] = WebSockHub.ServerInfos{
		Id:  serverId,
		Key: key,
		Iv:  iv,
		Ip:  addr,
		Os:  osType,
	}
}
func (s *FPanelHub) Init() *FPanelHub {
	readBufferSize := 51200
	writeBufferSize := 102400
	timeout, _ := time.ParseDuration("3600s")
	cachedMessageTTL, _ := time.ParseDuration("3600s")
	service := WebSockHub.NewService(
		readBufferSize,
		writeBufferSize,
		timeout,
		"DEBUG",
		WebSockHub.NewCache(cachedMessageTTL, s.Logger),
		s.Logger,
		s.Env,
		s.ServerInfos,
	)
	s.service = service
	return s
}

func (s *FPanelHub) Run() *FPanelHub {
	go s.service.Run(
		fmt.Sprintf("%s:%s", s.Host, s.Port),
		"",
		"",
	)
	Helpers.Sleep(500)
	return s
}
