package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (c RabbitMQConfig) ConnectionString() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
	)
}

type Publisher interface {
	Publish(
		ctx context.Context,
		event interface{},
		exchangeName, routingKey string,
	) error
	Close() error
}

type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQPublisher(config RabbitMQConfig) (*RabbitMQPublisher, error) {
	var conn *amqp.Connection
	var err error

	for attempt := 1; attempt <= 10; attempt++ {
		log.Printf("[rabbitmq] попытка подключения %d/10...", attempt)
		conn, err = amqp.Dial(config.ConnectionString())
		if err == nil {
			log.Println("[rabbitmq] подключение успешно")
			break
		}

		log.Printf("[rabbitmq] ошибка: %v", err)
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

	return &RabbitMQPublisher{
		conn:    conn,
		channel: channel,
	}, nil
}

func (p *RabbitMQPublisher) DeclareExchange(name, kind string) error {
	return p.channel.ExchangeDeclare(name, kind, true, false, false, false, nil)
}

func (p *RabbitMQPublisher) DeclareQueue(
	queueName, exchangeName, routingKey string,
) error {
	_, err := p.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("не удалось объявить очередь: %w", err)
	}

	err = p.channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		return fmt.Errorf("не удалось привязать очередь к exchange: %w", err)
	}

	return nil
}

func (p *RabbitMQPublisher) Publish(
	ctx context.Context,
	event interface{},
	exchangeName, routingKey string,
) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("не удалось сериализовать событие: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("не удалось опубликовать событие: %w", err)
	}

	log.Printf("[rabbitmq] событие опубликовано в %s/%s", exchangeName, routingKey)
	return nil
}

func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}

	return nil
}

func (p *RabbitMQPublisher) SetupOrderEventsExchange() error {
	err := p.DeclareExchange(OrderSubmittedExchange, "direct")
	if err != nil {
		return err
	}

	if err := p.DeclareQueue(
		OrderSubmittedQueue,
		OrderSubmittedExchange,
		OrderSubmittedEventType,
	); err != nil {
		return err
	}

	log.Printf(
		"[rabbitmq] настроен exchange '%s' и queue '%s'",
		OrderSubmittedExchange,
		OrderSubmittedQueue,
	)
	return nil
}

func (p *RabbitMQPublisher) PublishOrderSubmitted(
	ctx context.Context,
	event *OrderSubmittedEvent,
) error {
	return p.Publish(ctx, event, OrderSubmittedExchange, OrderSubmittedEventType)
}
