package ratchet

import "github.com/licaonfee/ratchet/processors"

// PipelineStage holds one or more DataProcessor instances.
type PipelineStage struct {
	p []*ProcessorWrapper
}

// NewPipelineStage creates a PipelineStage instance given a series
// of dataProcessors. dataProcessor (lower-case d) is a private wrapper around
// an object implementing the public DataProcessor interface. The
// syntax used to create PipelineLayouts abstracts this type
// away from your implementing code. For example:
//
//     layout, err := ratchet.NewPipelineLayout(
//             ratchet.NewPipelineStage(
//                      ratchet.Do(aDataProcessor).Outputs(anotherDataProcessor),
//                      // ...
//             ),
//             // ...
//     )
//
// Notice how the ratchet.Do() and Outputs() functions allow you to insert
// DataProcessor instances into your PipelineStages without having to
// worry about the internal dataProcessor type or how any of the
// channel management works behind the scenes.
//
// See the ratchet package documentation for more code examples.
func NewPipelineStage(p ...*ProcessorWrapper) *PipelineStage {
	return &PipelineStage{p}
}

//TODO: Review if iterate over index is better than foreach
func (s *PipelineStage) hasProcessor(p processors.DataProcessor) bool {
	for i := range s.p {
		if s.p[i].DataProcessor == p {
			return true
		}
	}
	return false
}

//TODO: Same as hasProcessor
func (s *PipelineStage) hasOutput(p processors.DataProcessor) bool {
	for i := range s.p {
		for j := range s.p[i].outputs {
			if s.p[i].outputs[j] == p {
				return true
			}
		}
	}
	return false
}
