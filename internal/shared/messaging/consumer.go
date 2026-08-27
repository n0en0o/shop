package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultMaxDeliveryAttempts = 5
	defaultInitialRetryDelay   = 5 * time.Second
	defaultMaxRetryDelay       = 30 * time.Second

	retryCountHeader = "x-retry-count"
	errorHeader      = "x-last-error"
	failedAtHeader   = "x-failed-at"
	deadAtHeader     = "x-dead-at"
)

type MessageHandler func(ctx context.Context, body []byte) error

type Consumer interface {
	Consume(ctx context.Context, queueName string, handler MessageHandler) error
	Close() error
}

type PermanentMessageError struct {
	err error
}

func NewPermanentMessageError(err error) error {
	if err == nil {
		return nil
	}

	return PermanentMessageError{err: err}
}

func (e PermanentMessageError) Error() string {
	return e.err.Error()
}

func (e PermanentMessageError) Unwrap() error {
	return e.err
}

func IsPermanentMessageError(err error) bool {
	var permanentErr PermanentMessageError
	return errors.As(err, &permanentErr)
}

type ConsumerOption func(*RabbitMQConsumer)

func WithMaxDeliveryAttempts(attempts int) ConsumerOption {
	return func(c *RabbitMQConsumer) {
		if attempts > 0 {
			c.maxDeliveryAttempts = attempts
		}
	}
}

func WithRetryDelays(initialDelay, maxDelay time.Duration) ConsumerOption {
	return func(c *RabbitMQConsumer) {
		if initialDelay > 0 {
			c.initialRetryDelay = initialDelay
		}
		if maxDelay > 0 {
			c.maxRetryDelay = maxDelay
		}
	}
}

type RabbitMQConsumer struct {
	conn                *amqp.Connection
	channel             *amqp.Channel
	queues              map[string]consumerQueue
	maxDeliveryAttempts int
	initialRetryDelay   time.Duration
	maxRetryDelay       time.Duration
}

type consumerQueue struct {
	exchangeName    string
	queueName       string
	routingKey      string
	retryExchange   string
	retryQueue      string
	retryRoutingKey string
	deadExchange    string
	deadQueue       string
	deadRoutingKey  string
}

func NewRabbitMQConsumer(config RabbitMQConfig, opts ...ConsumerOption) (*RabbitMQConsumer, error) {
	var conn *amqp.Connection
	var err error

	for attempt := 1; attempt <= 10; attempt++ {
		log.Printf("[rabbitmq-consumer] connection attempt %d/10...", attempt)
		conn, err = amqp.Dial(config.ConnectionString())
		if err == nil {
			log.Println("[rabbitmq-consumer] connected")
			break
		}

		log.Printf("[rabbitmq-consumer] connection error: %v", err)
		if attempt < 10 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("create RabbitMQ channel: %w", err)
	}

	if err := channel.Qos(1, 0, false); err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("set RabbitMQ QoS: %w", err)
	}

	consumer := &RabbitMQConsumer{
		conn:                conn,
		channel:             channel,
		queues:              make(map[string]consumerQueue),
		maxDeliveryAttempts: defaultMaxDeliveryAttempts,
		initialRetryDelay:   defaultInitialRetryDelay,
		maxRetryDelay:       defaultMaxRetryDelay,
	}

	for _, opt := range opts {
		opt(consumer)
	}

	return consumer, nil
}

func (c *RabbitMQConsumer) SetupQueue(
	exchangeName, exchangeType, queueName, routingKey string,
) error {
	queue := consumerQueue{
		exchangeName:    exchangeName,
		queueName:       queueName,
		routingKey:      routingKey,
		retryExchange:   exchangeName + ".retry",
		retryQueue:      queueName + ".retry",
		retryRoutingKey: routingKey + ".retry",
		deadExchange:    exchangeName + ".dead",
		deadQueue:       queueName + ".dead",
		deadRoutingKey:  routingKey + ".dead",
	}

	err := c.channel.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}

	_, err = c.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	err = c.channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}

	if err := c.setupRetryQueue(queue, exchangeType); err != nil {
		return err
	}

	if err := c.setupDeadQueue(queue, exchangeType); err != nil {
		return err
	}

	c.queues[queueName] = queue

	log.Printf(
		"[rabbitmq-consumer] queue '%s' configured on exchange '%s' with retry queue '%s' and dead queue '%s'",
		queueName,
		exchangeName,
		queue.retryQueue,
		queue.deadQueue,
	)
	return nil
}

