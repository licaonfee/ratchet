package ratchet

import (
	"reflect"
	"testing"

	"github.com/licaonfee/ratchet/data"
)

type dummyDataProcessor struct {
	q int
}

func (d *dummyDataProcessor) ProcessData(j data.JSON, o chan data.JSON, k chan error) {

}

func (d *dummyDataProcessor) Finish(o chan data.JSON, k chan error) {

}
func TestDo(t *testing.T) {
	type args struct {
		processor DataProcessor
	}
	dummy := &dummyDataProcessor{q: 0}
	tests := []struct {
		name string
		args args
		want *ProcessorWrapper
	}{
		{"Inmutable DataProcessor", args{processor: dummy}, &ProcessorWrapper{DataProcessor: dummy}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Do(tt.args.processor); !reflect.DeepEqual(got.DataProcessor, tt.want.DataProcessor) {
				t.Errorf("Do() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
