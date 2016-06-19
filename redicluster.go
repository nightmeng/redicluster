package redicluster

import (
	"github.com/garyburd/redigo/redis"
	"stathat.com/c/consistent"
	"sync"
	"time"
)

type RediCluster struct {
	chash    *consistent.Consistent
	redisMap map[string]*redis.Pool
	rwlock   sync.Mutex
}

func NewRediCluster(addrs []string) *RediCluster {
	chash := consistent.New()
	chash.Set(addrs)

	pools := make(map[string]*redis.Pool)
	for _, addr := range addrs {
		pools[addr] = &redis.Pool{
			MaxIdle:     512,
			MaxActive:   1024,
			IdleTimeout: 240 * time.Second,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(addr)
			},
		}
	}

	return &RediCluster{
		chash:    chash,
		redisMap: pools,
	}
}

func (rc *RediCluster) GetConn(key string) (redis.Conn, error) {
	addr, err := rc.chash.Get(key)
	if err != nil {
		return nil, err
	}

	pool, ok := rc.redisMap[addr]
	if !ok {

		rc.redisMap[addr] = pool
	}

	return pool.Get(), nil
}

func (rc *RediCluster) Close() (err error) {
	for _, pool := range rc.redisMap {
		err = pool.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
