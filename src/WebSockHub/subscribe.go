package WebSockHub

import (
	"github.com/gorilla/websocket"
	"sync"
)

var SubscriptionsMtx *sync.RWMutex

type SubUser struct {
	Conn   *websocket.Conn
	Topics map[string]bool
}
type subscriptions map[string]SubUser

func (s subscriptions) subscribeUser(token string, conn *websocket.Conn) {
	SubscriptionsMtx.Lock()
	topics := make(map[string]bool)
	s[token] = SubUser{
		Conn:   conn,
		Topics: topics,
	}
	SubscriptionsMtx.Unlock()
}
func (s subscriptions) subscribeUserToTopic(token string, topic string) {
	SubscriptionsMtx.Lock()
	if _, ok := s[token]; ok {
		s[token].Topics[topic] = true
	}
	SubscriptionsMtx.Unlock()
}

func (s subscriptions) getActiveSubscribers() int {
	SubscriptionsMtx.RLock()
	defer SubscriptionsMtx.RUnlock()
	total := 0
	for _, subUser := range s {
		if len(subUser.Topics) > 0 {
			total++
		}
	}
	return total
}

func (s subscriptions) publishMessage(topic string, message []byte) int {
	sent := 0
	SubscriptionsMtx.Lock()
	defer SubscriptionsMtx.Unlock()
	for token, _ := range s {
		if ok := s[token].Topics[topic]; ok && s[token].Topics[topic] {
			err := s[token].Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				delete(s, token)
			} else {
				sent++
			}
		}
	}
	return sent
}

func (s subscriptions) unsubscribeUser(token string, topic string) {
	SubscriptionsMtx.Lock()
	defer SubscriptionsMtx.Unlock()
	if ok := s[token].Topics[topic]; ok && s[token].Topics[topic] {
		delete(s[token].Topics, topic)
	}
}

func (s subscriptions) unsubscribeUserFromAll(token string) {
	SubscriptionsMtx.Lock()
	defer SubscriptionsMtx.Unlock()
	delete(s, token)
}
