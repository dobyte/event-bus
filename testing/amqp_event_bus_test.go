/**
 * @Author: fuxiao
 * @Author: 576101059@qq.com
 * @Date: 2021/7/20 11:51
 * @Desc: A testing for kafka event bus.
 */

package testing

import (
	"fmt"
	"testing"
	"time"

	eventbus "github.com/dobyte/event-bus"
)

func TestAmqpEventBus_Emit(t *testing.T) {
	var bus = newAmpqEventBus()

	for i := 0; i < 10; i++ {
		err := bus.Emit(defaultEventTopic, &Event{Name: "register"})
		if err != nil {
			t.Fatal(fmt.Sprintf("Emit event's message failed: %s.", err.Error()))
		}

		t.Log(fmt.Sprintf("Emit event's message success: %d", i))

		time.Sleep(2 * time.Second)
	}
}

func TestAmqpEventBus_Listen(t *testing.T) {
	var (
		err   error
		count int
		bus   = newAmpqEventBus()
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
