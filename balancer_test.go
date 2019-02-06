package golang_balancer

import (
	"github.com/pkg/errors"
	"testing"
	"time"
)

func getInstance() (b *Balancer) {
	var (
		ch      = make(chan interface{})
		handler = func(job interface{}) (err error) {
			println(job.(string))

			return
		}
		err = make(chan error)
	)

	b = NewBalancer(ch, handler, err, 10)

	return
}

func TestInit(t *testing.T) {
	b := getInstance()

	if b.CountPerSecond() != 10 {
		t.Fatal("Can't init balancer correctly")
	}
}

func TestIncrease(t *testing.T) {
	b := getInstance()

	b.Increase()
	b.Increase()

	if b.CountPerSecond() != 12 {
		t.Fatal("Can't increase amount goroutins in balancer")
	}
}

func TestDecrease(t *testing.T) {
	b := getInstance()

	b.Decrease()
	b.Decrease()

	if b.CountPerSecond() != 8 {
		t.Fatal("Can't decrease amount goroutins in balancer")
	}
}

func TestIncreaseWithLimitation(t *testing.T) {
	b := getInstance()

	b.SetMax(11)
	b.Increase()
	b.Increase()

	if b.CountPerSecond() != 11 {
		t.Fatal("Can't increase (with limitation) amount goroutins in balancer")
	}
}

func TestDecreaseWithLimitation(t *testing.T) {
	b := getInstance()

	b.SetMin(9)
	b.Decrease()
	b.Decrease()

	if b.CountPerSecond() != 9 {
		t.Fatal("Can't decrease (with limitation) amount goroutins in balancer")
	}
}

func TestErrChan(t *testing.T) {
	var (
		ch      = make(chan interface{})
		handler = func(job interface{}) (err error) {
			return errors.New("some error")
		}

		err = make(chan error)
	)

	NewBalancer(ch, handler, err, 10)

	ch <- "some job"

	time.AfterFunc(5*time.Second, func() {
		t.Fatal("Can't get error from channel")
	})

	<-err
}

func TestCountPerSecond(t *testing.T) {
	b := getInstance()
	b.SetCountPerSecond(11)

	if b.CountPerSecond() != 11 {
		t.Fatal("Can't get expected count per second value")
	}
}

func TestEfficiency(t *testing.T) {
	var (
		ch      = make(chan interface{})
		handler = func(job interface{}) (err error) {
			return
		}

		err = make(chan error)
	)

	b := NewBalancer(ch, handler, err, 10)

	time.AfterFunc(1100*time.Millisecond, func() {
		if b.Efficiency() != 10 {
			t.Fatal("Can't balance necessary amount jobs execution per second")
		}
	})

	for i := 0; i < 12; i++ {
		ch <- "some job"
	}
}
