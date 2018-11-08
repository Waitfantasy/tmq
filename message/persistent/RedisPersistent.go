package persistent

import (
	"github.com/Waitfantasy/tmq/config"
	"github.com/Waitfantasy/tmq/message"
	"github.com/Waitfantasy/tmq/util"
	"github.com/go-redis/redis"
	"strconv"
)

type RedisPersistent struct {
	cli *redis.Client
}

func NewRedisPersistent(c *config.Config) (*RedisPersistent) {
	return &RedisPersistent{
		cli: util.NewRedisClient(c),
	}
}

func (p *RedisPersistent) Store(msg *message.Message) error {
	key := key(msg.Id, msg.StatusToString())
	data, err := msg.MsgPackMarshal()
	if err !=nil {
		return err
	}

	// 存储消息
	if status := p.cli.Set(key, data, 0); status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (p *RedisPersistent) Update(msg *message.Message) error {
	return p.Store(msg)
}

func (p *RedisPersistent) Find(id uint64) (*message.Message, error) {
	key := key(id, "prepare")
	cmd := p.cli.Get(key)
	if result, err := cmd.Bytes(); err != nil {
		return nil, err

	} else {
		msg := new(message.Message)
		if err = msg.MsgPackUnmarshal(result); err != nil {
			return nil, err
		}
		return msg, nil
	}
}

func key(id uint64, status string) string {
	return "tmp:" + status + ":" + strconv.FormatUint(id, 10)
}
