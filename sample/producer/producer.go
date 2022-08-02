package main

import (
	"log"
	"rabbitmq"
	"rabbitmq/amqp091"
)

func main() {
	p, err := rabbitmq.New(amqp091.Config{
		UserName: "guest",
		Password: "guest",
		Host:     "127.0.0.1:5672",
	})
	if err != nil {
		log.Fatal(err)
	}

	for j := 1; j < 30; j++ {
		go func() {
			for i := 0; i < 100; i++ {
				err := p.Publish("exchange", "key", false, amqp091.Publishing{
					Body: []byte("xxx"),
				})
				if err != nil {
					log.Fatal(err)
				}
			}
		}()

	}

	for i := 0; i < 30000; i++ {
		err := p.Publish("exchange", "key", false, amqp091.Publishing{
			Body: []byte("xxx"),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	err = p.Close()
	if err != nil {
		log.Fatal(err)
	}
}
