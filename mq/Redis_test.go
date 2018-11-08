package mq

import (
	"fmt"
	"github.com/Waitfantasy/tmq/config"
	"github.com/Waitfantasy/tmq/message"
	"github.com/go-redis/redis"
	"strconv"
	"testing"
	"time"
)

func TestMq_Send(t *testing.T) {
	mq := NewNsqMq(&config.Config{
		Redis: struct {
			Addr     string
			Password string
			DB       int
		}{Addr: "127.0.0.1:6379", Password: "", DB: 0,},
	})
	topicName := "redis_publish_" + strconv.Itoa(int(time.Now().Unix()))
	stop := make(chan bool)
	go readMessage(topicName, stop)
	msg := message.NewPrepareMessage("redis_test_publish_case", "{foo: \"bar\"}")
	if err := mq.Send(topicName, msg); err != nil {
		t.Error(err)
	}
	<-stop
}

func readMessage(topicName string, c chan bool) {
	cli := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	pubSub := cli.Subscribe(topicName)
	msg := <-pubSub.Channel()
	fmt.Println("receive: ", msg.Payload)
	c <- true
}
