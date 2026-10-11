package WebSockHub

import (
	"Crypt"
	"Environments"
	"Logger"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Service struct {
	Gin                  *gin.Engine
	wsUpgrder            *websocket.Upgrader
	currentSubscriptions subscriptions
	timeout              time.Duration
	sync                 chan bool
	cache                *Cache
	Logging              Logger.Interface
	env                  *Environments.Env
	crypter              *Crypt.Crypter
	ServerInfos          map[string]ServerInfos
}

type ServerInfos struct {
	Id  string
	Key string
	Iv  string
	Ip  string
	Os  string
}

var netChan chan interface{}
var sysChan chan interface{}
var procChan chan interface{}
var logChan chan interface{}

func NewService(readBufferSize int, writeBufferSize int, timeout time.Duration,
	debugMode string, cache *Cache, logg Logger.Interface, env *Environments.Env,
	sInfos map[string]ServerInfos) *Service {
	if debugMode == "DEBUG" {
		gin.SetMode(gin.ReleaseMode)
	}
	SubscriptionsMtx = new(sync.RWMutex)
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		CheckOrigin:     func(request *http.Request) bool { return true },
	}

	c, _ := Crypt.NewCrypter([]byte(env.CliUsername), []byte(env.CliPassword))
	currentSubscriptions := make(map[string]SubUser)
	return &Service{
		gin.Default(),
		upgrader,
		currentSubscriptions,
		timeout,
		make(chan bool, 1),
		cache,
		logg,
		env,
		c,
		sInfos,
	}
}

func (s *Service) SetNetworkChan(ch chan interface{}) *Service {
	netChan = ch
	return s
}
func (s *Service) SetSysInfoChan(ch chan interface{}) *Service {
	sysChan = ch
	return s
}
func (s *Service) SetLogChan(ch chan interface{}) *Service {
	logChan = ch
	return s
}
func (s *Service) SetProcessChan(ch chan interface{}) *Service {
	procChan = ch
	return s
}

func (s *Service) Run(listen string, sslCertificate string, sslPrivateKey string) {

	v1 := s.Gin.Group("v1")
	{
		v1.GET("/subscribe", s.handleSubscription)
		v1.GET("/publish", s.handlePublish)
		v1.GET("/broadcast", s.handleBroadcast)
		v1.GET("/stats", s.handleStats)
		v1.GET("/stats/sessions", s.handleStatsSessions)
		v1.GET("/health", s.handleHealth)
	}

	s.Logging.Debug(1007, "Hub Service running on %s", listen)

	if sslCertificate != "" && sslPrivateKey != "" {
		err := s.Gin.RunTLS(
			listen,
			sslCertificate,
			sslPrivateKey,
		)
		if err != nil {
			return
		}
		return
	}

	err := s.Gin.Run(listen)
	if err != nil {
		return
	}
}

