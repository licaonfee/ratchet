package processors

import (
	"errors"
	"regexp"

	"github.com/licaonfee/ratchet/data"
	"github.com/licaonfee/ratchet/logger"
	"github.com/licaonfee/ratchet/util"
)

// RegexpMatcher checks if incoming data matches the given Regexp, and sends
// it on to the next stage only if it matches.
// It is using regexp.Match under the covers: https://golang.org/pkg/regexp/#Match
type RegexpMatcher struct {
	pattern string //TODO: this should be *regexp.Regexp
	// Set to true to log each match attempt (logger must be in debug mode).
	DebugLog bool
}

// NewRegexpMatcher returns a new RegexpMatcher initialized
// with the given pattern to match.
// func NewRegexpMatchers(pattern string) *RegexpMatcher {
// 	return &RegexpMatcher{pattern, false}
// }

//NewRegexpMatcher returns a new RegexpMatcher or nil if configuraion is invalid
func NewRegexpMatcher(opts ...Option) (DataProcessor, error) {
	m := &RegexpMatcher{}
	for _, o := range opts {
		o(m)
	}
	return m, nil
}

func WithPattern(pattern string) Option {
	return func(d DataProcessor) error {
		m, ok := d.(*RegexpMatcher)
		if !ok {
			return errors.New("must be a RegexpMatcher")
		}
		if _, err := regexp.Compile(pattern); err != nil {
			return err
		}
		m.pattern = pattern
		return nil
	}
}

// ProcessData sends the data it receives to the outputChan only if it matches the supplied regex
func (r *RegexpMatcher) ProcessData(d data.JSON, outputChan chan data.JSON, killChan chan error) {
	matches, err := regexp.Match(r.pattern, d)
	util.KillPipelineIfErr(err, killChan)
	if r.DebugLog {
		logger.Debug("RegexpMatcher: checking if", string(d), "matches pattern", r.pattern, ". MATCH=", matches)
	}
	if matches {
		outputChan <- d
	}
}

// Finish - see interface for documentation.
func (r *RegexpMatcher) Finish(outputChan chan data.JSON, killChan chan error) {
}

func (r *RegexpMatcher) String() string {
	return "RegexpMatcher"
}
