/**
 * @Author: fuxiao
 * @Author: 576101059@qq.com
 * @Date: 2021/7/19 17:13
 * @Desc: An event bus library implemented by redis.
 */

package eventbus

import (
	"context"
	"fmt"
	"log"
	"os"
	
	"github.com/go-redis/redis/v8"
	
	"github.com/dobyte/event-bus/internal/conv"
)

type RedisOptions struct {
	Addrs    []string
	Username string
	Password string
	DB       int
	Prefix   string
}

type RedisEventBus struct {
	BaseEventBus
	client redis.UniversalClient
	prefix string
}

func NewRedisEventBus(opt RedisOptions) EventBus {
	bus := &RedisEventBus{}
	bus.prefix = opt.Prefix
	bus.signalChan = make(chan os.Signal, 1)
	bus.doneChan = make(chan bool, 1)
	bus.client = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    opt.Addrs,
		Username: opt.Username,
		Password: opt.Password,
		DB:       opt.DB,
	})
	
	return bus
}

// Emit Send data to topic.
func (e *RedisEventBus) Emit(topic string, payload interface{}) error {
	return e.client.Publish(context.Background(), e.prefixTopic(topic), conv.String(payload)).Err()
}

// Listen Monitor data on topic.
func (e *RedisEventBus) Listen(topic string, handler listenHandlerFn) error {
	var (
		ctx    = context.Background()
		pubSub = e.client.Subscribe(ctx, e.prefixTopic(topic))
	)
	
	_, err := pubSub.Receive(ctx)
	if err != nil {
		return err
	}
	
	go func() {
		for message := range pubSub.Channel() {
			go func(message *redis.Message) {
				if err := handler(Payload{
					Value: []byte(message.Payload),
				}); err != nil {
					log.Printf("Handle event's message failed: %s\n", err.Error())
				}
			}(message)
		}
	}()
	
	return nil
}

// prefixTopic Add prefix to the front of topic.
func (e *RedisEventBus) prefixTopic(topic string) string {
	if e.prefix == "" {
		return topic
	} else {
		return fmt.Sprintf("%s:%s", e.prefix, topic)
	}
}