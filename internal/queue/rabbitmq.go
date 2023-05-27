package queue

import (
	"context"
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
	rc.cfg = cfg
	rc.conn, err = amqp091.Dial(cfg.URL)
	return rc, err
}

func (rc *RabbitConnection) Publish(msg []byte) error {
	c, err := rc.conn.Channel()
	if err != nil {
		return err
	}

	mp := amqp091.Publishing{
		ContentType: "text/plain",
		Timestamp:   rc.cfg.Timeout,
		Body:        msg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return c.PublishWithContext(ctx, "", rc.cfg.TopicName, false, false, mp)
}

func (rc *RabbitConnection) Consume(cdto chan<- QueueDTO) error {
	c, err := rc.conn.Channel()
	if err != nil {
		return err
	}

	q, err := c.QueueDeclare(rc.cfg.TopicName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := c.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range msgs {
		dto := QueueDTO{}
		dto.Unmarshal(msg.Body)
		cdto <- dto
	}

	return nil
}
