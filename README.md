# CompakQ

Compakq is a minimal queue manager written in Golang. The package is keep simple as possible
and provide mechanism for working with Redis queue or with any custom message broker.

## Installation

<pre>
    go get -u github.com/xmhscratch/compakq
</pre>

Import in go code

```go
    import (
        "github.com/xmhscratch/compakq"
    )
```

## Usage

Define queue item which is a heap stack for quick retrieving. The Index() and Key() interface implementations is mandatory

```go
type TFruitName (string)

type MyFruit struct {
	QItem[MyFruit]
	T TFruitName
}

// implements String() method for print out purposes
func (f TFruitName) Eat() string {
	const (
		Apple  TFruitName = "Apple"
		Banana TFruitName = "Banana"
		Cherry TFruitName = "Cherry"
		Grape  TFruitName = "Grape"
		Mango  TFruitName = "Mango"
		Orange TFruitName = "Orange"
	)

	return map[TFruitName]string{
		Apple:  "Apple",
		Banana: "Banana",
		Cherry: "Cherry",
		Grape:  "Grape",
		Mango:  "Mango",
		Orange: "Orange",
	}[f]
}

// heap stack retrievals is depending on item creation timestamp
func (ctx MyFruit) Index() int {
	return int(time.Now().Unix())
}

// provide unique identity for the item
func (ctx MyFruit) Key() string {
	return ctx.T.String()
}
```

## Queue spin-up

```go
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
```

# Queue options

## Capacity

Number of items are processing in parallel

```go
QueueOptions[MyFruit]{
    Capacity: 3,
}
```

[Back to TOC](#installation)

## Throttle

Heatbeat pulses in miliseconds. Used for new items pulling.

```go
QueueOptions[MyFruit]{
    Throttle: 500,
}
```

[Back to TOC](#installation)

## OnInit: OnInitFunc[I]

Execute before the queue start. Using for define default queue states.

```go
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
    QueueOptions[MyFruit]{
        OnInit: onInit,
    }
```

[Back to TOC](#installation)

## OnPulling: PullingFunc[I](*QueueStack[I]) (\*I, error)

Periodical pulling new items from message queue storage. (Eg. Redis cluster). New retrieval returned item will be pushed to the queue.

```go
	onPulling := func() PullingFunc[MyFruit] {
		return func(queue *QueueStack[MyFruit]) (*MyFruit, error) {
            // pulling item from redis queue

            var (
                err    error
                qItem  string
                fruit     *MyFruit
                hasFruit int64 = 0
            )

            rdmQueueKey := "eat-my-fruits:queue"
            if hasFruit, err = rdm.Exists(context.TODO(), rdmQueueKey).Result(); err != nil || hasFruit == 0 {
                return nil, nil
            }

            if qItem, err = rdm.SPop(
                context.TODO(),
                rdmQueueKey,
            ).Result(); err != nil {
                return nil, err
            } else {
                if err = json.Unmarshal([]byte(qItem), &fruit); err != nil {
                    return nil, err
                }
            }

            return &MyFruit{T:fruit}, err
		}
	}
```

[Back to TOC](#installation)

## OnHandling: HandlingFunc[I](*QueueStack[I], *I) error

Handling the heavy item processing. (Eg. minify image, encode video, ...)

```go
	onHandling := func() HandlingFunc[MyFruit] {
		return func(queue *QueueStack[MyFruit], fruit *MyFruit) error {
			time.Sleep(time.Duration(2) * time.Second)
            fmt.Println(fruit.Eat())

			return nil
		}
	}
```

[Back to TOC](#installation)

## PushBack: PushBackFunc[I] error

Happen on item acknowledgement was unsuccessful (Nack). Item is pushed back to the queue stack or customise pushing back to message storage (redis cluster).

```go
	onPushBack := func() PushBackFunc[MyFruit] error {
		return func(queue *QueueStack[MyFruit], fruit *MyFruit) error {
            rdmQueueKey := "eat-my-fruits:queue"

            if qItem, err := json.Marshal(fruit); err != nil {
                return err
            } else {
                if err := rdm.SAdd(
                    context.TODO(),
                    rdmQueueKey,
                    qItem,
                ).Err(); err != nil {
                    return err
                }
            }

            return nil
		}
	}
```

[Back to TOC](#installation)

## OnPulled: OnPulledFunc[I](*I)

Called when item is added to the queue.

```go
	onPulled := func() OnPulledFunc[MyFruit] {
		return func(fruit *MyFruit) {
			fmt.Println("Fruit added", fruit.Key())
		}
	}
```

[Back to TOC](#installation)

## OnAck: OnAckFunc[I](*I)

Called when item is done processed.

```go
	onAcknowledge := func() OnAckFunc[MyFruit] {
		return func(fruit *MyFruit) {
			fmt.Println("Fruit consumed", fruit.Key())
		}
	}
```

[Back to TOC](#installation)

## OnError: OnErrorFunc(error)

Error handling for the queue

```go
	onError := func() OnErrorFunc {
		return func(err error) {
			log.Fatal(err)
		}
	}
```

## 📜 License

This project is licensed under the Apache 2.0 - see the LICENSE file for details.
