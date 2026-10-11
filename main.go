package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"Command"
	"Crypt"
	"Defines"
	"Environments"
	"Helpers"
	"Logger"
	"PubSub"
	"Redis"
	"TCPServer"
	"WebSocket"
	"github.com/firstrow/tcp_server"
	"github.com/sevlyar/go-daemon"
)

var (
	addr    = flag.String("addr", ":8890", "WS Service Address")
	signal  = flag.Bool("q", false, "Terminate process")
	debug   = flag.Bool("i", false, "Print info level logs to stdout")
	systemD = flag.Bool("D", false, "Run for systemd")

	pidPath = "/var/run/fws.pid"

	stop    = make(chan struct{})
	done    = make(chan struct{})
	Env     *Environments.Env
	puber   *PubSub.FPanelPublisher
	suber   *PubSub.FPanelSubscriber
	Logging Logger.Interface
)

var webSock *WebSocket.WsServer

// var networkChan chan interface{}
// var processChan chan interface{}
// var sysInfoChan chan interface{}
// var logChan chan interface{}

func main() {
	daemon.AddCommand(daemon.BoolFlag(signal), syscall.SIGQUIT, termHandler)
	// defer profile.Start(profile.MemProfile).Stop()
	Environments.Init()
	Env = Environments.GetConfig()

	Logging = Environments.Logging
	Helpers.Logging = Logging
	flag.Parse()

	if !*systemD {
		if strings.Compare(Env.SystemOS, "windows") == 0 {
			pidPath = fmt.Sprintf("%s/var/fws.pid", Env.PanelDir)
		}
		cntxt := &daemon.Context{
			PidFileName: pidPath,
			PidFilePerm: 0644,
			WorkDir:     "./",
			Umask:       027,
			Args:        []string{},
		}

		if len(daemon.ActiveFlags()) > 0 {
			d, err := cntxt.Search()
			if err != nil {
				Logging.Error(1001, "Unable send signal to the daemon:")
			}

			defer func() {
				err := daemon.SendCommands(d)
				if err != nil {
					os.Exit(0)
				}
			}()

			return
		}

		d, err := cntxt.Reborn()
		if err != nil {
			Logging.Error(1002, "Error: %v", err)
		}
		if d != nil {
			return
		}
		defer func() {
			err := cntxt.Release()
			if err != nil {
				os.Exit(0)
			}
		}()

	}

	// test()
	startServer()
	Logging.Info(1007, "All services started")
	if !*systemD {
		err := daemon.ServeSignals()
		if err != nil {
			Logging.Error(1003, "Error: %v", err)
		}
	}
}

func termHandler(sig os.Signal) error {
	Logging.Error(1002, "Terminating")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	if daemon.ErrStop != nil {
		Logging.Debug(1007, "There is no process for stop")
	}
	return nil
}

func startServer() {
	Defines.Env = Env
	Defines.Logging = Logging
	if !Env.IsCloud {
		err := Defines.ConnectDatabase()
		Helpers.Sleep(1000)
		if err != nil {
			Helpers.HandleError(err, 1001, "Database connection failed.")
			os.Exit(1)
		}
	}
	Process()
	if Env.IsCloud {
		Logging.Info(1004, "Starting as Server Feed")
		serverTcp()
	} else {
		Logging.Info(1004, "Starting as WebSocket Server")
		serveWs()
	}
}

func Process() {
	Helpers.Sleep(1000)
	Command.Env = Env
	Command.Logging = Logging

	if Env.IsCloud {
		puber = PubSub.NewFPanelPublisher(&PubSub.FPanelPublisher{
			Host: Env.CloudMainServer,
			Port: "8891",
			Key:  Env.CliUsername,
			Iv:   Env.CliPassword,
		}).
			SetEnv(Env).
			SetLogger(Logging)

		puber.Run()
		Helpers.Sleep(1000)
	}
}

