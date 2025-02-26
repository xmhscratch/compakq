package compakq

import (
	"testing"
	"time"
)

func TestMe(t *testing.T) {
	q := NewQueue(QueueOptions[Fruit]{
		Capacity: 3,
		Throttle: 500,
		OnInit: func(queue *QueueStack[Fruit]) error {
			queue.Push(&Fruit{T: "Apple"})
			queue.Push(&Fruit{T: "Banana"})
			queue.Push(&Fruit{T: "Cherry"})
			queue.Push(&Fruit{T: "Grape"})
			queue.Push(&Fruit{T: "Mango"})
			queue.Push(&Fruit{T: "Orange"})

			return nil
		},
		Pulling: func(queue *QueueStack[Fruit]) (item *Fruit, err error) {
			queue.Push(&Fruit{T: "Banana"})
			queue.Push(&Fruit{T: "Mango"})

			return nil, err
		},
		Handling: func(queue *QueueStack[Fruit], item *Fruit) error {
			time.Sleep(time.Duration(2) * time.Second)
			return nil
		},
		OnPulled: func(item *Fruit) {
			t.Log("Item pushed", item.Key())
		},
		OnAck: func(item *Fruit) {
			t.Log("Item consumed", item.Key())
		},
		OnError: func(err error) {
			t.Fatal(err)
		},
	})

	q.WaitTermination()
}
