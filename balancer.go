package golang_balancer

import (
	"sync"
	"sync/atomic"
	"time"
)

type Balancer struct {
	countPerSecond int32
	min            int32
	max            int32
	efficiency     int32

	handler func(job interface{}) (err error)

	err chan error

	wg sync.WaitGroup
}

func NewBalancer(jobCh chan interface{}, handler func(job interface{}) (err error), errCh chan error, countPerSecond int32) (b *Balancer) {
	b = &Balancer{
		handler:        handler,
		countPerSecond: countPerSecond,
		max:            2 * countPerSecond,
		min:            0,

		err: errCh,
	}

	var counter, efficiency int32
	ticker := time.NewTicker(time.Second)

	go func() {
		for range ticker.C {
			efficiency = atomic.LoadInt32(&counter)
			atomic.StoreInt32(&counter, 0)

			atomic.StoreInt32(&b.efficiency, efficiency)
		}
	}()

	go func() {
		var (
			err     error
			sleep   time.Duration
			current int32
		)

		for job := range jobCh {
			b.wg.Add(1)

			go func() {
				err = handler(job)
				if err != nil {
					b.err <- err
				}

				b.wg.Done()

				atomic.AddInt32(&counter, 1)
			}()

			current = atomic.LoadInt32(&b.countPerSecond)

			sleep = time.Second / time.Duration(current)

			time.Sleep(sleep)
		}
	}()

	return
}

func (b *Balancer) SetMax(max int32) (ok bool) {
	if max < atomic.LoadInt32(&b.countPerSecond) {
		return
	}

	atomic.StoreInt32(&b.max, max)

	return true
}

func (b *Balancer) SetMin(min int32) (ok bool) {
	if min > atomic.LoadInt32(&b.countPerSecond) {
		return
	}

	atomic.StoreInt32(&b.min, min)

	return true
}

func (b *Balancer) Increase() (ok bool) {
	if atomic.LoadInt32(&b.countPerSecond) >= atomic.LoadInt32(&b.max) {
		return
	}

	atomic.AddInt32(&b.countPerSecond, 1)

	return true
}

func (b *Balancer) Decrease() (ok bool) {
	if atomic.LoadInt32(&b.countPerSecond) <= atomic.LoadInt32(&b.min) {
		return
	}

	atomic.AddInt32(&b.countPerSecond, -1)

	return true
}

func (b *Balancer) SetCountPerSecond(countPerSecond int32) (ok bool) {
	if countPerSecond > atomic.LoadInt32(&b.max) || countPerSecond < atomic.LoadInt32(&b.min) {
		return
	}

	atomic.StoreInt32(&b.countPerSecond, countPerSecond)

	return true
}

func (b *Balancer) CountPerSecond() (countPerSecond int32) {
	return atomic.LoadInt32(&b.countPerSecond)
}

func (b *Balancer) Efficiency() (efficiency int32) {
	return atomic.LoadInt32(&b.efficiency)
}

func (b *Balancer) Wait() {
	b.wg.Wait()
}
