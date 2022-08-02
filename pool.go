package rabbitmq

import (
	"rabbitmq/amqp091"
	"sync/atomic"
)

// 对Channel进行封装，因为pool一层有特殊处理
type poolChannel struct {
	amqp091.Channel
	p *pool
}

// 封装处理，优先放回到pool中
func (c *poolChannel) Close() error {
	// 如果Channel已经无效或者池子已经关闭就不必放回去了
	if c.IsClosed() || atomic.LoadUint32(&c.p.close) == 1 {
		return c.Channel.Close()
	}
	// 放回到pool中
	c.p.returnChannel(c)
	return nil
}

type pool struct {
	connPool connPool    // 连接池
	chl      channelList // Channel池，其实就是一个链表
	close    uint32      // 池子是否已关闭
}

func (p *pool) returnChannel(c *poolChannel) {
	p.chl.PushBack(c)
}

func (p *pool) Channel() (amqp091.Channel, error) {
	for {
		ch, ok := p.chl.PopFront()
		// 链表中已经没有数据了
		if !ok {
			break
		}

		// 如果有效的话直接返回
		if !ch.IsClosed() {
			return ch, nil
		}

		// 关闭后不必再放回
		_ = ch.Close()
	}

	conn, err := p.connPool.Connection()
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &poolChannel{
		Channel: ch,
		p:       p,
	}, nil
}

// 推送消息只是保证取出来的Channel是有效的
func (p *pool) Publish(exchange, key string, mandatory bool, msg amqp091.Publishing) error {
	ch, err := p.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.Publish(exchange, key, mandatory, false, msg)
}

func (p *pool) Close() error {
	if !atomic.CompareAndSwapUint32(&p.close, 0, 1) {
		return nil
	}

	for {
		ch, ok := p.chl.PopFront()
		if !ok {
			break
		}
		_ = ch.Close()
	}

	return p.connPool.Close()
}

func New(conf amqp091.Config) (Pool, error) {
	return &pool{
		chl:      newChannelList(),
		connPool: newConnPool(conf),
	}, nil
}
