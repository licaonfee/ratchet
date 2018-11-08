package processors

import (
	"github.com/licaonfee/ratchet/logger"

	"github.com/licaonfee/ratchet/util"
	"github.com/nats-io/go-nats-streaming"

	"github.com/licaonfee/ratchet/data"
)

// StanWriter extract messages from a nats-streaming server
type StanWriter struct {
	conn             stan.Conn
	queue            string
	BatchSize        int
	ConcurrencyLevel int
}

func (s *StanWriter) lostConnection(conn stan.Conn, err error) {
	logger.Info(err)
}

//NewStanWriter return a new Nats Streaming client associated with queue
func NewStanWriter(clusterID, clientID, queue, url string) *StanWriter {
	r := StanWriter{queue: queue}

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url),
		stan.SetConnectionLostHandler(r.lostConnection))
	if err != nil {
		panic(err)
	}
	r.conn = sc
	return &r
}

// ProcessData - see interface for documentation.
func (s *StanWriter) ProcessData(d data.JSON, outputChan chan data.JSON, killChan chan error) {

	err := s.conn.Publish(s.queue, d)
	util.KillPipelineIfErr(err, killChan)

}

// Finish - see interface for documentation.
func (s *StanWriter) Finish(outputChan chan data.JSON, killChan chan error) {

}

func (s *StanWriter) String() string {
	return "StanWriter"
}

// Concurrency defers to ConcurrentDataProcessor
func (s *StanWriter) Concurrency() int {
	return s.ConcurrencyLevel
}
