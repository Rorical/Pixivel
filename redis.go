package pixivel

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisSetting struct {
	MaxIdle     int
	IdleTimeout time.Duration
	Password    string
	redisURL    string
}

type RedisPool struct {
	RedisPool *redis.Pool
}

func NewRedisPool() *RedisPool {
	return &RedisPool{
		RedisPool: &redis.Pool{
			MaxIdle:     redisConf.MaxIdle,
			IdleTimeout: redisConf.IdleTimeout * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.DialURL(redisConf.redisURL)
				if err != nil {
					return nil, err
				}
				if redisConf.Password != "" {
					if _, authErr := c.Do("AUTH", redisConf.Password); authErr != nil {
						return nil, authErr
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
}

type RedisClient struct {
	conn redis.Conn
}

func (self *RedisPool) NewRedisClient() *RedisClient {
	conn := self.RedisPool.Get()
	return &RedisClient{
		conn: conn,
	}
}
func NewRedisClient() (*RedisClient, error) {
	conn, err := redis.DialURL(redisConf.redisURL)
	if err != nil {
		return nil, err
	}
	return &RedisClient{
		conn: conn,
	}, nil
}

func (self *RedisClient) GetValue(key string) (string, error) {
	value, err := redis.String(self.conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return "nil", err
		}
		return "", err
	}
	return value, nil
}

func (self *RedisClient) SetValue(key string, value string) error {
	_, err := self.conn.Do("SET", key, value)
	return err
}

func (self *RedisClient) SetExpire(key string, exp string) error {
	_, err := self.conn.Do("EXPIRE", key, exp)
	return err
}
func (self *RedisClient) GetExpire(key string) (string, error) {
	value, err := redis.String(self.conn.Do("TTL", key))
	if err != nil {
		if err == redis.ErrNil {
			return "nil", err
		}
		return "", err
	}
	return value, nil
}

func (self *RedisClient) KeyExist(key string) (bool, error) {
	exist, err := redis.Bool(self.conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	} else {
		return exist, nil
	}
}

func (self *RedisClient) BLAdd(key string) (bool, error) {
	exist, err := redis.Bool(self.conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	} else {
		return exist, nil
	}
}
