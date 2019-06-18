package ratchet

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/licaonfee/ratchet/data"
	"github.com/licaonfee/ratchet/processors"
)

// ProcessorWrapper is a type used internally to the Pipeline management
// code, and wraps a DataProcessor instance. DataProcessor is the main
// interface that should be implemented to perform work within the data
// pipeline, and this ProcessorWrapper type simply embeds it and adds some
// helpful channel management and other attributes.
type ProcessorWrapper struct {
	processors.DataProcessor
	ExecutionStat
	concurrentDataProcessor
	chanBrancher
	chanMerger
	outputs    []processors.DataProcessor
	inputChan  chan data.JSON
	outputChan chan data.JSON
}

type chanBrancher struct {
	branchOutChans []chan data.JSON
}

func (dp *ProcessorWrapper) branchOut() {
	go func() {
		for d := range dp.outputChan {
			for _, out := range dp.branchOutChans {
				// Make a copy to ensure concurrent stages
				// can alter data as needed.
				dc := make(data.JSON, len(d))
				copy(dc, d)
				out <- dc
			}
			dp.recordDataSent(d)
		}
		// Once all data is received, also close all the outputs
		for _, out := range dp.branchOutChans {
			close(out)
		}
	}()
}

type chanMerger struct {
	mergeInChans []chan data.JSON
	mergeWait    sync.WaitGroup
}

func (dp *ProcessorWrapper) mergeIn() {
	// Start a merge goroutine for each input channel.
	mergeData := func(c chan data.JSON) {
		for d := range c {
			dp.inputChan <- d
		}
		dp.mergeWait.Done()
	}
	dp.mergeWait.Add(len(dp.mergeInChans))
	for _, in := range dp.mergeInChans {
		go mergeData(in)
	}

	go func() {
		dp.mergeWait.Wait()
		close(dp.inputChan)
	}()
}

// Do takes a DataProcessor instance and returns the ProcessorWrapper
// type that will wrap it for internal ratchet processing. The details
// of the ProcessorWrapper type are abstracted away from the
// implementing end-user code. The "Do" function is named
// succinctly to provide a nicer syntax when creating a PipelineLayout.
// See the ratchet package documentation for code examples of creating
// a new branching pipeline layout.
func Do(p processors.DataProcessor) *ProcessorWrapper {
	dp := ProcessorWrapper{DataProcessor: p}
	dp.outputChan = make(chan data.JSON)
	dp.inputChan = make(chan data.JSON)

	if isConcurrent(p) {
		dp.concurrency = p.(ConcurrentDataProcessor).Concurrency()
		dp.workThrottle = make(chan workSignal, dp.concurrency)
		dp.workList = list.New()
		dp.doneChan = make(chan bool)
		dp.inputClosed = false
	}

	return &dp
}

// Outputs should be called to specify which DataProcessor instances the current
// processor should send it's output to. See the ratchet package
// documentation for code examples and diagrams.
func (dp *ProcessorWrapper) Outputs(p ...processors.DataProcessor) *ProcessorWrapper {
	dp.outputs = p
	return dp
}

// pass through String output to the DataProcessor
func (dp *ProcessorWrapper) String() string {
	return fmt.Sprintf("%v", dp.DataProcessor)
}
