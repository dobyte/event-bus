/**
 * @Author: fuxiao
 * @Author: 576101059@qq.com
 * @Date: 2021/7/16 11:46
 * @Desc: event bus interface define.
 */

package eventbus

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type EventBus interface {
	// Emit Send data to topic.
	Emit(topic string, payload interface{}) error
	// Listen Monitor data on topic.
	Listen(topic string, handler ListenHandlerFn) error
	// Wait Synchronous waiting.
	Wait()
	// Exit Exit monitoring.
	Exit()
}

const (
	RedisDriver = "redis"
	KafkaDriver = "kafka"
	AmqpDriver  = "amqp"
)

type ListenHandlerFn func(payload Payload) error

type Supports struct {
	Kafka *KafkaOptions
	Redis *RedisOptions
	Amqp  *AmqpOptions
}

type Options struct {
	Driver   string
	Prefix   string
	Supports Supports
}

// NewEventBus Create a event bus instance.
func NewEventBus(opt Options) EventBus {
	switch opt.Driver {
	case RedisDriver:
		return newRedisEventBus(opt)
	case KafkaDriver:
		return newKafkaEventBus(opt)
	case AmqpDriver:
		return newAmqpEventBus(opt)
	default:
		log.Fatalf("Invalid driver mode.")
	}

	return nil
}

// Create a event bus instance by redis.
func newRedisEventBus(opt Options) EventBus {
	option := RedisOptions{
		Addrs:    opt.Supports.Redis.Addrs,
		Username: opt.Supports.Redis.Username,
		Password: opt.Supports.Redis.Password,
		DB:       opt.Supports.Redis.DB,
		Prefix:   opt.Prefix,
	}

	if opt.Supports.Redis.Prefix != "" {
		option.Prefix = opt.Supports.Redis.Prefix
	}

	return NewRedisEventBus(option)
}

// Create a event bus instance by kafka.
func newKafkaEventBus(opt Options) EventBus {
	option := KafkaOptions{
		Addrs:  opt.Supports.Kafka.Addrs,
		Prefix: opt.Prefix,
	}

	if opt.Supports.Kafka.Prefix != "" {
		option.Prefix = opt.Supports.Kafka.Prefix
	}

	return NewKafkaEventBus(option)
}

// Create a event bus instance by amqp.
func newAmqpEventBus(opt Options) EventBus {
	option := AmqpOptions{
		Addr:     opt.Supports.Amqp.Addr,
		Username: opt.Supports.Amqp.Username,
		Password: opt.Supports.Amqp.Password,
		Vhost:    opt.Supports.Amqp.Vhost,
		Prefix:   opt.Prefix,
	}

	if opt.Supports.Amqp.Prefix != "" {
		option.Prefix = opt.Supports.Amqp.Prefix
	}

	return NewAmqpEventBus(option)
}

type BaseEventBus struct {
	signalChan chan os.Signal
	doneChan   chan bool
}

// Wait Synchronous waiting.
func (e *BaseEventBus) Wait() {
	signal.Notify(e.signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-e.signalChan
		e.doneChan <- true
	}()

	<-e.signalChan
}

// Exit Exit monitoring.
func (e *BaseEventBus) Exit() {
	e.doneChan <- true
}
