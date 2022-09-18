package data

import (
	"encoding/json"
	"github.com/eininst/rs"
	"github.com/eininst/scheduler/configs"
	"github.com/go-redis/redis/v8"
	"time"
)

type rconf struct {
	Work            int   `json:"work"`
	ReadCount       int64 `json:"readCount"`
	BlockTimeSecond int64 `json:"blockTimeSecond"`
	MaxRetries      int64 `json:"maxRetries"`
	TimeoutSecond   int64 `json:"timeoutSecond"`
}

type conf struct {
	Prefix  string          `json:"prefix"`
	Sender  rs.SenderConfig `json:"sender"`
	Receive rconf           `json:"receive"`
}

func NewRsClient(rcli *redis.Client) rs.Client {
	var c conf
	rstr := configs.Get("rs").String()
	_ = json.Unmarshal([]byte(rstr), &c)

	rcv := rs.DefaultReceiveCfg
	if c.Receive != (rconf{}) {
		r := c.Receive
		rcv = rs.ReceiveConfig{
			Work:       rs.Int(r.Work),
			MaxRetries: rs.Int64(r.MaxRetries),
			ReadCount:  rs.Int64(r.ReadCount),
			Timeout:    time.Second * time.Duration(r.TimeoutSecond),
			BlockTime:  time.Second * time.Duration(r.BlockTimeSecond),
		}
	}

	return rs.New(rcli, rs.Config{
		Prefix:  c.Prefix,
		Sender:  c.Sender,
		Receive: rcv,
	})
}
