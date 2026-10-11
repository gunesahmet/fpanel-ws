package PubSub

import (
	"Crypt"
	"Environments"
	"Helpers"
	"Logger"
	"encoding/json"
	"fmt"
	"github.com/sacOO7/gowebsocket"
	"strings"
	"sync"
)

type FPanelSubscriber struct {
	Host              string
	Port              string
	Key               string
	Iv                string
	IsConnected       bool
	Sub               *gowebsocket.Socket
	crypter           *Crypt.Crypter
	Env               *Environments.Env
	Logger            Logger.Interface
	NetworkChannels   NetworkChannels
	DomainLogChannels DomainLogChannels
	Mtx               *sync.RWMutex
	//SysInfoChannels   map[string]chan string
	//ProcessChannels   map[string]chan string
}

type DomainLogChannels struct {
	Channels map[string]map[string]chan interface{}
	Mtx      *sync.RWMutex
}
type NetworkChannels struct {
	Channels map[string]chan string
	Mtx      *sync.RWMutex
}

func NewFPanelSubscriber(config *FPanelSubscriber) *FPanelSubscriber {
	c, _ := Crypt.NewCrypter([]byte(config.Key), []byte(config.Iv))
	var networkChannels NetworkChannels
	var domainLogChannels DomainLogChannels
	domainLogChannels.Mtx = new(sync.RWMutex)
	//domainLogChannels.Mtx.Lock()
	//defer domainLogChannels.Mtx.Unlock()
	domainLogChannels.Channels = make(map[string]map[string]chan interface{})

	return &FPanelSubscriber{
		Host:              config.Host,
		Port:              config.Port,
		Key:               config.Key,
		Iv:                config.Iv,
		crypter:           c,
		DomainLogChannels: domainLogChannels,
		NetworkChannels:   networkChannels,
		Mtx:               new(sync.RWMutex),
	}
}
func (s *FPanelSubscriber) SetEnv(e *Environments.Env) *FPanelSubscriber {
	s.Env = e
	return s
}
func (s *FPanelSubscriber) SetLogger(logg Logger.Interface) *FPanelSubscriber {
	s.Logger = logg
	return s
}

func (s *FPanelSubscriber) Init() *FPanelSubscriber {
	sock := gowebsocket.New(fmt.Sprintf("ws://%s:%s/v1/subscribe", s.Host, s.Port))

	sock.OnConnected = func(socket gowebsocket.Socket) {
		authKey, _ := s.crypter.Encode(s.Key + s.Iv)
		socket.SendText(fmt.Sprintf("HELLO %s", authKey))
		s.IsConnected = true
	}

	sock.OnConnectError = func(err error, socket gowebsocket.Socket) {
		s.IsConnected = false
		s.Logger.Debug(1007, "Subscriber Connection error to hub: %s", socket.Url)
	}

	sock.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		decoded := message
		msg := make(map[string]interface{})
		err := json.Unmarshal([]byte(decoded), &msg)
		if err != nil {
			Helpers.HandleError(err, 1007, fmt.Sprintf("Json decode failed: %s", decoded))
		}
		if ok, _ := msg["status"].(bool); ok {
			s.Logger.Debug(1007, "Subscription is success")
		} else {
			//s.Logger.Error(1007, "Incoming decoded message from hub: %s", decoded)
			if strings.Compare(msg["chan"].(string), "domain-log") == 0 {
				s.DomainLogChannels.Mtx.Lock()
				defer s.DomainLogChannels.Mtx.Unlock()
				messageData := msg["data"]
				for subKey, _ := range s.DomainLogChannels.Channels {
					if subKey == msg["chan_name"].(string) {
						for clientId, _ := range s.DomainLogChannels.Channels[subKey] {
							//s.Logger.Info(1007, "Data setting to client: %s", clientId)
							s.DomainLogChannels.Channels[subKey][clientId] <- messageData
						}
					}
				}
			}
		}
	}

	sock.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		s.IsConnected = false
		s.Logger.Debug(1007, "Subscriber disconnected from hub")
		return
	}
	//sock.EnableLogging()
	s.Sub = &sock
	return s
}
func (s *FPanelSubscriber) Connect() *FPanelSubscriber {
	go s.Sub.Connect()
	return s
}

func (s *FPanelSubscriber) SubscribeTopic(topic string, topicType string, clientId string) bool {
	s.Mtx.Lock()
	defer s.Mtx.Unlock()
	registered := true
	if strings.Compare(topicType, "domain-log") == 0 {
		s.DomainLogChannels.Mtx.Lock()
		defer s.DomainLogChannels.Mtx.Unlock()
		if _, ok := s.DomainLogChannels.Channels[topic]; !ok {
			s.DomainLogChannels.Channels[topic] = make(map[string]chan interface{})
			registered = false
		}
		s.DomainLogChannels.Channels[topic][clientId] = make(chan interface{})
	}

	if !registered {
		s.Sub.SendText(fmt.Sprintf("subscribe:%s", topic))
	}
	return true
}

func (s *FPanelSubscriber) UnSub(topic string, clientId string) {
	s.Mtx.Lock()
	defer s.Mtx.Unlock()
	if _, ok := s.DomainLogChannels.Channels[topic]; ok {
		delete(s.DomainLogChannels.Channels[topic], clientId)
	}

	if len(s.DomainLogChannels.Channels[topic]) == 0 {
		s.Logger.Debug(1007, "Total sub client: %d", len(s.DomainLogChannels.Channels[topic]))
		s.Logger.Debug(1007, "Unsubscribing from topic %s", topic)
		s.Sub.SendText(fmt.Sprintf("unsubscribe:%s", topic))
	}
}
