package processors

import "github.com/licaonfee/ratchet/data"

// DataProcessor is the interface that should be implemented to perform data-related
// tasks within a Pipeline. DataProcessors are responsible for receiving, processing,
// and then sending data on to the next stage of processing.
type DataProcessor interface {
	// ProcessData will be called for each data sent from the previous stage.
	// ProcessData is called with a data.JSON instance, which is the data being received,
	// an outputChan, which is the channel to send data to, and a killChan,
	// which is a channel to send unexpected errors to (halting execution of the Pipeline).
	ProcessData(d data.JSON, outputChan chan data.JSON, killChan chan error)

	// Finish will be called after the previous stage has finished sending data,
	// and no more data will be received by this DataProcessor. Often times
	// Finish can be an empty function implementation, but sometimes it is
	// necessary to perform final data processing.
	Finish(outputChan chan data.JSON, killChan chan error)
}

//Builder is an interface for DataProcessor creation, this allow to test full Object creation
type Builder func(...Option) (DataProcessor, error)

//Option represent an initialization action in DataProcessor
type Option func(DataProcessor) error

var avaliable = []DataProcessor{&IoReader{}}
