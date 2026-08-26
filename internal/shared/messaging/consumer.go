package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageHandler func(ctx context.Context, body []byte) error

type Consumer interface {
	Consume(ctx context.Context, queueName string, handler MessageHandler) error
	Close() error
}

type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQConsumer(config RabbitMQConfig) (*RabbitMQConsumer, error) {
	var conn *amqp.Connection
	var err error

	for attempt := 1; attempt <= 10; attempt++ {
		log.Printf("[rabbitmq-consumer] попытка подключения %d/10...", attempt)
		conn, err = amqp.Dial(config.ConnectionString())
		if err == nil {
			log.Println("[rabbitmq-consumer] подключение успешно")
			break
		}
		log.Printf("[rabbitmq-consumer] ошибка: %v", err)
		if attempt < 10 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("не удалось создать канал: %w", err)
	}

	if err := channel.Qos(1, 0, false); err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("не удалось установить QoS: %w", err)
	}

	return &RabbitMQConsumer{
		conn:    conn,
		channel: channel,
	}, nil
}

func (c *RabbitMQConsumer) SetupQueue(
	exchangeName, exchangeType, queueName, routingKey string,
) error {
	err := c.channel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("не удалось объявить exchange: %w", err)
	}

	_, err = c.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("не удалось объявить очередь: %w", err)
	}

	err = c.channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		return fmt.Errorf("не удалось привязать очередь: %w", err)
	}

	log.Printf(
		"[rabbitmq-consumer] настроена очередь '%s' на exchange '%s'",
		queueName,
		exchangeName,
	)
	return nil
}

func (c *RabbitMQConsumer) Consume(
	ctx context.Context,
	queueName string,
	handler MessageHandler,
) error {
	msgs, err := c.channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("не удалось начать потребление: %w", err)
	}

	log.Printf("[rabbitmq-consumer] начато потребление из очереди '%s'", queueName)

	for {
		select {
		case <-ctx.Done():
			log.Println("[rabbitmq-consumer] получен сигнал остановки")
			return ctx.Err()

		case msg, ok := <-msgs:
			if !ok {
				log.Println("[rabbitmq-consumer] канал сообщений закрыт")
				return nil
			}

			log.Printf("[rabbitmq-consumer] получено сообщение: %d bytes", len(msg.Body))

			if err := handler(ctx, msg.Body); err != nil {
				log.Printf("[rabbitmq-consumer] ошибка обработки: %v", err)
				msg.Nack(false, true)
			} else {
				log.Println("[rabbitmq-consumer] сообщение обработано успешно")
				msg.Ack(false)
			}
		}
	}
}

func (c *RabbitMQConsumer) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

func ParseMessage[T any](body []byte) (*T, error) {
	var msg T
	if err := json.Unmarshal(body, &msg); err != nil {
		return nil, fmt.Errorf("не удалось десериализовать сообщение: %w", err)
	}

	return &msg, nil
}
