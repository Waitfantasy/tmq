package mq

import (
	"fmt"
	"github.com/Waitfantasy/tmq/message"
	"github.com/nsqio/go-nsq"
	"strconv"
	"testing"
	"time"
)

type ConsumerHandler struct {
	t              *testing.T
	q              *nsq.Consumer
	messagesGood   int
	messagesFailed int
}

func (h *ConsumerHandler) LogFailedMessage(message *nsq.Message) {
	h.messagesFailed++
	h.q.Stop()
}

func (h *ConsumerHandler) HandleMessage(message *nsq.Message) error {
	msg := string(message.Body)
	fmt.Println("receive msg: ", msg)
	h.messagesGood++
	h.q.Stop()
	return nil
}

func TestNsqMq_Send(t *testing.T) {
	topicName := "publish_" + strconv.Itoa(int(time.Now().Unix()))
	mq := NewNsqMq("127.0.0.1:4150", time.Second*10)
	msg := message.NewPrepareMessage(1,"test_sync_send_case", "")
	if err := mq.Send(topicName, msg); err != nil {
		t.Error(err)
		return
	}

	if msg.Status != message.Commit {
		t.Error("message commit fail.")
	}

	readMessages(topicName, t, 1)
}

func readMessages(topicName string, t *testing.T, msgCount int) {
	config := nsq.NewConfig()
	config.DefaultRequeueDelay = 0
	config.MaxBackoffDuration = 50 * time.Millisecond
	q, _ := nsq.NewConsumer(topicName, "ch", config)

	h := &ConsumerHandler{
		t: t,
		q: q,
	}
	q.AddHandler(h)

	err := q.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		t.Fatalf(err.Error())
	}
	<-q.StopChan

	if h.messagesGood != msgCount {
		t.Fatalf("end of test. should have handled a diff number of messages %d != %d", h.messagesGood, msgCount)
	}

	if h.messagesFailed != 0 {
		t.Fatal("failed message not done")
	}
}
