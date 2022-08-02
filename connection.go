package rabbitmq

import (
	"fmt"
	"rabbitmq/amqp091"
)

type connChannel struct {
	amqp091.Channel
	conn *connection
}

// 之前测试过判断关闭的话需要两者配合来判断
func (c *connChannel) IsClosed() bool {
	if c.Channel.IsClosed() || c.conn.IsClosed() {
		return true
	}
	return false
}

type connection struct {
	amqp091.Connection                     // 原生句柄
	closeChan          chan *amqp091.Error // 监听关闭的chan
}

// Channel需要控制在2K以下，封装此函数主要是考虑将来可能要做p2c之类的逻辑，目前来说没有什么必须封装一层的必要性
func (c *connection) Channel() (amqp091.Channel, error) {
	ch, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}

	return &connChannel{
		Channel: ch,
		conn:    c,
	}, nil
}

// 这个函数的封装只是为了将异常关闭的信息打印出来，closeChan并不是必须的
func (c *connection) IsClosed() bool {
	if c.Connection.IsClosed() {
		return true
	}

	select {
	case err, ok := <-c.closeChan:
		if ok {
			fmt.Printf("rabbitmq connection closed: %v\n", err)
		}
		return true
	default:
	}
	return false
}

func newConnection(conf amqp091.Config) (*connection, error) {
	// 创建连接
	conn, err := amqp091.New(conf)
	if err != nil {
		return nil, err
	}

	// 注册关闭函数
	closeChan := make(chan *amqp091.Error)
	conn.NotifyClose(closeChan)

	return &connection{
		Connection: conn,
		closeChan:  closeChan,
	}, nil
}
