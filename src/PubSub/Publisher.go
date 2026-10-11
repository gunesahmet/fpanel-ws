package PubSub

import (
	"Crypt"
	"Environments"
	"Logger"
	"encoding/json"
	"fmt"
	"github.com/sacOO7/gowebsocket"
	"sync"
)

type FPanelPublisher struct {
	Host              string
	Port              string
	Key               string
	Iv                string
	IsConnected       bool
	Pub               *gowebsocket.Socket
	crypter           *Crypt.Crypter
	Logger            Logger.Interface
	NetChan           chan string
	ProcessChan       chan string
	DomainLogChan     chan string
	SysInfoChan       chan string
	Env               *Environments.Env
	StartedPublishers PublishStatusStruct
}

type MessageStruct struct {
	Chan     string      `json:"chan"`
	ChanName string      `json:"chan_name"`
	SockId   string      `json:"sock_id"`
	Data     interface{} `json:"data"`
}

type PublishStatusStruct struct {
	Mtx    *sync.RWMutex
	Pubers map[string]bool
}

func NewFPanelPublisher(config *FPanelPublisher) *FPanelPublisher {
	var publishStatusStruct PublishStatusStruct
	publishStatusStruct.Mtx = new(sync.RWMutex)
	publishStatusStruct.Pubers = make(map[string]bool)
	c, _ := Crypt.NewCrypter([]byte(config.Key), []byte(config.Iv))
	return &FPanelPublisher{
		Host:              config.Host,
		Port:              config.Port,
		Key:               config.Key,
		Iv:                config.Iv,
		crypter:           c,
		StartedPublishers: publishStatusStruct,
	}
}
func (s *FPanelPublisher) SetEnv(e *Environments.Env) *FPanelPublisher {
	s.Env = e
	return s
}
func (s *FPanelPublisher) SetLogger(logg Logger.Interface) *FPanelPublisher {
	s.Logger = logg
	return s
}
func (s *FPanelPublisher) Run() *FPanelPublisher {
	sock := gowebsocket.New(fmt.Sprintf("ws://%s:%s/v1/publish", s.Host, s.Port))
	authKey, _ := s.crypter.Encode(s.Key + s.Iv)
	sock.RequestHeader.Add("Token", authKey)
	sock.OnConnected = func(socket gowebsocket.Socket) {
		s.Logger.Debug(1007, "Publisher connected to hub")
		s.IsConnected = true
	}

	sock.OnConnectError = func(err error, socket gowebsocket.Socket) {
		s.IsConnected = false
		s.Logger.Debug(1007, "Publisher Connection error to hub: %s %v", socket.Url, err)
	}

	sock.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		//@TODO message'i parse et stop veya start gelirse processleri tetikle veya durdur
		s.Logger.Debug(1007, "Publisher got message from hub: %s", message)
	}

	sock.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		s.IsConnected = false
		s.Logger.Debug(1007, "Disconnected")
		return
	}
	s.Pub = &sock
	sock.Connect()
	return s
}

func (s *FPanelPublisher) Publish(chanType string, topic string, sockId string, data interface{}) bool {
	var message MessageStruct
	message.Chan = chanType
	message.ChanName = topic
	message.SockId = sockId
	message.Data = data
	result, _ := json.Marshal(message)
	s.Pub.SendText(string(result))
	return true
}
