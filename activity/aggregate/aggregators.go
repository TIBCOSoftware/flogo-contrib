package aggregate

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-contrib/activity/aggregate/window"
	"github.com/TIBCOSoftware/flogo-contrib/activity/aggregate/window/functions"
)

func NewTumblingWindow(function string, windowSize int) (window.Window, error) {
	switch function {
	case "avg":
		return window.NewTumblingWindow(functions.AddSampleSum, functions.AggregateSingleAvg, windowSize), nil
	case "sum":
		return window.NewTumblingWindow(functions.AddSampleSum, functions.AggregateSingleNoopFunc, windowSize), nil
	case "min":
		return window.NewTumblingWindow(functions.AddSampleMin, functions.AggregateSingleNoopFunc, windowSize), nil
	case "max":
		return window.NewTumblingWindow(functions.AddSampleMax, functions.AggregateSingleNoopFunc, windowSize), nil
	case "count":
		return window.NewTumblingWindow(functions.AddSampleCount, functions.AggregateSingleNoopFunc, windowSize), nil
	case "accumulate":
		return window.NewTumblingWindow(functions.AddSampleAccum, functions.AggregateSingleNoopFunc, windowSize), nil
	default:
		return nil, fmt.Errorf("unsupported function: %s", function)
	}
}

// NewTumblingTimeWindow creates a new tumbling time window, all time windows are managed
// externally and are progressed using the NextBlock() method
func NewTumblingTimeWindow(function string, windowTime int, externalTimer bool) (window.TimeWindow, error) {
	switch function {
	case "avg":
		return window.NewTumblingTimeWindow(functions.AddSampleSum, functions.AggregateSingleAvg, windowTime, externalTimer), nil
	case "sum":
		return window.NewTumblingTimeWindow(functions.AddSampleSum, functions.AggregateSingleNoopFunc, windowTime, externalTimer), nil
	case "min":
		return window.NewTumblingTimeWindow(functions.AddSampleMin, functions.AggregateSingleNoopFunc, windowTime, externalTimer), nil
	case "max":
		return window.NewTumblingTimeWindow(functions.AddSampleMax, functions.AggregateSingleNoopFunc, windowTime, externalTimer), nil
	case "count":
		return window.NewTumblingTimeWindow(functions.AddSampleCount, functions.AggregateSingleNoopFunc, windowTime, externalTimer), nil
	case "accumulate":
		return window.NewTumblingTimeWindow(functions.AddSampleAccum, functions.AggregateSingleNoopFunc, windowTime, externalTimer), nil
	default:
		return nil, fmt.Errorf("unsupported function: %s", function)
	}
}

func NewSlidingWindow(function string, windowSize int) (window.Window, error) {
	switch function {
	case "avg":
		return window.NewSlidingWindow(functions.AggregateBlocksAvg, windowSize, 1), nil
	case "sum":
		return window.NewSlidingWindow(functions.AggregateBlocksSum, windowSize, 1), nil
	case "min":
		return window.NewSlidingWindow(functions.AggregateBlocksMin, windowSize, 1), nil
	case "max":
		return window.NewSlidingWindow(functions.AggregateBlocksMax, windowSize, 1), nil
	case "count":
		return window.NewSlidingWindow(functions.AggregateBlocksCount, windowSize, 1), nil
	case "accumulate":
		return window.NewSlidingWindow(functions.AggregateBlocksAccumulate, windowSize, 1), nil
	default:
		return nil, fmt.Errorf("unsupported function: %s", function)
	}
}

// NewSlidingTimeWindow creates a new sliding time window, all time windows are managed
// externally and are progressed using the NextBlock() method
func NewSlidingTimeWindow(function string, windowTime int, resolution int, externalTimer bool) (window.TimeWindow, error) {
	switch function {
	case "avg":
		return window.NewSlidingTimeWindow(functions.AddSampleSum, functions.AggregateBlocksAvg, windowTime, resolution, externalTimer), nil
	case "sum":
		return window.NewSlidingTimeWindow(functions.AddSampleSum, functions.AggregateBlocksSum, windowTime, resolution, externalTimer), nil
	case "min":
		return window.NewSlidingTimeWindow(functions.AddSampleMin, functions.AggregateBlocksMin, windowTime, resolution, externalTimer), nil
	case "max":
		return window.NewSlidingTimeWindow(functions.AddSampleMax, functions.AggregateBlocksMax, windowTime, resolution, externalTimer), nil
	case "count":
		return window.NewSlidingTimeWindow(functions.AddSampleCount, functions.AggregateBlocksSum, windowTime, resolution, externalTimer), nil
	default:
		return nil, fmt.Errorf("unsupported function: %s", function)
	}
}
