package mq

import (
	"github.com/Waitfantasy/tmq/config"
	"github.com/Waitfantasy/tmq/message"
	"github.com/Waitfantasy/tmq/util"
	"github.com/go-redis/redis"
)

type RedisMq struct {
	cli *redis.Client
}

func NewRedisMq(c *config.Config) *RedisMq {
	return &RedisMq{
		cli: util.NewRedisClient(c),
	}
}

func (mq *RedisMq) Send(topic string, msg *message.Message) error {
	data, err := msg.JsonMarshal()
	if err != nil {
		return err
	}

	cmd := mq.cli.Publish(topic, string(data))
	if _, err = cmd.Result(); err != nil {
		return err
	}

	return nil
}
