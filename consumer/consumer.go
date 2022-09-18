package consumer

import (
	"github.com/eininst/rs"
	"github.com/eininst/scheduler/consumer/subscribe"
	"github.com/gofiber/fiber/v2/utils"
	"strings"
)

type Consumer struct {
	Cli rs.Client                `inject:""`
	Ts  *subscribe.TaskSubscribe `inject:""`
}

func (cs *Consumer) Init() {
	var gid = func() string { return strings.Replace(utils.UUIDv4(), "-", "", -1) }

	cs.Cli.Receive(rs.Rctx{
		Stream:  "task_register",
		Group:   gid(),
		Handler: cs.Ts.Register,
	})

	cs.Cli.Receive(rs.Rctx{
		Stream:  "task_request",
		Group:   gid(),
		Handler: cs.Ts.Request,
	})

	cs.Cli.Receive(rs.Rctx{
		Stream:  "task_stop",
		Group:   gid(),
		Handler: cs.Ts.Stop,
	})
}