func serveWs() {
	var ctx = context.Background()
	Redis.ConnectRedis(Env.RedisAuth)
	_, err := Redis.RedisClient.Ping(ctx).Result()
	if err != nil {
		Logging.Error(1004, "Failed to connect redis: %v", err)
		os.Exit(0)
	}
	flag.Parse()

	if Env.IsCloud {
		if Env.IsCloudMain {
			webSock = WebSocket.NewWebSocket(&WebSocket.WsConfig{
				Host: "0.0.0.0",
				Port: 8890,
			}).
				SetEnv(Env).
				SetRedis(Redis.RedisClient).
				SetLogger(Logging)

			suber = PubSub.NewFPanelSubscriber(&PubSub.FPanelSubscriber{
				Host: "127.0.0.1",
				Port: "8891",
				Key:  Env.CliUsername,
				Iv:   Env.CliPassword,
			}).
				SetEnv(Env).
				SetLogger(Logging).
				Init().
				Connect()
			Helpers.Sleep(1000)
			webSock.AddSubscriber(suber)
		} else {
			PubSub.NewHub(&PubSub.FPanelHub{
				Host: "0.0.0.0",
				Port: "8891",
				Key:  Env.CliUsername,
				Iv:   Env.CliPassword,
			}).
				SetEnv(Env).
				SetLogger(Logging).
				Init().
				Run()
			Logging.Debug(1007, "FPanel Hub started")
		}
	} else {
		webSock = WebSocket.NewWebSocket(&WebSocket.WsConfig{
			Host: "0.0.0.0",
			Port: 8890,
		}).
			SetEnv(Env).
			SetRedis(Redis.RedisClient).
			SetLogger(Logging)
	}

	Command.WebSock = webSock
	go Command.ReadSystemInfo()
	go Command.NetworkCounters()
	go Command.ReadProcesses()
	go Command.StartScreenshots()

	webSock.Start()
}

func serverTcp() {
	server := TCPServer.NewTcpServer(&TCPServer.TConfig{
		Host: "0.0.0.0",
		Port: "8890",
		Key:  Env.CliUsername,
		Iv:   Env.CliPassword,
	}).
		SetEnv(Env).
		SetLogger(Logging).
		SetPublisher(puber)
	server.Run()
}

func connectTcpServer(tcpHost string, tcpPort string) {
	Logging.Info(1004, "Tcp server started on %s:%s", tcpHost, tcpPort)
	server := tcp_server.New(tcpHost + ":" + tcpPort)

	server.OnNewClient(func(c *tcp_server.Client) {
		// Helpers.LogInfo.Info("Connected a client")
	})
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		// Messages.Mtx.Lock()
		// defer Messages.Mtx.Unlock()
		// Helpers.LogDebug.Info("Request Data: " + message)
		// err := c.Send("{\"send\":true}")

		// if err != nil {
		//	Helpers.LogError.Error("Can not send message to client: ", err)
		// }
		// Messages.ParseCommands(c, message)
		// c.Close()
		// c.Conn().Close()
	})

	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		c.Close()
	})
	server.Listen()
}

func test() {
	str := "{\"command\":\"domain-log\",\"proc\":\"start\",\"connection\":\"open\"," +
		"\"data\":{\"domain_id\":\"4444E91D-3D61-543C-4BA3-A9C6510F5AF3\",\"domain_name\":\"domain.com\",\"sock_id\":\"056279da-a729-49be-ab9c-10d3651ca531\"}}"

	clientCrypter, _ := Crypt.NewCrypter([]byte("60a5b17150ffa83c63db8e1b9324f5c9"), []byte("8a53de3d0d106725"))
	command := clientCrypter.EncodeToString(str)
	cmd := clientCrypter.DecodeToString(command)
	fmt.Println("command: " + command)
	fmt.Println("cmd    : " + cmd)
	os.Exit(0)
}

func test1(test1C chan string) {
	time.Sleep(time.Second * 2)
	t := <-test1C
	fmt.Println(t)
}

func test2(test2C chan string) {
	time.Sleep(time.Second * 1)
	t := <-test2C
	fmt.Println(t)
}
func test3(ch chan string, content string) {
	ch <- content
}
