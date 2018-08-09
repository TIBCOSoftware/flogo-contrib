package window

import (
	"fmt"
	"time"
)

///////////////////
// Tumbling Window

func NewTumblingWindow(addFunc AddSampleFunc, aggFunc AggregateSingleFunc, windowSize int) Window {

	return &TumblingWindow{addFunc: addFunc, aggFunc: aggFunc, windowSize: windowSize}
}

//note:  using interface{} 4x slower than using specific types, starting with interface{} for expediency
type TumblingWindow struct {
	addFunc    AddSampleFunc
	aggFunc    AggregateSingleFunc
	windowSize int

	data       interface{}
	numSamples int
}

func (w *TumblingWindow) AddSample(sample interface{}) (bool, interface{}) {

	//sample size should match data size
	w.data = w.addFunc(w.data, sample)
	w.numSamples++

	if w.numSamples == w.windowSize {
		// aggregate and emit
		val := w.aggFunc(w.data, w.windowSize)

		w.numSamples = 0
		w.data, _ = zero(w.data)

		return true, val
	}

	return false, nil
}

///////////////////////
// Tumbling Time Window

func NewTumblingTimeWindow(addFunc AddSampleFunc, aggFunc AggregateSingleFunc, windowTime int, externalTimer bool) TimeWindow {
	return &TumblingTimeWindow{addFunc: addFunc, aggFunc: aggFunc, windowTime: windowTime, externalTimer: externalTimer}
}

// TumblingTimeWindow - A tumbling window based on time. Relies on external entity moving window along
// by calling NextBlock at the appropriate time.
//note:  using interface{} 4x slower than using specific types, starting with interface{} for expediency
type TumblingTimeWindow struct {
	addFunc       AddSampleFunc
	aggFunc       AggregateSingleFunc
	windowTime    int
	externalTimer bool

	data       interface{}
	maxSamples int
	numSamples int

	nextEmit int
	lastAdd  int
}

func (w *TumblingTimeWindow) AddSample(sample interface{}) (bool, interface{}) {

	w.data = w.addFunc(w.data, sample)
	w.numSamples++

	if w.numSamples > w.maxSamples {
		w.maxSamples = w.numSamples
	}

	if !w.externalTimer {
		currentTime := getTimeMillis()

		//todo what do we do if this greatly exceeds the nextEmit time?
		if currentTime >= w.nextEmit {
			w.nextEmit = + w.windowTime
			return w.NextBlock()
		}
	}

	return false, nil
}

func (w *TumblingTimeWindow) NextBlock() (bool, interface{}) {

	// aggregate and emit
	val := w.aggFunc(w.data, w.maxSamples) //num samples or max samples?

	w.numSamples = 0
	w.data, _ = zero(w.data)

	return true, val
}

///////////////////
// Sliding Window

func NewSlidingWindow(aggFunc AggregateBlocksFunc, windowSize int, resolution int) Window {

	w := &SlidingWindow{aggFunc: aggFunc, windowSize: windowSize, resolution: resolution}
	w.blocks = make([]interface{}, windowSize)

	return w
}

//note:  using interface{} 4x slower than using specific types, starting with interface{} for expediency
// todo split external vs on-add timer
type SlidingWindow struct {
	aggFunc    AggregateBlocksFunc
	windowSize int
	resolution int

	blocks       []interface{}
	numSamples   int
	currentBlock int
	canEmit      bool
}

func (w *SlidingWindow) AddSample(sample interface{}) (bool, interface{}) {

	//sample size should match data size
	w.blocks[w.currentBlock] = sample //no addSampleFunc required, just tracking all values

	if !w.canEmit {
		if w.currentBlock == w.windowSize-1 {
			w.canEmit = true
		}
	}

	w.numSamples++

	if w.canEmit && w.numSamples >= w.resolution {

		// aggregate and emit
		val := w.aggFunc(w.blocks, 1)

		w.numSamples = 0
		w.currentBlock++

		w.currentBlock = w.currentBlock % w.windowSize

		return true, val
	}

	w.currentBlock++

	return false, nil
}

//////////////////////
// Sliding Time Window

func NewSlidingTimeWindow(addFunc AddSampleFunc, aggFunc AggregateBlocksFunc, windowTime int, resolution int, externalTimer bool) TimeWindow {

	numBlocks := windowTime / resolution

	w := &SlidingTimeWindow{addFunc: addFunc, aggFunc: aggFunc, numBlocks: numBlocks, windowTime: windowTime,
		windowResolution: resolution, externalTimer: externalTimer}

	w.blocks = make([]interface{}, numBlocks)

	return w
}

// SlidingTimeWindow - A sliding window based on time. Relies on external entity moving window along
// by calling NextBlock at the appropriate time.
// note:  using interface{} 4x slower than using specific types, starting with interface{} for expediency
type SlidingTimeWindow struct {
	addFunc          AddSampleFunc
	aggFunc          AggregateBlocksFunc
	numBlocks        int
	windowTime       int
	windowResolution int
	externalTimer    bool

	blocks       []interface{}
	maxSamples   int
	numSamples   int
	currentBlock int
	canEmit      bool

	nextBlockTime int
	lastAdd       int
}

func (w *SlidingTimeWindow) AddSample(sample interface{}) (bool, interface{}) {

	//sample size should match data size
	w.blocks[w.currentBlock] = w.addFunc(w.blocks[w.currentBlock], sample)

	w.numSamples++

	if w.numSamples > w.maxSamples {
		w.maxSamples = w.numSamples
	}

	if !w.externalTimer {
		currentTime := getTimeMillis()

		if currentTime > w.nextBlockTime {
			w.nextBlockTime += w.windowResolution
			return w.NextBlock()
		}

		return false, nil
	}

	return false, nil
}

func (w *SlidingTimeWindow) NextBlock() (bool, interface{}) {

	if !w.canEmit {
		if w.currentBlock == w.numBlocks-1 {
			w.canEmit = true
		}
	}

	w.numSamples = 0
	w.currentBlock++

	if w.canEmit {

		// aggregate and emit
		val := w.aggFunc(w.blocks, w.maxSamples)

		w.currentBlock = w.currentBlock % w.numBlocks
		w.blocks[w.currentBlock], _ = zero(w.blocks[w.currentBlock])
		return true, val
	}

	return false, nil
}

///////////////////
// utils

func zero(a interface{}) (interface{}, error) {
	switch x := a.(type) {
	case int:
		return 0, nil
	case float64:
		return 0.0, nil
	case []int:
		for idx := range x {
			x[idx] = 0
		}
		return x, nil
	case []float64:
		for idx := range x {
			x[idx] = 0.0
		}
		return x, nil
	}

	return nil, fmt.Errorf("unsupported type")
}

func getTimeMillis() int {
	now := time.Now()
	nano := now.Nanosecond()
	return nano / 1000000
}
