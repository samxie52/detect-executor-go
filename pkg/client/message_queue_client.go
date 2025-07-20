package client

import (
	"context"
	"detect-executor-go/pkg/config"
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// MessageQueueClient 消息队列客户端接口
type MessageQueueClient interface {
	BaseClient
	Publish(ctx context.Context, exchange, routingKey string, message interface{}) error
	Consume(ctx context.Context, queue string, handler MessageHandler) error
	// 声明队列 durable: 是否持久化
	DeclareQueue(ctx context.Context, queue string, durable bool) error
	// 声明交换机 exchangeType: 交换机类型 durable: 是否持久化
	DeclareExchange(ctx context.Context, exchange, exchangeType string, durable bool) error
}

// MessageHandler 消息处理函数
type MessageHandler func(ctx context.Context, message []byte) error

// messageQueueClient 消息队列客户端实现
type messageQueueClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  config.MessageQueueConfig
}

func NewMessageQueueClient(config config.MessageQueueConfig) (MessageQueueClient, error) {
	conn, err := amqp.Dial(config.URL)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to message queue: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %v", err)
	}

	return &messageQueueClient{
		conn:    conn,
		channel: channel,
		config:  config,
	}, nil

}

func (c *messageQueueClient) Publish(ctx context.Context, exchange, routingKey string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	err = c.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

func (c *messageQueueClient) Consume(ctx context.Context, queue string, handler MessageHandler) error {
	msgs, err := c.channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume message: %v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				//handler 处理消息
				if err := handler(ctx, msg.Body); err != nil {
					//msg.Nack(false, true) 表示拒绝消息
					msg.Nack(false, true)
				} else {
					//msg.Ack(false) 表示确认消息
					msg.Ack(false)
				}

			}
		}
	}()

	return nil
}

// 声明队列 durable: 是否持久化
func (c *messageQueueClient) DeclareQueue(ctx context.Context, queue string, durable bool) error {
	_, err := c.channel.QueueDeclare(
		queue,
		durable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}
	return nil
}

// 声明交换机 exchangeType: 交换机类型 durable: 是否持久化
func (c *messageQueueClient) DeclareExchange(ctx context.Context, exchange string, exchangeType string, durable bool) error {
	err := c.channel.ExchangeDeclare(
		exchange,
		exchangeType,
		durable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}
	return nil
}

// Close 关闭客户端
func (c *messageQueueClient) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}

// Ping 测试连接
func (c *messageQueueClient) Ping(ctx context.Context) error {
	if c.conn == nil || c.conn.IsClosed() {
		return fmt.Errorf("connection is closed")
	}
	if c.channel == nil {
		return fmt.Errorf("channel is nil")
	}
	return nil
}
