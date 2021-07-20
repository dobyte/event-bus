/**
 * @Author: fuxiao
 * @Author: 576101059@qq.com
 * @Date: 2021/7/19 14:19
 * @Desc: TODO
 */

package main

import (
	"fmt"

	"github.com/dobyte/event-bus"
)

func main() {
	var (
		err error
		bus = eventbus.NewEventBus(eventbus.Options{
			Driver: eventbus.KafkaDriver,
			Prefix: "aoe",
			Supports: eventbus.Supports{
				Kafka: &eventbus.KafkaOptions{
					Addrs: []string{":9092"},
				},
			},
		})
	)

	err = bus.Emit("kafka", "this is a test log")
	if err != nil {
		fmt.Println(fmt.Sprintf("Emit event's message failed: %s.", err.Error()))
		return
	}

	fmt.Println("Emit event's message success.")
}
