package ratchet_test

import (
	"reflect"
	"testing"

	"github.com/licaonfee/ratchet"
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
		processor ratchet.DataProcessor
	}
	dummy := &dummyDataProcessor{q: 0}
	tests := []struct {
		name string
		args args
		want *ratchet.ProcessorWrapper
	}{
		{"Inmutable DataProcessor", args{processor: dummy}, &ratchet.ProcessorWrapper{DataProcessor: dummy}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ratchet.Do(tt.args.processor); !reflect.DeepEqual(got.DataProcessor, tt.want.DataProcessor) {
				t.Errorf("Do() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
