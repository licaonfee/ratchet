package processors

import (
	"github.com/licaonfee/ratchet/data"
	"github.com/licaonfee/ratchet/util"
)

// CSVTransformer converts data.JSON objects into a CSV string object
// and sends it on to the next stage. In use-cases where
// you simply want to write to a CSV file, use CSVWriter instead.
//
// CSVTransformer is for more complex use-cases where you need to
// generate CSV data and perhaps send it to multiple output stages.
type CSVTransformer struct {
	Parameters util.CSVParameters
}

// NewCSVTransformer returns a new CSVTransformer wrapping the given io.Writer object
// func NewCSVTransformer() *CSVTransformer {
// 	return &CSVTransformer{
// 		Parameters: util.CSVParameters{
// 			Writer:        util.NewCSVWriter(),
// 			WriteHeader:   true,
// 			HeaderWritten: false,
// 			SendUpstream:  true,
// 		},
// 	}
// }
func NewCSVTransformer(opts ...Option) (DataProcessor, error) {
	t := &CSVTransformer{}
	t.Parameters = util.CSVParameters{Writer: util.NewCSVWriter(), WriteHeader: true, SendUpstream: true}
	for _, o := range opts {
		if err := o(t); err != nil {
			return nil, err
		}
	}
	return t, nil
}

// ProcessData defers to util.CSVProcess
func (w *CSVTransformer) ProcessData(d data.JSON, outputChan chan data.JSON, killChan chan error) {
	util.CSVProcess(&w.Parameters, d, outputChan, killChan)
}

// Finish - see interface for documentation.
func (w *CSVTransformer) Finish(outputChan chan data.JSON, killChan chan error) {
}

func (w *CSVTransformer) String() string {
	return "CSVTransformer"
}
