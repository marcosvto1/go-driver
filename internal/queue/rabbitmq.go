package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	URL       string
	TopicName string
	Timeout   time.Time
}

type RabbitConnection struct {
	cfg  RabbitMQConfig
	conn *amqp091.Connection
}

func newRabbitConn(cfg RabbitMQConfig) (rc *RabbitConnection, err error) {
	rc = new(RabbitConnection)
	rc.cfg = cfg
	rc.conn, err = amqp091.Dial(cfg.URL)
	return rc, err
}

func (rc *RabbitConnection) Publish(msg []byte) error {
	log.Println("Publish WITH Rabbit")

	c, err := rc.conn.Channel()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	mp := amqp091.Publishing{
		ContentType: "text/plain",
		Timestamp:   rc.cfg.Timeout,
		Body:        msg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = c.PublishWithContext(ctx, "", rc.cfg.TopicName, false, false, mp)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (rc *RabbitConnection) Consume(cdto chan<- QueueDTO) error {
	log.Println("Consume WITH Rabbit")

	c, err := rc.conn.Channel()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	defer c.Close()

	q, err := c.QueueDeclare(rc.cfg.TopicName, true, false, false, false, nil)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	msgs, err := c.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for msg := range msgs {
		log.Printf("Received message")

		dto := QueueDTO{}
		dto.Unmarshal(msg.Body)

		cdto <- dto
	}

	return nil
}
