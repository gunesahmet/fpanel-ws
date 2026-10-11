package TCPServer

import (
	"Crypt"
	"Defines"
	"Environments"
	"Logger"
	"PubSub"
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/gunesahmet/gev"
	"github.com/gunesahmet/gev/log"
	//"github.com/gunesahmet/pub0sub/publisher"
	"github.com/gunesahmet/toolkit/sync/atomic"
	"os"
	"strings"
	"sync"
)

type GevServer struct {
	conn     *list.List
	mu       sync.RWMutex
	server   *gev.Server
	Count    atomic.Int64
	SrvCrypt *Crypt.Crypter
	Logger   Logger.Interface
	Puber    *PubSub.FPanelPublisher
	env      *Environments.Env
}

type TServer struct {
	host   string
	port   string
	Config TConfig
	server GevServer
	Logger Logger.Interface
	Puber  *PubSub.FPanelPublisher
	env    *Environments.Env
}

type TConfig struct {
	Host string
	Port string
	Key  string
	Iv   string
}

func (s *GevServer) OnConnect(c *gev.Connection) {
	fmt.Println(" OnConnect ： ", c.PeerAddr())
	s.mu.Lock()
	s.Count.Add(1)
	log.Info("connections :", s.Count.Get())
	s.mu.Unlock()
}
func (s *GevServer) OnMessage(conn *gev.Connection, ctx interface{}, data []byte) (out interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(string(data)) > 0 {
		message := s.SrvCrypt.DecodeToString(string(data))
		if len(message) > 0 {
			var msg Defines.TcpServerMessage
			err := json.Unmarshal([]byte(message), &msg)
			if err != nil {
				s.Logger.Error(1007, "Json decode failed. Message: %s", message)
				_ = conn.Close()
				os.Exit(0)
				return
			}
			if strings.Compare(msg.Connection, "close") == 0 {
				_ = conn.Close()
			}

			s.handleMessages(msg, conn)

		} else {
			s.Logger.Error(1007, "Message decode failed")
			_ = conn.Close()
			os.Exit(0)
			return
		}
	} else {
		s.Logger.Error(1007, "Message is empty")
		_ = conn.Close()
		os.Exit(0)
		return
	}
	return
}
func (s *GevServer) OnClose(c *gev.Connection) {
	s.Count.Add(-1)
	//log.Println("OnClose")
}
func NewTcpServer(config *TConfig) *TServer {
	return &TServer{
		host: config.Host,
		port: config.Port,
		Config: TConfig{
			Key: config.Key,
			Iv:  config.Iv,
		},
	}
}
func (s *TServer) SetEnv(e *Environments.Env) *TServer {
	s.env = e
	return s
}
func (s *TServer) SetLogger(logg Logger.Interface) *TServer {
	s.Logger = logg
	return s
}
func (s *TServer) SetPublisher(pub *PubSub.FPanelPublisher) *TServer {
	s.Puber = pub
	return s
}
func (s *TServer) Run() {
	scrypt, _ := Crypt.NewCrypter([]byte(s.env.CliUsername), []byte(s.env.CliPassword))
	handler := new(GevServer)
	handler.SrvCrypt = scrypt
	handler.Logger = s.Logger
	handler.Puber = s.Puber
	handler.env = s.env
	srv, err := gev.NewServer(handler,
		gev.Network("tcp"),
		gev.Address(":"+s.port),
	)
	if err != nil {
		s.Logger.Debug(1007, "TCP Server start has an error: %v", err)
		return
	}

	srv.Start()
}

func (s *GevServer) handleMessages(message Defines.TcpServerMessage, conn *gev.Connection) {
	switch message.Command {
	case "domain-log":
		s.domainLogs(message, conn)
		err := conn.Send([]byte("\n"))
		if err != nil {
			fmt.Println(err)
			_ = conn.Close()
		}
		break
	case "network":
		s.network(message, conn)
		err := conn.Send([]byte("\n"))
		if err != nil {
			fmt.Println(err)
			_ = conn.Close()
		}
		break
	case "sysinfo":
		s.sysinfo(message, conn)
		err := conn.Send([]byte("\n"))
		if err != nil {
			fmt.Println(err)
			_ = conn.Close()
		}
		break
	case "process":
		s.process(message, conn)
		err := conn.Send([]byte("\n"))
		if err != nil {
			fmt.Println(err)
			_ = conn.Close()
		}
		break
	case "syslog":
		s.syslog(message, conn)
		err := conn.Send([]byte("\n"))
		if err != nil {
			fmt.Println(err)
			_ = conn.Close()
		}
		break
	}
}
func (s *GevServer) domainLogs(message Defines.TcpServerMessage, conn *gev.Connection) {
	if _, ok := message.Data["domain_id"]; !ok {
		s.Logger.Error(1007, "domain_id missed")
		_ = conn.Close()
		return
	}
	if _, ok := message.Data["domain_name"]; !ok {
		s.Logger.Error(1007, "domain_name missed", 1007)
		_ = conn.Close()
		return
	}
	if _, ok := message.Data["sock_id"]; !ok {
		s.Logger.Error(1007, "sock_id missed", 1007)
		_ = conn.Close()
		return
	}

	pubKey := message.Data["domain_id"] + "domain-log"

	if strings.Compare(message.Proc, "start") == 0 {
		s.Puber.StartedPublishers.Mtx.Lock()
		defer s.Puber.StartedPublishers.Mtx.Unlock()

		if ok := s.Puber.StartedPublishers.Pubers[pubKey]; !ok {
			s.Puber.StartedPublishers.Pubers[pubKey] = true
			s.Logger.Debug(1007, "Publish starting")
			go s.startDomainLog(message.Data["domain_id"], message.Data["domain_name"], message.Data["sock_id"])
			return
		}

		s.Logger.Success(1007, "PubKey status :%v", s.Puber.StartedPublishers.Pubers[pubKey])

		if s.Puber.StartedPublishers.Pubers[pubKey] {
			s.Logger.Error(1007, "Domain Log Publish already started")
		} else {
			s.Logger.Error(1007, "PubKey is not true")
		}

	} else {
		s.Puber.StartedPublishers.Pubers[pubKey] = false
		s.Logger.Debug(1007, "Publish canceled")
	}
}

