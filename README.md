# Event Bus
git

## Use

Download and install

```shell script
go get github.com/dobyte/event-bus
```

API

```go
// Emit Send data to topic.
Emit(topic string, payload interface{}) error
// Listen Monitor data on topic.
Listen(topic string, handler listenHandlerFn) error
// Wait Synchronous waiting.
Wait()
// Exit Exit monitoring.
Exit()
```

Listen Dome

```go
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
```

Emit Dome
```go
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
```

## Example

- View emit demo [example/emit/emit.go](example/emit/emit.go)
- View listen demo [example/listen/listen.go](example/listen/listen.go)