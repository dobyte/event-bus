/**
 * @Author: fuxiao
 * @Author: 576101059@qq.com
 * @Date: 2021/7/20 10:02
 * @Desc: An event bus library implemented by amqp.
 */

package eventbus

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"

	"github.com/dobyte/event-bus/internal/conv"
)

type AmqpOptions struct {
	Addr     string
	Username string
	Password string
	Vhost    string
	Prefix   string
}

type AmqpEventBus struct {
	BaseEventBus
	prefix             string
	connection         *amqp.Connection
	channel            *amqp.Channel
	isDeclaredExchange bool
}

func NewAmqpEventBus(opt AmqpOptions) EventBus {
	vhost := opt.Vhost
	if vhost == "" {
		vhost = "/"
	}

	connection, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s%s", opt.Username, opt.Password, opt.Addr, vhost))
	if err != nil {
		log.Fatalf("Connect amqp's server failed: %s\n", err.Error())
	}

	channel, err := connection.Channel()
	if err != nil {
		_ = connection.Close()
		log.Fatalf("Open amqp's channel failed: %s\n", err.Error())
	}

	bus := &AmqpEventBus{}
	bus.connection = connection
	bus.channel = channel
	bus.prefix = opt.Prefix
	bus.signalChan = make(chan os.Signal, 1)
	bus.doneChan = make(chan bool, 1)

	return bus
}

// Emit Send data to topic.
func (e *AmqpEventBus) Emit(topic string, payload interface{}) error {
	if e.isDeclaredExchange == false {
		err := e.channel.ExchangeDeclare(e.prefixTopic(topic), amqp.ExchangeFanout, true, false, false, false, nil)
		if err != nil {
			return err
		}

		e.isDeclaredExchange = true
	}

	return e.channel.Publish(e.prefixTopic(topic), "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(conv.String(payload)),
	})
}

// Listen Monitor data on topic.
func (e *AmqpEventBus) Listen(topic string, handler listenHandlerFn) error {
	topic = e.prefixTopic(topic)

	err := e.channel.ExchangeDeclare(topic, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		return err
	}

	queue, err := e.channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return err
	}

	err = e.channel.QueueBind(queue.Name, "", topic, false, nil)
	if err != nil {
		return err
	}

	messages, err := e.channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range messages {
			go func(delivery amqp.Delivery, handler listenHandlerFn) {
				if err := handler(Payload{
					Value: delivery.Body,
				}); err != nil {
					log.Printf("Handle event's message failed: %s\n", err.Error())
				}
			}(delivery, handler)
		}
	}()

	return nil
}

// prefixTopic Add prefix to the front of topic.
func (e *AmqpEventBus) prefixTopic(topic string) string {
	if e.prefix == "" {
		return topic
	} else {
		return fmt.Sprintf("%s.%s", e.prefix, topic)
	}
}
