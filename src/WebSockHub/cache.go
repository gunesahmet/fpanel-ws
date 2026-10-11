package WebSockHub

import (
	"Logger"
	"github.com/rs/xid"
	"sort"
	"time"
)

type Cache struct {
	messageTTL time.Duration
	messages   map[string]map[string]*CacheObject
	semaphore  chan bool
	Logger     Logger.Interface
}

type CacheObject struct {
	message interface{}
	addDate int64
}

func NewCache(TTL time.Duration, logg Logger.Interface) *Cache {
	return &Cache{
		messageTTL: TTL,
		messages:   make(map[string]map[string]*CacheObject),
		semaphore:  make(chan bool, 1),
		Logger:     logg,
	}
}

func (c *Cache) put(authToken string, message interface{}) {

	//c.Logger.Debug(1007,
	//	"caching message %s for %s",
	//	message,
	//	authToken,
	//)

	c.semaphore <- true
	messageCacheID := xid.New().String()

	messageCache := make(map[string]*CacheObject)
	if len(c.messages[authToken]) == 0 {

		messageCache[messageCacheID] = &CacheObject{
			message,
			time.Now().UnixNano(),
		}

		c.messages[authToken] = messageCache

	} else {
		messageCache = c.messages[authToken]
		messageCache[messageCacheID] = &CacheObject{
			message,
			time.Now().UnixNano(),
		}

	}

	c.messages[authToken] = messageCache

	<-c.semaphore

	go func() {
		ttlTimer := time.NewTimer(c.messageTTL)

		defer ttlTimer.Stop()

		<-ttlTimer.C

		//c.Logger.Debug(1007,
		//	"message %s for %s TTL expired, it will be evicted",
		//	message,
		//	c.messageTTL,
		//)

		c.semaphore <- true
		switch len(c.messages[authToken]) {
		case 0:
			c.Logger.Debug(1007,
				"messages for %s was already evicted, nothing to delete",
				authToken,
			)

			<-c.semaphore
			return
		case 1:
			delete(c.messages, authToken)
			<-c.semaphore

			//c.Logger.Debug(1007,
			//	"message cache for %s was purged",
			//	authToken,
			//)
			return
		default:
			delete(c.messages[authToken], messageCacheID)

			<-c.semaphore
		}

		//c.Logger.Debug(1007,
		//	"message %s successfully deleted by TTL expiration",
		//	message,
		//)
	}()

	defer func() {
		//c.Logger.Debug(1007,
		//	"message %s successfuly cached",
		//	message,
		//)
	}()
}

func (c *Cache) pop(authToken string) []*CacheObject {

	var cachedMessages []*CacheObject

	if len(c.messages[authToken]) == 0 {
		c.Logger.Debug(1007,
			"no cached messages was found for %s",
			authToken,
		)

		return cachedMessages
	}

	c.semaphore <- true

	for _, message := range c.messages[authToken] {
		cachedMessages = append(cachedMessages, message)
	}

	delete(c.messages, authToken)
	<-c.semaphore

	sort.Slice(cachedMessages, func(i, j int) bool {
		return cachedMessages[i].addDate <= cachedMessages[j].addDate
	})

	return cachedMessages
}
