/**
 * @Author: fuxiao
 * @Author: 576101059@qq.com
 * @Date: 2021/7/19 14:08
 * @Desc: TODO
 */

package main

import (
	"fmt"

	"github.com/dobyte/event-bus"
)

func main() {
	var (
		err   error
		count = 0
		bus   = eventbus.NewEventBus(eventbus.Options{
			Driver: eventbus.KafkaDriver,
			Prefix: "aoe",
			Supports: eventbus.Supports{
				Kafka: &eventbus.KafkaOptions{
					Addrs: []string{":9092"},
				},
			},
		})
	)

	err = bus.Listen("kafka", func(payload eventbus.Payload) error {
		count++

		fmt.Println(fmt.Sprintf("Receive event's message, count: %d, content:%s", count, payload.String()))

		return nil
	})

	if err != nil {
		fmt.Println(fmt.Sprintf("Listen topic failed: %s", err.Error()))
		return
	}

	bus.Wait()
}
