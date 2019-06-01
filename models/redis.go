package models

import (
	"fmt"
	"time"

	log "github.com/cihub/seelog"
	rlib "github.com/garyburd/redigo/redis"
)

var (
	redisPool *rlib.Pool
)

func Open(host string, port int, password string) {
	redisPool = newPool(host, port, password)
}

func GetConn() rlib.Conn {
	return redisPool.Get()
}

func Close() error {
	return redisPool.Close()
}

func newPool(host string, port int, password string) *rlib.Pool {
	return &rlib.Pool{
		MaxIdle:     19999,
		IdleTimeout: 240 * time.Second,
		Dial: func() (rlib.Conn, error) {
			c, err := rlib.Dial("tcp", fmt.Sprintf("%s:%d", host, port), rlib.DialPassword(password))
			if err != nil {
				log.Error("failed to dial redis server:", err)
				return nil, err
			}
			return c, err
		},
	}
}

func LoadCap(shortcode, key string) (uint64, error) {
	conn := GetConn()
	defer func() {
		conn.Close()
	}()
	return rlib.Uint64(conn.Do("GET", "cap_"+shortcode+"_"+key))
}

func SetCap() {
	conn := GetConn()
	defer func() {
		conn.Close()
	}()
	conn.Do("SET", "cap_4556066_FY1", 0)
	conn.Do("SET", "cap_4556067_GZ", 0)
	conn.Do("SET", "cap_4556067_GY", 0)
	conn.Do("SET", "cap_4556068_WZ", 0)

}

func IncrCap(shortcode, key string) {
	conn := GetConn()
	defer func() {
		conn.Close()
	}()
	conn.Do("INCRBY", "cap_"+shortcode+"_"+key, 1)
}

func LoadPostback(camp_id string) (uint64, error) {
	conn := GetConn()
	defer func() {
		conn.Close()
	}()
	return rlib.Uint64(conn.Do("GET", "postback_kc_"+camp_id))
}

func SetPostback() {
	conn := GetConn()
	defer func() {
		conn.Close()
	}()
	conn.Do("SET", "postback_kc_10000", 0)
	conn.Do("SET", "postback_kc_10001", 0)
	conn.Do("SET", "postback_kc_10002", 0)
	conn.Do("SET", "postback_kc_10003", 0)
	conn.Do("SET", "postback_kc_10004", 0)
	conn.Do("SET", "postback_kc_10005", 0)
	conn.Do("SET", "postback_kc_10006", 0)
	conn.Do("SET", "postback_kc_10007", 0)
	conn.Do("SET", "postback_kc_10008", 0)

	conn.Do("SET", "postback_kc_21101", 0)
	conn.Do("SET", "postback_kc_21102", 0)
	conn.Do("SET", "postback_kc_21103", 0)
	conn.Do("SET", "postback_kc_21104", 0)
	conn.Do("SET", "postback_kc_21105", 0)
	conn.Do("SET", "postback_kc_21106", 0)
	conn.Do("SET", "postback_kc_21107", 0)
	conn.Do("SET", "postback_kc_21108", 0)
	conn.Do("SET", "postback_kc_21109", 0)
	conn.Do("SET", "postback_kc_21110", 0)
	conn.Do("SET", "postback_kc_21111", 0)
	conn.Do("SET", "postback_kc_21112", 0)
	conn.Do("SET", "postback_kc_21113", 0)

	conn.Do("SET", "postback_kc_5601", 0)
	conn.Do("SET", "postback_kc_5602", 0)
	conn.Do("SET", "postback_kc_5603", 0)
	conn.Do("SET", "postback_kc_5604", 0)
	conn.Do("SET", "postback_kc_5605", 0)
	conn.Do("SET", "postback_kc_5606", 0)
	conn.Do("SET", "postback_kc_5608", 0)
	conn.Do("SET", "postback_kc_5609", 0)
	conn.Do("SET", "postback_kc_5610", 0)
	conn.Do("SET", "postback_kc_5611", 0)
	conn.Do("SET", "postback_kc_5612", 0)
	conn.Do("SET", "postback_kc_5613", 0)
	conn.Do("SET", "postback_kc_5614", 0)
	conn.Do("SET", "postback_kc_5615", 0)
	conn.Do("SET", "postback_kc_5616", 0)
	conn.Do("SET", "postback_kc_5617", 0)
	conn.Do("SET", "postback_kc_5618", 0)
	conn.Do("SET", "postback_kc_5619", 0)

}

func IncrPostback(camp_id string) {
	conn := GetConn()
	defer func() {
		conn.Close()
	}()
	conn.Do("INCRBY", "postback_kc_"+camp_id, 1)
}
