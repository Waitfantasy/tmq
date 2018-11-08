package mq

import "github.com/Waitfantasy/tmq/message"

const (
	NsqType   = "nsq"
	RedisType = "redis"
)

type Mqer interface {
	Send(topic string, msg *message.Message) error
}
