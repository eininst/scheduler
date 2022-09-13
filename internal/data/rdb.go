package data

import (
	"context"
	"encoding/json"
	"github.com/eininst/flog"
	"github.com/eininst/scheduler/configs"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

func NewRedisClient() *redis.Client {
	var ctx = context.TODO()
	var rconf struct {
		Addr         string `json:"addr"`
		Db           int    `json:"db"`
		PoolSize     int    `json:"poolSize"`
		MinIdleCount int    `json:"minIdleCount"`
		Password     string `json:"password"`
	}
	rstr := configs.Get("redis").String()
	_ = json.Unmarshal([]byte(rstr), &rconf)

	rcli := redis.NewClient(&redis.Options{
		Addr:         rconf.Addr,
		Password:     rconf.Password,
		DB:           rconf.Db,
		DialTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     rconf.PoolSize,
		MinIdleConns: rconf.MinIdleCount,
		PoolTimeout:  30 * time.Second,
	})
	_, err := rcli.Ping(ctx).Result()

	if err != nil {
		log.Fatal("Unbale to connect to Redis", err)
	}
	flog.With(flog.Fields{
		"addr":     rconf.Addr,
		"db":       rconf.Db,
		"poolSize": rconf.PoolSize,
	}).Debug("Connected to Redis server...")

	return rcli
}
