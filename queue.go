package compakq

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type OnInitFunc[I QItem[I]] func(queue *QueueStack[I]) error
type PullingFunc[I QItem[I]] func(queue *QueueStack[I]) (*I, error)
type HandlingFunc[I QItem[I]] func(queue *QueueStack[I], item *I) error
type PushBackFunc[I QItem[I]] func(queue *QueueStack[I], item *I) error
type OnPulledFunc[I QItem[I]] func(item *I)
type OnAckFunc[I QItem[I]] func(item *I)
type OnErrorFunc func(err error)

type QueueOptions[I QItem[I]] struct {
	Capacity int
	Throttle int64
	OnInit   OnInitFunc[I]
	Pulling  PullingFunc[I]
	Handling HandlingFunc[I]
	PushBack PushBackFunc[I]
	OnPulled OnPulledFunc[I]
	OnAck    OnAckFunc[I]
	OnError  OnErrorFunc
}

type Queue[I QItem[I]] struct {
	*QueueOptions[I]
	items *QueueStack[I]
	queue chan *I
}

func NewQueue[I QItem[I]](opts QueueOptions[I]) *Queue[I] {
	if opts.Capacity < 1 {
		opts.Capacity = 1
	}

	if opts.OnError == nil {
		opts.OnError = func(err error) { log.Panic(err) }
	}
	if opts.OnInit == nil {
		opts.OnInit = func(*QueueStack[I]) error { return nil }
	}
	if opts.OnPulled == nil {
		opts.OnPulled = func(*I) {}
	}
	if opts.OnAck == nil {
		opts.OnAck = func(*I) {}
	}

	if opts.PushBack == nil {
		opts.PushBack = func(items *QueueStack[I], item *I) error {
			items.Push(item)
			return nil
		}
	}

	q := &Queue[I]{
		QueueOptions: &opts,
	}

	q.items = NewQueueStack[I]()
	q.queue = make(chan *I, q.Capacity)

	return q
}

func (q *Queue[I]) startQueue() {
	var (
		stopSignal bool = false
		mu         sync.Mutex
		wg         *sync.WaitGroup = &sync.WaitGroup{}
		NUMCPU     int             = runtime.NumCPU()
	)

	defer close(q.queue)

	q.OnInit(q.items)

	defer wg.Wait()
	runtime.GOMAXPROCS(NUMCPU)

	for cpu := 1; cpu <= NUMCPU; cpu++ {
		wg.Add(1)

		runtime.LockOSThread()
		go func(cpu int) {
			defer wg.Done()

		exitLoop:
			for {
				time.Sleep(time.Duration(q.Throttle) * time.Millisecond)

				switch true {
				case cpu == Clamp(1, 1, NUMCPU):
					{
						if item, err := q.Pulling(q.items); err != nil {
							q.OnError(err)
							break exitLoop
						} else {
							if item == nil {
								continue
							}
							// litter.D(item)
							q.items.Push(item)
						}
						break
					}
				case cpu == Clamp(2, 1, NUMCPU):
					{
					exitFillCap:
						for i := 0; i < q.Capacity; i++ {
							time.Sleep(time.Duration(q.Throttle) * time.Millisecond)

							// continue popping item until queue reach its capacity
							if s := func() int {
								if q.items.Len() > 0 {
									item := q.items.Pop().(I)
									q.queue <- &item
									q.OnPulled(&item)
								}
								return len(q.queue)
							}(); s >= q.Capacity {
								// waiting for queue empty slot
								break exitFillCap
							} else {
								// continue filling up the queue
								continue
							}
						}
						break
					}
				default:
					{
						s := len(q.queue)
						// Handling leftover items on the queue
						if q.items.Len() == 0 && s > 0 {
							continue
						}
						mu.Lock()
						item := <-q.queue
						if err := q.Handling(q.items, item); err != nil {
							// put back item to the item pool
							if q.PushBack != nil {
								q.PushBack(q.items, item)
							}
							mu.Unlock()

							q.OnError(err)
							continue
						}
						mu.Unlock()

						go q.OnAck(item)
						break
					}
				}
				if stopSignal {
					break exitLoop
				}
			}
		}(cpu)
		runtime.UnlockOSThread()
	}
}

func (q *Queue[I]) WaitTermination() {
	go q.startQueue()

	exit := make(chan struct{})
	SignalC := make(chan os.Signal, 4)

	signal.Notify(
		SignalC,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		for s := range SignalC {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				close(exit)
				return
			}
		}
	}()

	<-exit
	os.Exit(0)
}
