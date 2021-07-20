/**
 * @Author: fuxiao
 * @Author: 576101059@qq.com
 * @Date: 2021/7/16 11:58
 * @Desc: An event bus library implemented by kafka.
 */

package eventbus

import (
	"fmt"
	"log"
	"os"
	"sync"
	
	"github.com/Shopify/sarama"
	
	"github.com/dobyte/event-bus/internal/conv"
)

type KafkaOptions struct {
	Addrs  []string
	Prefix string
}

type KafkaEventBus struct {
	BaseEventBus
	Listener sarama.Consumer
	Emitter  sarama.AsyncProducer
	prefix   string
	locker   sync.Mutex
	handlers map[string][]ListenHandlerFn
}

func NewKafkaEventBus(opt KafkaOptions) EventBus {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	
	consumer, err := sarama.NewConsumer(opt.Addrs, config)
	if err != nil {
		log.Fatalf("Create kafka's consumer failed: %s\n", err.Error())
	}
	
	producer, err := sarama.NewAsyncProducer(opt.Addrs, config)
	if err != nil {
		log.Fatalf("Create kafka's producer failed: %s\n", err.Error())
	}
	
	bus := &KafkaEventBus{}
	bus.Listener = consumer
	bus.Emitter = producer
	bus.prefix = opt.Prefix
	bus.signalChan = make(chan os.Signal, 1)
	bus.doneChan = make(chan bool, 1)
	bus.handlers = make(map[string][]ListenHandlerFn)
	
	return bus
}

// Emit Send data to topic.
func (e *KafkaEventBus) Emit(topic string, payload interface{}) error {
	message := &sarama.ProducerMessage{
		Topic: e.prefixTopic(topic),
		Value: sarama.StringEncoder(conv.String(payload)),
	}
	
	e.Emitter.Input() <- message
	
	select {
	case <-e.Emitter.Successes():
		return nil
	case err := <-e.Emitter.Errors():
		return err
	}
}

// Listen Monitor data on topic.
func (e *KafkaEventBus) Listen(topic string, handler ListenHandlerFn) error {
	isListened := false
	
	e.locker.Lock()
	if _, ok := e.handlers[topic]; ok {
		isListened = true
	} else {
		e.handlers[topic] = make([]ListenHandlerFn, 0)
	}
	e.handlers[topic] = append(e.handlers[topic], handler)
	e.locker.Unlock()
	
	if !isListened {
		partitions, err := e.Listener.Partitions(e.prefixTopic(topic))
		if err != nil {
			return err
		}
		
		for _, partition := range partitions {
			pc, err := e.Listener.ConsumePartition(e.prefixTopic(topic), partition, sarama.OffsetNewest)
			if err != nil {
				return err
			}
			
			go func(pc sarama.PartitionConsumer) {
				defer pc.AsyncClose()
				
				for message := range pc.Messages() {
					for _, handler := range e.handlers[topic] {
						go func(message *sarama.ConsumerMessage, handler ListenHandlerFn) {
							if err := handler(Payload{
								Value: message.Value,
							}); err != nil {
								log.Printf("Handle event's message failed: %s\n", err.Error())
							}
						}(message, handler)
					}
				}
			}(pc)
		}
	}
	
	return nil
}

// prefixTopic Add prefix to the front of topic.
func (e *KafkaEventBus) prefixTopic(topic string) string {
	if e.prefix == "" {
		return topic
	} else {
		return fmt.Sprintf("%s_%s", e.prefix, topic)
	}
}
