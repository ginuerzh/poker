package main

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

const (
	redisUserTokens = "sports:user:tokens"
)

var (
	pool *redis.Pool
)

func redisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", RedisAddr)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			/*
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			*/
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func getConn() redis.Conn {
	if pool == nil {
		pool = redisPool()
	}

	return pool.Get()
}

func onlineUser(token string) (id string) {
	conn := getConn()
	id, _ = redis.String(conn.Do("HGET", redisUserTokens, token))
	return
}