func (s *Service) handleSubscription(context *gin.Context) {
	s.Logging.Debug(1007,
		"Subscription handling new request from %s",
		context.Request.RemoteAddr,
	)

	authInfo := make(chan string)
	isDisconnected := make(chan bool, 1)
	terminateAuthentication := make(chan bool, 1)
	terminateKeepAlive := make(chan bool, 1)
	incomingMessages := make(chan string)
	subscription := make(chan string)
	unsubscription := make(chan string)

	connection, err := s.wsUpgrder.Upgrade(context.Writer, context.Request, nil)

	if err != nil {
		s.Logging.Debug(1007,
			"something happening wrong during create websocket connection, reason %s",
			err.Error(),
		)
		return
	}

	defer func() {
		err := connection.Close()

		if err != nil {
			fmt.Println(err)
			return
		}

		s.Logging.Debug(1007,
			"Subscription websocket connection successfully closed for client %s",
			connection.RemoteAddr().String(),
		)
		return
	}()

	_ = connection.SetReadDeadline(time.Now().Add(s.timeout))
	_ = connection.SetWriteDeadline(time.Now().Add(s.timeout))

	var authToken string

	go s.readMessage(connection, isDisconnected, subscription, unsubscription, incomingMessages)
	go s.authenticateUser(
		incomingMessages,
		authInfo,
		terminateAuthentication,
	)

	for {
		select {
		case authToken = <-authInfo:
			s.Logging.Debug(1007,
				"Subscriber auth token received: %s",
				authToken,
			)

			go s.keepAlive(connection, s.timeout, terminateKeepAlive, s.sync)

			s.Logging.Debug(1007,
				"user %s was subscribed",
				connection.RemoteAddr().String(),
			)

			s.currentSubscriptions.subscribeUser(authToken, connection)

		case <-isDisconnected:
			s.sync <- true

			s.Logging.Debug(1007,
				"client %s with token %s will be disconnected",
				context.Request.RemoteAddr,
				authToken,
			)

			terminateAuthentication <- true
			terminateKeepAlive <- true
			s.Logging.Debug(1007,
				"terminate message was successfully sent into keepalive"+
					" and authentication processes",
			)
			s.currentSubscriptions.unsubscribeUserFromAll(authToken)
			<-s.sync
			return

		case topic := <-subscription:
			s.sync <- true
			err := connection.WriteJSON(map[string]interface{}{
				"status":   true,
				"response": fmt.Sprintf("Subscription is success channel: [%s]", topic),
			})
			if err != nil {
				s.Logging.Error(1007,
					"%s: can`t write message to websocket, reason %s",
					reflect.TypeOf(err).String(),
					err.Error(),
				)
				<-s.sync

				return
			}
			s.currentSubscriptions.subscribeUserToTopic(authToken, topic)
			s.Logging.Info(1007, "User: %s Subscribed to [%s] is success", authToken, topic)
			s.Logging.Error(1007, "Total subscribers [%d]", s.currentSubscriptions.getActiveSubscribers())
			<-s.sync
		case topic := <-unsubscription:
			s.sync <- true
			err := connection.WriteJSON(map[string]interface{}{
				"status":   true,
				"response": fmt.Sprintf("Unsubscription is success from channel: [%s]", topic),
			})
			if err != nil {
				s.Logging.Error(1007,
					"%s: can`t write message to websocket, reason %s",
					reflect.TypeOf(err).String(),
					err.Error(),
				)
				<-s.sync

				return
			}
			s.currentSubscriptions.unsubscribeUser(authToken, topic)
			s.Logging.Info(1007, "User: %s Unsubscribed to [%s] is success", authToken, topic)
			<-s.sync
		}
	}
}
func (s *Service) handlePublish(context *gin.Context) {

	authToken := strings.TrimSpace(
		context.Request.Header.Get("Token"),
	)
	s.Logging.Debug(1007,
		"Publish handling new request from %s. Token: %s",
		context.Request.RemoteAddr, authToken,
	)

	isDisconnected := make(chan bool, 1)
	terminateAuthentication := make(chan bool, 1)
	terminateKeepAlive := make(chan bool, 1)
	incomingMessages := make(chan []byte)

	connection, err := s.wsUpgrder.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		s.Logging.Error(1007,
			"Something happening wrong during create websocket connection, reason %s",
			err.Error(),
		)
		return
	}

	defer func() {
		err := connection.Close()

		if err != nil {
			fmt.Println(err)
			return
		}

		s.Logging.Debug(1007,
			"Publisher websocket connection successfully closed for client %s",
			connection.RemoteAddr().String(),
		)
		return
	}()

	go s.keepAlive(connection, s.timeout, terminateKeepAlive, s.sync)
	go s.readPublisherMessage(connection, isDisconnected, incomingMessages)

	_ = connection.SetReadDeadline(time.Now().Add(s.timeout))
	_ = connection.SetWriteDeadline(time.Now().Add(s.timeout))

	for {
		select {
		case <-isDisconnected:
			s.sync <- true

			s.Logging.Debug(1007,
				"Publisher %s with token %s disconnected",
				context.Request.RemoteAddr,
				authToken,
			)

			terminateAuthentication <- true
			terminateKeepAlive <- true
			<-s.sync
			return
		case message := <-incomingMessages:
			s.sync <- true

			publishRequest := make(map[string]interface{})
			err := json.Unmarshal(message, &publishRequest)
			if err != nil {
				context.JSON(makeInternalServerError("Message decode failed: ", err.Error()))
				return
			}

			if s.currentSubscriptions.getActiveSubscribers() > 0 {
				total := s.currentSubscriptions.publishMessage(publishRequest["chan_name"].(string), message)
				s.Logging.Debug(1007,
					"Publisher message published to [%d] subscribers chan: [%v]",
					total,
					publishRequest["chan_name"].(string),
				)
			}
			<-s.sync
		}
	}
}

func (s *Service) handleStats(context *gin.Context) {
	stat := ServiceStats{
		SubscribersCount: len(s.currentSubscriptions),
	}

	context.JSON(makeStatusOk(
		context.Request.URL.String(), stat,
	))

	return
}
func (s *Service) handleStatsSessions(context *gin.Context) {

	sessionsCount := make(map[string]int)

	s.sync <- true
	defer func() {
		<-s.sync
	}()

	//for authToken, sessions := range s.currentSubscriptions {
	//	sessionsCount[authToken] = len(s.currentSubscriptions)
	//
	//}
	//sortedCurrentSessions := newSortedSessions(sessionsCount)
	//sortedCurrentSessions.Sort()

	stat := ActiveSessions{
		SubscribersSessionsCount: sessionsCount,
	}

	context.JSON(makeStatusOk(
		context.Request.URL.String(), stat,
	))

	return
}
func (s *Service) handleHealth(context *gin.Context) {
	context.JSON(makeStatusOk(
		context.Request.URL.String(), "OK",
	))
	return
}
