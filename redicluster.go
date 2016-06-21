package redicluster

import (
	"github.com/garyburd/redigo/redis"
	"hash/crc32"
	"sync"
	"time"
)

type RediCluster interface {
	Get(key string) redis.Conn
	Close() error
}

type cluster struct {
	pools  []*redis.Pool
	len    uint32
	rwlock sync.Mutex
}

func NewRediCluster(addrs []string) RediCluster {
	pools := make([]*redis.Pool, len(addrs))
	for i, addr := range addrs {
		func(i int, addr string) {
			pools[i] = &redis.Pool{
				MaxIdle:     512,
				MaxActive:   1024,
				IdleTimeout: 240 * time.Second,
				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					_, err := c.Do("PING")
					return err
				},
				Dial: func() (redis.Conn, error) {
					return redis.DialURL(addrs[i])
				},
			}
		}(i, addr)
	}

	return &cluster{
		pools: pools,
		len:   uint32(len(pools)),
	}
}

func (c *cluster) Get(key string) redis.Conn {
	return c.pools[crc32.ChecksumIEEE([]byte(key))%c.len].Get()
}

func (c *cluster) Close() (err error) {
	for _, pool := range c.pools {
		err = pool.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