func (s *GevServer) network(message Defines.TcpServerMessage, conn *gev.Connection) {
	pubKey := message.Data["sid"] + "network"

	if strings.Compare(message.Proc, "start") == 0 {
		s.Puber.StartedPublishers.Mtx.Lock()
		defer s.Puber.StartedPublishers.Mtx.Unlock()

		if ok := s.Puber.StartedPublishers.Pubers[pubKey]; !ok {
			s.Puber.StartedPublishers.Pubers[pubKey] = true
			s.Logger.Debug(1007, "Network Publish starting")
			go s.startNetwork(message.Data["sid"], message.Data["sock_id"])
			return
		}

		s.Logger.Success(1007, "Network PubKey status :%v", s.Puber.StartedPublishers.Pubers[pubKey])

		if s.Puber.StartedPublishers.Pubers[pubKey] {
			s.Logger.Error(1007, "Network Publish already started")
		} else {
			s.Logger.Error(1007, "Network PubKey is not true")
		}

	} else {
		s.Puber.StartedPublishers.Pubers[pubKey] = false
		s.Logger.Debug(1007, "Network Publish canceled")
	}
	_ = conn.Close()
}

func (s *GevServer) sysinfo(message Defines.TcpServerMessage, conn *gev.Connection) {
	pubKey := message.Data["sid"] + "sysinfo"

	if strings.Compare(message.Proc, "start") == 0 {
		s.Puber.StartedPublishers.Mtx.Lock()
		defer s.Puber.StartedPublishers.Mtx.Unlock()

		if ok := s.Puber.StartedPublishers.Pubers[pubKey]; !ok {
			s.Puber.StartedPublishers.Pubers[pubKey] = true
			s.Logger.Debug(1007, "SysInfo Publish starting")
			go s.startSysInfo(message.Data["sid"], message.Data["sock_id"])
			return
		}

		s.Logger.Success(1007, "SysInfo PubKey status :%v", s.Puber.StartedPublishers.Pubers[pubKey])

		if s.Puber.StartedPublishers.Pubers[pubKey] {
			s.Logger.Error(1007, "SysInfo Publish already started")
		} else {
			s.Logger.Error(1007, "SysInfo PubKey is not true")
		}

	} else {
		s.Puber.StartedPublishers.Pubers[pubKey] = false
		s.Logger.Debug(1007, "SysInfo Publish canceled")
	}
	_ = conn.Close()
}
func (s *GevServer) process(message Defines.TcpServerMessage, conn *gev.Connection) {
	pubKey := message.Data["sid"] + "process"

	if strings.Compare(message.Proc, "start") == 0 {
		s.Puber.StartedPublishers.Mtx.Lock()
		defer s.Puber.StartedPublishers.Mtx.Unlock()

		if ok := s.Puber.StartedPublishers.Pubers[pubKey]; !ok {
			s.Puber.StartedPublishers.Pubers[pubKey] = true
			s.Logger.Debug(1007, "Process Publish starting")
			go s.startProcess(message.Data["sid"], message.Data["sock_id"])
			return
		}

		s.Logger.Success(1007, "Process PubKey status :%v", s.Puber.StartedPublishers.Pubers[pubKey])

		if s.Puber.StartedPublishers.Pubers[pubKey] {
			s.Logger.Error(1007, "Process Publish already started")
		} else {
			s.Logger.Error(1007, "Process PubKey is not true")
		}

	} else {
		s.Puber.StartedPublishers.Pubers[pubKey] = false
		s.Logger.Debug(1007, "Process Publish canceled")
	}
	_ = conn.Close()
}
func (s *GevServer) syslog(message Defines.TcpServerMessage, conn *gev.Connection) {
}
