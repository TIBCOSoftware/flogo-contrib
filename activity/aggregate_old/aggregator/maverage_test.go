package aggregator

import (
	"fmt"
	"testing"
)

func TestMovingAverage_Add(t *testing.T) {

	agg := NewMovingAverage(5)

	report, avg := agg.Add(10)
	if report {
		t.Error("Window should not report yet")
	}
	report, avg = agg.Add(20)
	if report {
		t.Error("Window should not report yet")
	}
	report, avg = agg.Add(30)
	if report {
		t.Error("Window should not report yet")
	}
	report, avg = agg.Add(40)
	if report {
		t.Error("Window should not report yet")
	}
	report, avg = agg.Add(50)

	if avg != 30.0 {
		t.Error("Average should be 30")
	}

	report, avg = agg.Add(60)

	fmt.Println("avg:", avg)

	if avg != 40.0 {
		t.Error("Average should be 40")
	}
}
