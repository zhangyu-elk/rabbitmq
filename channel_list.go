package rabbitmq

import (
	"container/list"
	"rabbitmq/amqp091"
	"sync"
)

type channelList struct {
	m sync.Mutex
	l *list.List
}

func (c *channelList) PushBack(ch amqp091.Channel) {
	c.m.Lock()
	defer c.m.Unlock()

	c.l.PushBack(ch)
}

func (c *channelList) PopFront() (amqp091.Channel, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.l.Len() == 0 {
		return nil, false
	}

	e := c.l.Front()
	c.l.Remove(e)
	return e.Value.(amqp091.Channel), true
}

func newChannelList() channelList {
	return channelList{
		l: list.New(),
	}
}
