package processors

import (
	"bufio"
	"compress/gzip"
	"errors"
	"io"

	"github.com/licaonfee/ratchet/data"
	"github.com/licaonfee/ratchet/util"
)

// IoReader wraps an io.Reader and reads it.
type IoReader struct {
	Reader     io.Reader
	LineByLine bool // defaults to false
	BufferSize int
	Gzipped    bool
}

//NewIoReader create a new IO
func NewIoReader(o ...Option) (DataProcessor, error) {
	r := &IoReader{Reader: nil, LineByLine: false, BufferSize: 1024, Gzipped: false}
	for _, opt := range o {
		if err := opt(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func toIoReader(p DataProcessor) (*IoReader, error) {
	r, ok := p.(*IoReader)
	if !ok {
		return nil, errors.New("must be an IoReader")
	}
	return r, nil
}

func WithReader(rd io.Reader) Option {
	return func(p DataProcessor) error {
		r, err := toIoReader(p)
		if err != nil {
			return err
		}
		r.Reader = rd
		return nil
	}
}

func LineByLine() Option {
	return func(p DataProcessor) error {
		r, err := toIoReader(p)
		if err != nil {
			return err
		}
		r.LineByLine = true
		return nil
	}
}

func WithBufferSize(n int) Option {
	return func(p DataProcessor) error {
		r, err := toIoReader(p)
		if err != nil {
			return err
		}
		r.BufferSize = n
		return nil
	}
}

func Gzipped() Option {
	return func(p DataProcessor) error {
		r, err := toIoReader(p)
		if err != nil {
			return err
		}
		r.Gzipped = true
		return nil
	}
}

// ProcessData overwrites the reader if the content is Gzipped, then defers to ForEachData
func (r *IoReader) ProcessData(d data.JSON, outputChan chan data.JSON, killChan chan error) {
	if r.Gzipped {
		gzReader, err := gzip.NewReader(r.Reader)
		util.KillPipelineIfErr(err, killChan)
		r.Reader = gzReader
	}
	r.ForEachData(killChan, func(d data.JSON) {
		outputChan <- d
	})
}

// Finish - see interface for documentation.
func (r *IoReader) Finish(outputChan chan data.JSON, killChan chan error) {

}

// ForEachData either reads by line or by buffered stream, sending the data
// back to the anonymous func that ultimately shoves it onto the outputChan
func (r *IoReader) ForEachData(killChan chan error, forEach func(d data.JSON)) {
	if r.LineByLine {
		r.scanLines(killChan, forEach)
	} else {
		r.bufferedRead(killChan, forEach)
	}
}

func (r *IoReader) scanLines(killChan chan error, forEach func(d data.JSON)) {
	scanner := bufio.NewScanner(r.Reader)
	for scanner.Scan() {
		forEach(data.JSON(scanner.Bytes()))
	}
	err := scanner.Err()
	util.KillPipelineIfErr(err, killChan)
}

func (r *IoReader) bufferedRead(killChan chan error, forEach func(d data.JSON)) {
	reader := bufio.NewReader(r.Reader)
	d := make([]byte, r.BufferSize)
	for {
		n, err := reader.Read(d)
		if err != nil && err != io.EOF {
			killChan <- err
		}
		if n == 0 {
			break
		}
		forEach(d)
	}
}

func (r *IoReader) String() string {
	return "IoReader"
}
