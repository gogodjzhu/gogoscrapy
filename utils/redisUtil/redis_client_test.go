package redisUtil

import (
	"github.com/gogodjzhu/gogoscrapy/utils"
	"testing"
)

func TestRedisInit(t *testing.T) {
	conf := Config{
		Host: "ubuntu01:6379",
	}
	if err := Init(conf); err != nil {
		t.Errorf("test failed @ TestRedisInit, err:%+v", err)
	}
}

func TestGetConn(t *testing.T) {
	TestRedisInit(t)
	conn := GetConn()
	if _, err := conn.Do("SET", "name", "gogodjzhu"); err != nil {
		t.Errorf("test failed @ TestGetConn, err:%+v", err)
	}
	if res, err := conn.Do("GET", "name"); err != nil || utils.Uint8ToString(res.([]uint8)) != "gogodjzhu" {
		t.Errorf("test failed @ TestGetConn, err:%+v", err)
	}
	if _, err := conn.Do("DEL", "name", "gogodjzhu"); err != nil {
		t.Errorf("test failed @ TestGetConn, err:%+v", err)
	}
	if res, err := conn.Do("GET", "name"); err != nil || res != nil {
		t.Errorf("test failed @ TestGetConn, err:%+v", err)
	}
}

func TestPFADD(t *testing.T) {
	TestRedisInit(t)
	conn := GetConn()
	if res, err := conn.Do("PFADD", "pfadd", "gogodjzhu"); err != nil {
		t.Errorf("test failed @ TestGetConn, err:%+v", err)
	} else {
		t.Log(res)
	}
}
