package amqp091

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
)

type Config struct {
	UserName string // 用户名
	Password string // 密码
	Host     string // 地址
	VHost    string //Vhost
}

// 创建一个新的连接
func New(conf Config) (Connection, error) {
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s/%s", conf.UserName, conf.Password, conf.Host, conf.VHost)

	conn, err := amqp091.DialConfig(amqpURL, amqp091.Config{})
	if err != nil {
		return nil, err
	}

	return &connection{conn}, nil
}
