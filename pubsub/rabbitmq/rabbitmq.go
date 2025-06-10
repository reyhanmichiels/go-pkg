package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/operator"
)

type ExchangeType string

// Exchange type
const (
	ExchangeTypeDirect  ExchangeType = "direct"
	ExchangeTypeFanout  ExchangeType = "fanout"
	ExchangeTypeTopic   ExchangeType = "topic"
	ExchangeTypeHeaders ExchangeType = "headers"
)

// Config holds RabbitMQ connection configuration
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	SSL      bool
	Options  ConnOptions
}

// ConnOptions holds additional configuration options
type ConnOptions struct {
	ReconnectDelay time.Duration
}

// MessageHandler is a function type for handling received messages
type MessageHandler func(ctx context.Context, exchangeName, routingKey, message string) error

// Interface defines the RabbitMQ wrapper methods
type Interface interface {
	Publish(ctx context.Context, exchangeName string, routingKey string, body string) error
	Subscribe(ctx context.Context, queueName string, handler MessageHandler)
	CreateExchange(name string, kind ExchangeType, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	CreateQueue(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error)
	BindQueue(queueName, exchangeName, routingKey string, noWait bool, args amqp.Table) error
	MonitorConnection()
	Stop()
}

type rabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	cfg        Config
	endOnce    *sync.Once
	log        log.Interface
}

// Init creates a new RabbitMQ client
func Init(cfg Config, log log.Interface) Interface {
	rmq := &rabbitMQ{
		cfg:     cfg,
		endOnce: &sync.Once{},
		log:     log,
	}

	rmq.init()
	return rmq
}

func (r *rabbitMQ) init() {
	var err error

	// Attempt to connect with retries
	retryTime := operator.Ternary(r.cfg.Options.ReconnectDelay == 0, 5*time.Second, r.cfg.Options.ReconnectDelay)
	for {
		r.connection, err = amqp.Dial(r.getConnectionString())
		if err == nil {
			break
		}

		r.log.Error(context.Background(), fmt.Sprintf("Failed to connect to RabbitMQ: %s Retrying in %d second", err, retryTime))

		time.Sleep(retryTime)
	}

	r.channel, err = r.connection.Channel()
	if err != nil {
		r.connection.Close()
		r.log.Fatal(context.Background(), fmt.Sprintf("Failed to open a RabbitMQ channel: %v", err))
	}

	r.log.Info(context.Background(), fmt.Sprintf("Connected to RabbitMQ: %s", r.getConnectionString()))
}

func (r *rabbitMQ) getConnectionString() string {
	protocol := "amqp"
	if r.cfg.SSL {
		protocol = "amqps"
	}

	connString := fmt.Sprintf("%s://%s:%s@%s:%d/",
		protocol,
		r.cfg.Username,
		r.cfg.Password,
		r.cfg.Host,
		r.cfg.Port)

	return connString
}

func (r *rabbitMQ) MonitorConnection() {
	connClosed := r.connection.NotifyClose(make(chan *amqp.Error, 1))
	chanClosed := r.channel.NotifyClose(make(chan *amqp.Error, 1))

	select {
	case err := <-connClosed:
		r.log.Error(context.Background(), fmt.Sprintf("RabbitMQ connection closed: %v", err))
		r.reconnect()
	case err := <-chanClosed:
		r.log.Error(context.Background(), fmt.Sprintf("RabbitMQ channel closed: %v", err))
		r.reconnect()
	}
}

func (r *rabbitMQ) reconnect() {
	// Close existing connections
	if r.channel != nil {
		r.channel.Close()
	}
	if r.connection != nil {
		r.connection.Close()
	}

	// Reconnect
	r.init()
}

// Stop closes the connection to RabbitMQ
func (r *rabbitMQ) Stop() {
	r.endOnce.Do(func() {
		if r.channel != nil {
			r.channel.Close()
		}
		if r.connection != nil {
			r.connection.Close()
		}
	})
}

// CreateExchange creates a new exchange
func (r *rabbitMQ) CreateExchange(name string, kind ExchangeType, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return r.channel.ExchangeDeclare(
		name,
		string(kind),
		durable,
		autoDelete,
		internal,
		noWait,
		args,
	)
}

// CreateQueue creates a new queue
func (r *rabbitMQ) CreateQueue(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)
}

// BindQueue binds a queue to an exchange with a routing key
func (r *rabbitMQ) BindQueue(queueName, exchangeName, routingKey string, noWait bool, args amqp.Table) error {
	return r.channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		noWait,
		args,
	)
}

// Publish sends a message to an exchange with a routing key
func (r *rabbitMQ) Publish(ctx context.Context, exchangeName, routingKey string, body string) error {
	return r.channel.PublishWithContext(
		ctx,
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			Body: []byte(body),
		},
	)
}

// Subscribe starts consuming messages from a queue
func (r *rabbitMQ) Subscribe(ctx context.Context, queueName string, handler MessageHandler) {
	delivery, err := r.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.log.Error(context.Background(), fmt.Sprintf("Failed to register consumer: %v", err))
	}

	go func() {
		for msg := range delivery {
			r.log.Info(ctx, fmt.Sprintf("Received message with exchange name: %v routing key:%v and body: %s", msg.Exchange, msg.RoutingKey, msg.Body))

			if err := handler(ctx, msg.Exchange, msg.RoutingKey, string(msg.Body)); err != nil {
				r.log.Error(ctx, fmt.Sprintf("Error handling message with exchange name: %v routing key:%v and body: %s", msg.Exchange, msg.RoutingKey, msg.Body))
			} else {
				r.log.Info(ctx, fmt.Sprintf("Success handling message with exchange name: %v routing key:%v and body: %s", msg.Exchange, msg.RoutingKey, msg.Body))
			}
		}
	}()

}
