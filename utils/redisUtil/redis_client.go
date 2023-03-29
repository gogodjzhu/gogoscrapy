package redisUtil

import (
	"fmt"
	"github.com/gogodjzhu/gogoscrapy/utils"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"time"
)

var LOG = utils.NewLogger()

var redisPool *redis.Pool
var isStarted bool

type Config struct {
	Host     string
	Password string // set to "" if no need
	Db       int
}

func Init(conf Config) error {
	if isStarted {
		return nil
	}
	redisPool = &redis.Pool{
		MaxIdle:     20,
		MaxActive:   5000,
		IdleTimeout: 120 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			if conn, err := redis.Dial("tcp", conf.Host,
				redis.DialPassword(conf.Password),
				redis.DialDatabase(conf.Db)); err != nil {
				return nil, err
			} else {
				return conn, nil
			}
		},
	}
	if reply, err := redisPool.Get().Do("PING"); err != nil || reply.(string) != "PONG" {
		return errors.New(fmt.Sprintf("Failed to connect redis, conf:%v, err:%+v, reply:%s",
			conf, err, reply))
	}
	return nil
}

func GetConn() redis.Conn {
	return redisPool.Get()
}

func Close() {
	if err := redisPool.Close(); err != nil {
		LOG.Warnf("failed to close redis pool, err:%+v", err)
	}
}
