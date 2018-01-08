package aggregator

import (
	"testing"
	"sync"
	"time"
	"fmt"
)

func TestTimeBlockAverage_Add(t *testing.T) {

	agg := NewTimeBlockAverage(25)
	results := make(chan float64)
	var wg sync.WaitGroup

	wg.Add(6)
	go func() {
		defer wg.Done()
		_,avg := agg.Add(5.0)
		results <- avg
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 5)
		_,avg := agg.Add(10.0)
		results <- avg
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 10)
		_,avg := agg.Add(15.0)
		results <- avg
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 30)
		_,avg := agg.Add(6.0)
		results <- avg
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 35)
		_,avg := agg.Add(12.0)
		results <- avg
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 40)
		_,avg := agg.Add(18.0)
		results <- avg
	}()

	go func() {
		for i := range results {
			fmt.Println(i)
		}
	}()

	wg.Wait()
}
