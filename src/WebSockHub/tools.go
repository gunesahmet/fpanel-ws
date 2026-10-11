package WebSockHub

import (
	"Helpers"
	"github.com/gorilla/websocket"
	"reflect"
	"sort"
	"strings"
	"time"
)

const helloMessage = "HELLO"
const subscribeMessage = "subscribe:"
const ubSubscribeMessage = "unsubscribe:"

func (s *Service) readMessage(connection *websocket.Conn, isDisconnected chan bool,
	subscription chan string, unsubscription chan string, incomingMessages chan string) {
	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			switch err.(type) {
			case *websocket.CloseError:
				s.Logging.Debug(1007,
					"Handle CLOSE frame, subscriber [%s] will be disconnected",
					connection.RemoteAddr().String(),
				)
				isDisconnected <- true

			default:
				s.Logging.Debug(1007,
					"Subscriber sent unexpected error type [%s] with message [%s]",
					reflect.TypeOf(err).String(),
					err.Error(),
				)
				isDisconnected <- true
			}
			return
		}

		switch messageType {
		case websocket.PingMessage, websocket.PongMessage:
			s.Logging.Debug(1007,
				"got PING-PONG frame, ignore it")
			continue

		case websocket.TextMessage:
			if strings.HasPrefix(string(message), subscribeMessage) {
				channel := Helpers.Trim(Helpers.RemoveAll(string(message), subscribeMessage))
				s.Logging.Debug(1007,
					"Subscription User subscribed channel: [%s]. Token: ...",
					channel,
				)
				subscription <- channel
				continue
			} else if strings.HasPrefix(string(message), ubSubscribeMessage) {
				channel := Helpers.Trim(Helpers.RemoveAll(string(message), ubSubscribeMessage))
				s.Logging.Debug(1007,
					"Subscription User unsubscribed channel: [%s]. Token: ...",
					channel,
				)
				unsubscription <- channel
				continue
			} else {
				incomingMessages <- string(message)
			}

		default:
			s.Logging.Debug(1007,
				"Subscriber unexpected message type %d",
				messageType,
			)
		}
	}
}
func (s *Service) readPublisherMessage(connection *websocket.Conn, isDisconnected chan bool,
	incomingMessages chan []byte) {
	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			switch err.(type) {
			case *websocket.CloseError:
				s.Logging.Debug(1007,
					"Handle CLOSE frame, publisher [%s] will be disconnected",
					connection.RemoteAddr().String(),
				)
				isDisconnected <- true

			default:
				s.Logging.Error(1007,
					"Publisher sent unexpected error type [%s] with message [%s]",
					reflect.TypeOf(err).String(),
					err.Error(),
				)
				isDisconnected <- true
			}
			return
		}

		switch messageType {
		case websocket.PingMessage, websocket.PongMessage:
			s.Logging.Debug(1007,
				"got PING-PONG frame, ignore it")
			continue

		case websocket.TextMessage:
			incomingMessages <- message

		default:
			s.Logging.Debug(1007,
				"Publisher unexpected message type %d",
				messageType,
			)
		}
	}
}

func (s *Service) authenticateUser(incomingMessages chan string, isAuthenticated chan string,
	terminateAuthentication chan bool) {
	isAlreadyAuthenticated := false
	for {
		select {
		case <-terminateAuthentication:
			s.Logging.Debug(1007,
				"authentication process terminated",
			)
			return

		case incomingMessage := <-incomingMessages:
			//s.Logging.Debug(1007,
			//	"Got new message from subscriber: %s",
			//	incomingMessage,
			//)

			if strings.HasPrefix(incomingMessage, helloMessage) {
				if isAlreadyAuthenticated {
					s.Logging.Debug(1007,
						"Session already authenticated",
					)
					continue
				}

				tokens := strings.Fields(incomingMessage)
				if len(tokens) != 2 {
					s.Logging.Debug(1007,
						"Unexpected HELLO message: %s",
						incomingMessage,
					)
					continue
				}

				authToken := strings.TrimSpace(tokens[1])
				//s.Logging.Debug(1007,
				//	"Subscription User authenticated. Token: %s",
				//	authToken,
				//)

				isAuthenticated <- authToken
				isAlreadyAuthenticated = true

				continue
			}
		}
	}
}

