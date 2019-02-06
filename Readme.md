# Golang way to control number goroutines executing per second with dynamic control

## Usage

```go
func main() {
	var (
		ch      = make(chan interface{})
		handler = func(job interface{}) (err error) {
			println(job.(string))

			return
		}
		err = make(chan error)
	)

	b = NewBalancer(ch, handler, err, 10)
	
	go func(){
		for e := range err {
		    log.Fatal(e)
		} 
	}()
	
	for i:= 0; i < 20; i++{
		ch <- "some job"
		
		if i%2 == 0{
			b.Increase()
		} else {
			b.Decrease()
		}
	}
}
```

## List of all methods

* golang_balancer.NewBalancer(jobCh, handler, errCh, countPerSecond) - Creates new balancer with jobs channel, function to handle new job, errors channel and number executing goroutines per second 
* golang_balancer.SetMax(number) - Sets maximum of goroutines 
* golang_balancer.SetMin(number) - Set minimum of goroutine
* golang_balancer.Increase() - Increases number of goroutine
* golang_balancer.Decrease() - Decreases number of goroutine
* golang_balancer.SetCountPerSecond(number) - Sets number of goroutine
* golang_balancer.CountPerSecond() - Returns number of goroutine
* golang_balancer.Efficiency() - Returns real efficiency: number executed jobs per second (refreshes every second)
* golang_balancer.Wait() - Waits of all left jobs execution