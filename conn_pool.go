package rabbitmq

import (
	"rabbitmq/amqp091"
	"sync"
)

// 通过connPool来维护Rabbitmq的连接，目前实现的是单连接RabbitmqConnection池
type connPool struct {
	lock sync.Mutex
	conn *connection

	conf amqp091.Config
}

// 创建一个连接
func (c *connPool) Connection() (amqp091.Connection, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// 优先取现有的连接
	if c.conn != nil && !c.conn.IsClosed() {
		return c.conn, nil
	}
	c.conn = nil

	// 尝试获取新的Connection
	conn, err := newConnection(c.conf)
	if err != nil {
		return nil, err
	}

	c.conn = conn
	return conn, nil
}

func (c *connPool) Close() error {
	return c.conn.Close()
}

func newConnPool(conf amqp091.Config) connPool {
	return connPool{
		conf: conf,
	}
}