func (s *Service) authenticatePublisher(incomingMessages chan string, isAuthenticated chan string,
	terminateAuthentication chan bool) {
	isAlreadyAuthenticated := false
	for {
		select {
		case <-terminateAuthentication:
			s.Logging.Debug(1007,
				"Publisher authentication process terminated",
			)
			return

		case incomingMessage := <-incomingMessages:
			s.Logging.Debug(1007,
				"Publisher sent new message: %s",
				incomingMessage,
			)

			if strings.HasPrefix(incomingMessage, helloMessage) {
				if isAlreadyAuthenticated {
					s.Logging.Debug(1007,
						"Publisher session already authenticated",
					)
					continue
				}

				tokens := strings.Fields(incomingMessage)
				if len(tokens) != 2 {
					s.Logging.Debug(1007,
						"Publisher sent unexpected HELLO message: %s",
						incomingMessage,
					)
					continue
				}

				authToken := strings.TrimSpace(tokens[1])
				s.Logging.Debug(1007,
					"Publisher authenticated. Token: %s",
					authToken,
				)

				isAuthenticated <- authToken
				isAlreadyAuthenticated = true

				continue
			}
		}
	}
}

func (s *Service) keepAlive(connection *websocket.Conn, timeout time.Duration, terminateKeepAlive chan bool, sync chan bool) {
	lastResponse := time.Now()
	keepaliveStartTime := lastResponse

	connection.SetPongHandler(func(msg string) error {
		//s.Logging.Debug(1007,
		//	"handle PONG frame from %s with message %s",
		//	connection.RemoteAddr().String(),
		//	msg,
		//)

		lastResponse = time.Now()
		return nil
	})

	for {
		select {
		case <-terminateKeepAlive:
			s.Logging.Debug(1007,
				"keepalive process of %s terminated",
				connection.RemoteAddr().String(),
			)
			return

		default:
			_ = connection.SetReadDeadline(time.Now().Add(timeout))
			_ = connection.SetWriteDeadline(time.Now().Add(timeout))

			//s.Logging.Debug(1007,
			//	"sending PING frame to client %s",
			//	connection.RemoteAddr().String(),
			//)

			sync <- true
			err := connection.WriteMessage(
				websocket.PingMessage,
				[]byte("ping"),
			)
			<-sync
			if err != nil {
				s.Logging.Debug(1007,
					"unable to write PING frame to %s, reason %s",
					connection.RemoteAddr().String(),
					err.Error(),
				)

				return
			}

			time.Sleep(timeout / 2)

			if time.Now().Sub(lastResponse) > timeout {
				s.Logging.Debug(1007,
					"ping timeout exceeded from %s, keepalive session length %s\n",
					connection.RemoteAddr().String(),
					time.Now().Sub(keepaliveStartTime).String(),
				)

				return
			}
		}
	}
}

type sortedSessions struct {
	AuthTokens []string
	Sessions   []int
}

func newSortedSessions(rawSessions map[string]int) *sortedSessions {
	sessions := &sortedSessions{
		AuthTokens: make([]string, 0, len(rawSessions)),
		Sessions:   make([]int, 0, len(rawSessions)),
	}

	for authToken, sess := range rawSessions {
		sessions.AuthTokens = append(sessions.AuthTokens, authToken)
		sessions.Sessions = append(sessions.Sessions, sess)
	}
	return sessions
}

func (sessions *sortedSessions) Sort() {
	sort.Sort(sessions)
}

func (sessions *sortedSessions) Len() int {
	return len(sessions.Sessions)
}

func (sessions *sortedSessions) Less(i, j int) bool {
	return sessions.Sessions[i] < sessions.Sessions[j]
}

func (sessions *sortedSessions) Swap(i, j int) {
	sessions.Sessions[i], sessions.Sessions[j] =
		sessions.Sessions[j], sessions.Sessions[i]

	sessions.AuthTokens[i], sessions.AuthTokens[j] =
		sessions.AuthTokens[j], sessions.AuthTokens[i]
}
