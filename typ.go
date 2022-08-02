package rabbitmq

import "rabbitmq/amqp091"

// 后续可以考虑对Consume之类的进行封装处理
type Pool interface {
	// 获取一个Channel，注意必须Close
	Channel() (amqp091.Channel, error)
	// 发布消息
	Publish(exchange, key string, mandatory bool, msg amqp091.Publishing) error // 发布消息
	// 关闭池子，目前此函数不可重入
	Close() error
}
