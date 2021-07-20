/**
 * @Author: fuxiao
 * @Author: 576101059@qq.com
 * @Date: 2021/7/19 11:36
 * @Desc: A testing for kafka event bus.
 */

package testing

import (
	"fmt"
	"testing"
	"time"

	"github.com/dobyte/event-bus"
)

func TestKafkaEventBus_Emit(t *testing.T) {
	var bus = newKafkaEventBus()

	for i := 0; i < 10; i++ {
		err := bus.Emit(defaultEventTopic, &Event{Name: "register"})
		if err != nil {
			t.Fatal(fmt.Sprintf("Emit event's message failed: %s.", err.Error()))
		}

		t.Log(fmt.Sprintf("Emit event's message success: %d", i))

		time.Sleep(2 * time.Second)
	}
}

func TestKafkaEventBus_Listen(t *testing.T) {
	var (
		err   error
		count int
		bus   = newKafkaEventBus()
	)

	err = bus.Listen(defaultEventTopic, func(payload eventbus.Payload) error {
		count++

		fmt.Println(fmt.Sprintf("Receive event's message, id: %d count: %d, content:%s", 1, count, payload.String()))

		return nil
	})

	err = bus.Listen(defaultEventTopic, func(payload eventbus.Payload) error {
		count++

		fmt.Println(fmt.Sprintf("Receive event's message, id: %d count: %d, content:%s", 2, count, payload.String()))

		return nil
	})

	if err != nil {
		fmt.Println(fmt.Sprintf("Listen topic failed: %s", err.Error()))
		return
	}

	bus.Wait()
}
