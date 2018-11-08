package mq

import (
	"context"
	"github.com/Waitfantasy/tmq/message"
	"github.com/nsqio/go-nsq"
	"time"
)

type NsqMq struct {
	addr  string
	async bool
	c     *nsq.Config
	delay time.Duration
}

func NewNsqMq(addr string, delay time.Duration) *NsqMq {
	return &NsqMq{
		addr:  addr,
		async: true,
		c:     nsq.NewConfig(),
		delay: delay,
	}
}

func (mq *NsqMq) Send(topic string, msg *message.Message) error {
	err := mq.publish(topic, msg)
	return err
	//if mq.async {
	//	err := mq.publishAsync(topic, msg, func(ctx context.Context, msg *message.Message, c chan *nsq.ProducerTransaction) {
	//		select {
	//		case <-ctx.Done():
	//			// TODO write log
	//			fmt.Println("msg publish fail: ", msg)
	//		case <-c:
	//			msg.Commit()
	//			fmt.Println("msg commit")
	//		}
	//	})
	//
	//	return err
	//} else {
	//	err := mq.publish(topic, msg)
	//	return err
	//}
}

func (mq *NsqMq) SetAsync(ok bool) {
	mq.async = ok
}

func (mq *NsqMq) publish(topic string, msg *message.Message) error {
	producer, err := nsq.NewProducer(mq.addr, mq.c)
	if err != nil {
		return err
	}
	defer producer.Stop()
	data, err := msg.JsonMarshal()
	if err != nil {
		return err
	}

	if err := producer.Publish(topic, data); err != nil {
		return err
	}

	msg.Commit()
	return nil
}

func (mq *NsqMq) publishAsync(topic string, msg *message.Message, callback func(ctx context.Context,
	msg *message.Message, c chan *nsq.ProducerTransaction)) error {
	producer, err := nsq.NewProducer(mq.addr, mq.c)
	if err != nil {
		return err
	}

	data, err := msg.JsonMarshal()
	if err != nil {
		return err
	}

	responseChan := make(chan *nsq.ProducerTransaction)
	if err := producer.PublishAsync(topic, data, responseChan); err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), mq.delay)

	go callback(ctx, msg, responseChan)

	return nil
}
