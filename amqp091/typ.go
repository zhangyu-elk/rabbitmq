package amqp091

import "github.com/rabbitmq/amqp091-go"

// 原生接口定义，这里主要是封装出Interface进行单测

type (
	Publishing = amqp091.Publishing
	Table      = amqp091.Table
	Delivery   = amqp091.Delivery
	Error      = amqp091.Error
	Queue      = amqp091.Queue
)

type Channel interface {
	Publish(exchange, key string, mandatory, immediate bool, msg Publishing) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table) (<-chan amqp091.Delivery, error)

	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table) (Queue, error)
	NotifyClose(c chan *Error) chan *Error

	IsClosed() bool
	Close() error
}

type Connection interface {
	Channel() (Channel, error)
	NotifyClose(receiver chan *Error) chan *Error
	IsClosed() bool
	Close() error
}

type connection struct {
	*amqp091.Connection
}

func (c *connection) Channel() (Channel, error) {
	return c.Connection.Channel()
}
