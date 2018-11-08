package persistent

import "github.com/Waitfantasy/tmq/message"

const (
	RedisPersistentType = "redis"
	MysqlPersistentType = "mysql"
	MongoPersistentType = "mongodb"
)

type Persistenter interface {
	Store(msg *message.Message) error
	Update(msg *message.Message) error
	Find(id uint64) (*message.Message, error)
}
