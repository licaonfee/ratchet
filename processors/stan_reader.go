package processors

import (
	"github.com/licaonfee/ratchet/util"
	"github.com/nats-io/go-nats-streaming"

	"github.com/licaonfee/ratchet/data"
)

// StanReader extract messages from a nats-streaming server
type StanReader struct {
	conn             stan.Conn
	queue            string
	cdata            chan data.JSON
	kill             chan bool
	BatchSize        int
	ConcurrencyLevel int
}

func (s *StanReader) lostConnection(conn stan.Conn, err error) {
	close(s.cdata)
}

//NewStanReader return a new Nats Streaming client associated with queue
func NewStanReader(clusterID, clientID, queue, url string) *StanReader {
	r := StanReader{queue: queue}
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url),
		stan.SetConnectionLostHandler(r.lostConnection))
	if err != nil {
		panic(err)
	}
	r.conn = sc
	r.cdata = make(chan data.JSON)
	r.kill = make(chan bool, 1)
	return &r
}

// ProcessData - see interface for documentation.
func (s *StanReader) ProcessData(d data.JSON, outputChan chan data.JSON, killChan chan error) {
	s.forEachData(d, killChan, func(d data.JSON) {
		outputChan <- d
	})
}

func (s *StanReader) forEachData(d data.JSON, killChan chan error, forEach func(data.JSON)) {
	sub, err := s.conn.Subscribe(s.queue, func(m *stan.Msg) {
		s.cdata <- m.Data
	}, stan.StartWithLastReceived())
	if err != nil {
		util.KillPipelineIfErr(err, killChan)
		return
	}
	defer sub.Close()
	defer s.conn.Close()
	var derr interface{}

	for d := range s.cdata {
		if err := data.ParseJSONSilent(d, &derr); err != nil {
			util.KillPipelineIfErr(err, killChan)
		} else {
			forEach(d)
		}

	}

}

// Finish - see interface for documentation.
func (s *StanReader) Finish(outputChan chan data.JSON, killChan chan error) {

}

func (s *StanReader) String() string {
	return "StanReader"
}

// Concurrency defers to ConcurrentDataProcessor
func (s *StanReader) Concurrency() int {
	return s.ConcurrencyLevel
}
