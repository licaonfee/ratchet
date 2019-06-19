package processors

import (
	"errors"

	"github.com/licaonfee/ratchet/data"
)

// FuncTransformer executes the given function on each data
// payload, sending the resuling data to the next stage.
//
// While FuncTransformer is useful for simple data transformation, more
// complicated tasks justify building a custom implementation of DataProcessor.
type FuncTransformer struct {
	transform        Transformer
	Name             string // can be set for more useful log output
	ConcurrencyLevel int    // See ConcurrentDataProcessor
}

//Transformer define an especific transformation
type Transformer func(data.JSON) data.JSON

//NewFuncTransformer instantiates a new instance of func transformer
func NewFuncTransformer(opts ...Option) (DataProcessor, error) {
	t := &FuncTransformer{}
	for _, o := range opts {
		if err := o(t); err != nil {
			return nil, err
		}
	}
	return t, nil
}

func WithFunc(f Transformer) Option {
	return func(d DataProcessor) error {
		t, ok := d.(*FuncTransformer)
		if !ok {
			return errors.New("must be a FuncTransformer")
		}
		t.transform = f
		return nil
	}
}

// ProcessData runs the supplied func and sends the returned value to outputChan
func (t *FuncTransformer) ProcessData(d data.JSON, outputChan chan data.JSON, killChan chan error) {
	outputChan <- t.transform(d)
}

// Finish - see interface for documentation.
func (t *FuncTransformer) Finish(outputChan chan data.JSON, killChan chan error) {
}

func (t *FuncTransformer) String() string {
	if t.Name != "" {
		return t.Name
	}
	return "FuncTransformer"
}

// Concurrency defers to ConcurrentDataProcessor
func (t *FuncTransformer) Concurrency() int {
	return t.ConcurrencyLevel
}
