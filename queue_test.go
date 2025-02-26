package compakq

import (
	"log"
	"testing"
	"time"
)

func TestMe(t *testing.T) {
	onInit := func() OnInitFunc[MyFruit] {
		return func(queue *QueueStack[MyFruit]) error {
			queue.Push(&MyFruit{T: "Apple"})
			queue.Push(&MyFruit{T: "Banana"})
			queue.Push(&MyFruit{T: "Cherry"})
			queue.Push(&MyFruit{T: "Grape"})
			queue.Push(&MyFruit{T: "Mango"})
			queue.Push(&MyFruit{T: "Orange"})

			return nil
		}
	}
	onPulling := func() PullingFunc[MyFruit] {
		return func(queue *QueueStack[MyFruit]) (*MyFruit, error) {
			return nil, nil
		}
	}
	onHandling := func() HandlingFunc[MyFruit] {
		return func(queue *QueueStack[MyFruit], item *MyFruit) error {
			time.Sleep(time.Duration(2) * time.Second)
			return nil
		}
	}
	onPulled := func() OnPulledFunc[MyFruit] {
		return func(item *MyFruit) {
			t.Log("Item pushed", item.Key())
		}
	}
	onAcknowledge := func() OnAckFunc[MyFruit] {
		return func(item *MyFruit) {
			t.Log("Item consumed", item.Key())
		}
	}
	onError := func() OnErrorFunc {
		return func(err error) {
			log.Fatal(err)
		}
	}

	q := NewQueue(QueueOptions[MyFruit]{
		Capacity: 3,
		Throttle: 500,
		OnInit:   onInit(),
		Pulling:  onPulling(),
		Handling: onHandling(),
		OnPulled: onPulled(),
		OnAck:    onAcknowledge(),
		OnError:  onError(),
	})

	q.WaitTermination()
}
