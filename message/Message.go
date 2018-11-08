package message

import (
	"encoding/json"
	"github.com/vmihailenco/msgpack"
	"time"
)

const (
	Prepare int = iota
	Commit
	Cancel
)

var status = [3]string{"prepare", "commit", "cancel"}

type Message struct {
	Id              uint64 `json:"id"`
	Status          int    `json:"status"`
	Topic           string `json:"topic"`
	Retries         int    `json:"retries"`
	RetrySecond     int    `json:"retry_second"`
	CreateTimestamp int64  `json:"create_timestamp"`
	UpdateTimestamp int64  `json:"update_timestamp"`
	Body            string `json:"body"`
}

func NewPrepareMessage(id uint64, topic string, retrySecond int, body string) *Message {
	message := &Message{}
	message.Id = id
	message.Status = Prepare
	message.Topic = topic
	message.RetrySecond = retrySecond
	message.CreateTimestamp = time.Now().Unix()
	message.UpdateTimestamp = message.CreateTimestamp
	message.Body = body
	return message
}

func (m *Message) MsgPackMarshal() ([]byte, error)  {
	return msgpack.Marshal(m)
}

func (m *Message) MsgPackUnmarshal(data []byte) error {
	return msgpack.Unmarshal(data, m)
}

func (m *Message) JsonMarshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Message) Commit() {
	m.Status = Commit
	m.UpdateTimestamp = time.Now().Unix()
}

func (m *Message) IsPrepare() bool {
	return m.Status == Prepare
}

func (m *Message) IsCommit() bool {
	return m.Status == Commit
}

func (m *Message) IsCancelOk() bool {
	return m.Status == Prepare
}

func (m *Message) StatusToString() string {
	return status[m.Status]
}

func (m *Message) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"status":           m.Status,
		"topic":            m.Topic,
		"retries":          m.Retries,
		"retry_second":     m.RetrySecond,
		"create_timestamp": m.CreateTimestamp,
		"update_timestamp": m.UpdateTimestamp,
		"body":             m.Body,
	}
}
