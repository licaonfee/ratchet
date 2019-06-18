package processors

import (
	"errors"
	"fmt"
	"io"

	"github.com/licaonfee/ratchet/data"
	"github.com/licaonfee/ratchet/logger"
	"github.com/licaonfee/ratchet/util"
)

// IoWriter wraps any io.Writer object.
// It can be used to write data out to a File, os.Stdout, or
// any other task that can be supported via io.Writer.
type IoWriter struct {
	Writer     io.Writer
	AddNewline bool
}

// NewIoWriter returns a new IoWriter wrapping the given io.Writer object
func NewIoWriter(writer io.Writer) *IoWriter {
	return &IoWriter{Writer: writer, AddNewline: false}
}

func NewIOWriter(opts ...Option) (DataProcessor, error) {
	w := &IoWriter{}
	for _, o := range opts {
		if err := o(w); err != nil {
			return nil, err
		}
	}
	return w, nil
}

func toIOWriter(p DataProcessor) (*IoWriter, error) {
	w, ok := p.(*IoWriter)
	if !ok {
		return nil, errors.New("must be an IoWriter")
	}
	return w, nil
}

func WithWriter(w io.Writer) Option {
	return func(p DataProcessor) error {
		d, err := toIOWriter(p)
		if err != nil {
			return err
		}
		d.Writer = w
		return nil
	}
}

func AddNewline() Option {
	return func(p DataProcessor) error {
		d, err := toIOWriter(p)
		if err != nil {
			return err
		}
		d.AddNewline = true
		return nil
	}
}

// ProcessData writes the data
func (w *IoWriter) ProcessData(d data.JSON, outputChan chan data.JSON, killChan chan error) {
	var bytesWritten int
	var err error
	if w.AddNewline {
		bytesWritten, err = fmt.Fprintln(w.Writer, string(d))
	} else {
		bytesWritten, err = w.Writer.Write(d)
	}
	util.KillPipelineIfErr(err, killChan)
	logger.Debug("IoWriter:", bytesWritten, "bytes written")
}

// Finish - see interface for documentation.
func (w *IoWriter) Finish(outputChan chan data.JSON, killChan chan error) {
}

func (w *IoWriter) String() string {
	return "IoWriter"
}