func (c *RabbitMQConsumer) setupRetryQueue(queue consumerQueue, exchangeType string) error {
	err := c.channel.ExchangeDeclare(queue.retryExchange, exchangeType, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare retry exchange: %w", err)
	}

	_, err = c.channel.QueueDeclare(
		queue.retryQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    queue.exchangeName,
			"x-dead-letter-routing-key": queue.routingKey,
		},
	)
	if err != nil {
		return fmt.Errorf("declare retry queue: %w", err)
	}

	err = c.channel.QueueBind(queue.retryQueue, queue.retryRoutingKey, queue.retryExchange, false, nil)
	if err != nil {
		return fmt.Errorf("bind retry queue: %w", err)
	}

	return nil
}

func (c *RabbitMQConsumer) setupDeadQueue(queue consumerQueue, exchangeType string) error {
	err := c.channel.ExchangeDeclare(queue.deadExchange, exchangeType, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare dead exchange: %w", err)
	}

	_, err = c.channel.QueueDeclare(queue.deadQueue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare dead queue: %w", err)
	}

	err = c.channel.QueueBind(queue.deadQueue, queue.deadRoutingKey, queue.deadExchange, false, nil)
	if err != nil {
		return fmt.Errorf("bind dead queue: %w", err)
	}

	return nil
}

func (c *RabbitMQConsumer) Consume(
	ctx context.Context,
	queueName string,
	handler MessageHandler,
) error {
	queue, ok := c.queues[queueName]
	if !ok {
		return fmt.Errorf("queue %q is not configured; call SetupQueue first", queueName)
	}

	msgs, err := c.channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("start consuming queue %q: %w", queueName, err)
	}

	log.Printf("[rabbitmq-consumer] consuming from queue '%s'", queueName)

	for {
		select {
		case <-ctx.Done():
			log.Println("[rabbitmq-consumer] stop signal received")
			return ctx.Err()

		case msg, ok := <-msgs:
			if !ok {
				log.Println("[rabbitmq-consumer] delivery channel closed")
				return nil
			}

			log.Printf("[rabbitmq-consumer] received message: %d bytes", len(msg.Body))

			if err := handler(ctx, msg.Body); err != nil {
				log.Printf("[rabbitmq-consumer] processing error: %v", err)
				if ctx.Err() != nil {
					if nackErr := msg.Nack(false, true); nackErr != nil {
						return fmt.Errorf("nack canceled message: %w", nackErr)
					}

					return ctx.Err()
				}

				if err := c.handleFailedMessage(ctx, queue, msg, err); err != nil {
					return err
				}
			} else {
				log.Println("[rabbitmq-consumer] message processed successfully")
				if err := msg.Ack(false); err != nil {
					return fmt.Errorf("ack message: %w", err)
				}
			}
		}
	}
}

func (c *RabbitMQConsumer) handleFailedMessage(
	ctx context.Context,
	queue consumerQueue,
	msg amqp.Delivery,
	processingErr error,
) error {
	retryCount := readRetryCount(msg.Headers)
	failedAttempt := retryCount + 1

	if IsPermanentMessageError(processingErr) || failedAttempt >= c.maxDeliveryAttempts {
		if err := c.publishDead(ctx, queue, msg, processingErr, failedAttempt); err != nil {
			if nackErr := msg.Nack(false, true); nackErr != nil {
				return fmt.Errorf("dead publish failed: %w; nack failed: %v", err, nackErr)
			}

			return fmt.Errorf("dead publish failed: %w", err)
		}

		if err := msg.Ack(false); err != nil {
			return fmt.Errorf("ack dead-lettered message: %w", err)
		}

		log.Printf(
			"[rabbitmq-consumer] message moved to dead queue '%s' after attempt %d/%d",
			queue.deadQueue,
			failedAttempt,
			c.maxDeliveryAttempts,
		)
		return nil
	}

	delay := c.retryDelay(failedAttempt)
	if err := c.publishRetry(ctx, queue, msg, processingErr, failedAttempt, delay); err != nil {
		if nackErr := msg.Nack(false, true); nackErr != nil {
			return fmt.Errorf("retry publish failed: %w; nack failed: %v", err, nackErr)
		}

		return fmt.Errorf("retry publish failed: %w", err)
	}

	if err := msg.Ack(false); err != nil {
		return fmt.Errorf("ack retried message: %w", err)
	}

	log.Printf(
		"[rabbitmq-consumer] message scheduled for retry %d/%d after %s",
		failedAttempt,
		c.maxDeliveryAttempts,
		delay,
	)
	return nil
}

