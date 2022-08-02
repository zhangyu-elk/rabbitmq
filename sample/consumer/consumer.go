package main

import (
	"fmt"
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

	ch, err := p.Channel()
	if err != nil {
		log.Fatal(err)
	}

	_, err = ch.QueueDeclare("queue", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	dch, err := ch.Consume("queue", "consumer", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for {
		v, ok := <-dch
		if !ok {
			log.Fatal("closed")
		}

		count++
		fmt.Println("recv: ", string(v.Body), "count: ", count)
	}
}
