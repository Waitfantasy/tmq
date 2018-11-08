package manager

import (
	"errors"
	"github.com/Waitfantasy/tmq/message"
	"github.com/Waitfantasy/tmq/message/persistent"
	"github.com/Waitfantasy/tmq/mq"
	"github.com/Waitfantasy/unicorn/rpc/client"
)

type Manager struct {
	persistent persistent.Persistenter
	cli        *client.Client
	mq         mq.Mqer
}

func New(mqer mq.Mqer, persistenter persistent.Persistenter, idRpcCli *client.Client) *Manager {
	return &Manager{
		mq:         mqer,
		persistent: persistenter,
		cli:        idRpcCli,
	}
}

func (manager *Manager) Prepare(topic string, retrySecond int, body string) (*message.Message, error) {
	response, err := manager.cli.MakeUUID()
	if err != nil {
		return nil, err
	}

	// 创建Prepare消息
	msg := message.NewPrepareMessage(response.Uuid, topic, retrySecond, body)
	if err = manager.persistent.Store(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (manager *Manager) Send(id uint64) (*message.Message, error) {
	msg, err := manager.persistent.Find(id)
	if err != nil {
		return nil, err
	}

	if !msg.IsPrepare() {
		return nil, errors.New("the message cannot be commit")
	}

	if err := manager.mq.Send(msg.Topic, msg); err != nil {
		return nil, err
	}

	// redis pub/sub需要等待业务被动方回传ack后才可以确认
	if _, ok := manager.mq.(*mq.RedisMq); ok {
		return msg, nil
	}

	if !msg.IsCommit() {
		return nil, errors.New("the message not commit")
	}

	// 将提交状态的消息持久化
	if err = manager.persistent.Update(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (manager *Manager) ConsumerCommit(id uint64) (*message.Message, error) {
	// 查找消息
	msg, err := manager.persistent.Find(id)
	if err != nil {
		return nil, err
	}

	// 提交消息
	if !msg.IsPrepare() {
		return nil, errors.New("the message not commit")
	}

	msg.Commit()

	// 将提交状态的消息持久化
	if err = manager.persistent.Update(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (manager *Manager) Cancel(id uint64) (*message.Message, error) {
	// 查找消息
	msg, err := manager.persistent.Find(id)
	if err != nil {
		return nil, err
	}

	if !msg.IsCancelOk() {
		return nil, errors.New("the message can not rollback")
	}

	// delete msg
	//manager.persistent.Destory(id)
	return msg, nil
}
