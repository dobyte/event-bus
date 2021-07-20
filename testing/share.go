/**
 * @Author: fuxiao
 * @Author: 576101059@qq.com
 * @Date: 2021/7/20 12:49
 * @Desc: TODO
 */

package testing

import eventbus "github.com/dobyte/event-bus"

type Event struct {
	Name string
}

const defaultEventTopic = "abced"

func newAmpqEventBus() eventbus.EventBus {
	return eventbus.NewAmqpEventBus(eventbus.AmqpOptions{
		Addr:     ":5672",
		Username: "guest",
		Password: "guest",
		Prefix:   "aoe",
	})
}

func newRedisEventBus() eventbus.EventBus {
	return eventbus.NewRedisEventBus(eventbus.RedisOptions{
		Addrs:  []string{":7000"},
		Prefix: "aoe",
	})
}

func newKafkaEventBus() eventbus.EventBus {
	return eventbus.NewKafkaEventBus(eventbus.KafkaOptions{
		Addrs:  []string{":9092"},
		Prefix: "aoe",
	})
}