package discover

import (
	"github.com/Waitfantasy/tmq/message/persistent"
	"github.com/Waitfantasy/tmq/mq"
)

type Discover struct {
	mq mq.Mqer
	p  persistent.Persistenter
}

func New(p persistent.Persistenter) *Discover {
	return &Discover{
		p: p,
	}
}