func (c *RabbitMQConsumer) publishRetry(
	ctx context.Context,
	queue consumerQueue,
	msg amqp.Delivery,
	processingErr error,
	failedAttempt int,
	delay time.Duration,
) error {
	publishing := clonePublishing(msg)
	publishing.Headers[retryCountHeader] = int32(failedAttempt)
	publishing.Headers[errorHeader] = processingErr.Error()
	publishing.Headers[failedAtHeader] = time.Now().UTC().Format(time.RFC3339Nano)
	publishing.Expiration = strconv.FormatInt(delay.Milliseconds(), 10)

	return c.channel.PublishWithContext(
		ctx,
		queue.retryExchange,
		queue.retryRoutingKey,
		false,
		false,
		publishing,
	)
}

func (c *RabbitMQConsumer) publishDead(
	ctx context.Context,
	queue consumerQueue,
	msg amqp.Delivery,
	processingErr error,
	failedAttempt int,
) error {
	publishing := clonePublishing(msg)
	publishing.Headers[retryCountHeader] = int32(failedAttempt)
	publishing.Headers[errorHeader] = processingErr.Error()
	publishing.Headers[deadAtHeader] = time.Now().UTC().Format(time.RFC3339Nano)

	return c.channel.PublishWithContext(
		ctx,
		queue.deadExchange,
		queue.deadRoutingKey,
		false,
		false,
		publishing,
	)
}

func clonePublishing(msg amqp.Delivery) amqp.Publishing {
	headers := make(amqp.Table, len(msg.Headers)+4)
	for key, value := range msg.Headers {
		headers[key] = value
	}

	return amqp.Publishing{
		Headers:         headers,
		ContentType:     msg.ContentType,
		ContentEncoding: msg.ContentEncoding,
		DeliveryMode:    amqp.Persistent,
		Priority:        msg.Priority,
		CorrelationId:   msg.CorrelationId,
		ReplyTo:         msg.ReplyTo,
		Expiration:      msg.Expiration,
		MessageId:       msg.MessageId,
		Timestamp:       time.Now().UTC(),
		Type:            msg.Type,
		UserId:          msg.UserId,
		AppId:           msg.AppId,
		Body:            msg.Body,
	}
}

func readRetryCount(headers amqp.Table) int {
	if headers == nil {
		return 0
	}

	switch value := headers[retryCountHeader].(type) {
	case int:
		return value
	case int8:
		return int(value)
	case int16:
		return int(value)
	case int32:
		return int(value)
	case int64:
		return int(value)
	case uint:
		return int(value)
	case uint8:
		return int(value)
	case uint16:
		return int(value)
	case uint32:
		return int(value)
	case uint64:
		return int(value)
	case string:
		retryCount, err := strconv.Atoi(value)
		if err == nil {
			return retryCount
		}
	}

	return 0
}

func (c *RabbitMQConsumer) retryDelay(failedAttempt int) time.Duration {
	delay := c.initialRetryDelay
	for i := 1; i < failedAttempt; i++ {
		delay *= 2
		if delay >= c.maxRetryDelay {
			return c.maxRetryDelay
		}
	}

	if delay > c.maxRetryDelay {
		return c.maxRetryDelay
	}

	return delay
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
		return nil, NewPermanentMessageError(
			fmt.Errorf("deserialize message: %w", err),
		)
	}

	return &msg, nil
}
