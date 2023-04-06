package redisUtil

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"time"
)

type Config struct {
	Addr     string
	Password string // set to "" if no need
	Db       int
}

type RedisClient struct {
	redisPool *redis.Pool
}

func NewRedisClient(conf Config) (*RedisClient, error) {
	redisPool := &redis.Pool{
		MaxIdle:     20,
		MaxActive:   5000,
		IdleTimeout: 120 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			if conn, err := redis.Dial("tcp", conf.Addr,
				redis.DialPassword(conf.Password),
				redis.DialDatabase(conf.Db)); err != nil {
				return nil, err
			} else {
				return conn, nil
			}
		},
	}
	if reply, err := redisPool.Get().Do("PING"); err != nil || reply.(string) != "PONG" {
		return nil, errors.New(fmt.Sprintf("Failed to connect redis, conf:%v, err:%+v, reply:%s",
			conf, err, reply))
	}
	return &RedisClient{
		redisPool: redisPool,
	}, nil
}

func (rc *RedisClient) GetConn() (redis.Conn, error) {
	return rc.GetConn()
}

func (rc *RedisClient) Close() error {
	return rc.redisPool.Close()
}
