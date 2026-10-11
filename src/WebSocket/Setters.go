package WebSocket

import (
	"Crypt"
	"Environments"
	"Logger"
	"github.com/go-redis/redis"
)

func (s *WsServer) SetEnv(env *Environments.Env) *WsServer {
	Env = env
	s.crypt, _ = Crypt.NewCrypter([]byte(Env.CliUsername), []byte(Env.CliPassword))
	return s
}

func (s *WsServer) SetRedis(redis *redis.Client) *WsServer {
	Redis = redis
	return s
}

func (s *WsServer) SetLogger(logg Logger.Interface) *WsServer {
	s.Logger = logg
	return s
}
