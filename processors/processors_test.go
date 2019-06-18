package processors_test

import (
	"reflect"
	"testing"

	"github.com/licaonfee/ratchet/data"
	"github.com/licaonfee/ratchet/processors"
)

func callProcess(p processors.DataProcessor, d []data.JSON) (chan data.JSON, chan error) {
	out := make(chan data.JSON, len(d)) //Buffered
	killC := make(chan error, 1)
	go func() {
		for _, t := range d {
			p.ProcessData(t, out, killC)
		}
	}()
	return out, killC
}

func TestProccesor_ProcessData(t *testing.T) {

	tests := []struct {
		name    string
		builder processors.Builder
		opts    []processors.Option
		args    []data.JSON
		want    []data.JSON
		killed  bool
	}{

		//Generic Test cases goes here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := tt.builder(tt.opts...)
			if err != nil {
				t.Errorf("ProcessData() err=%s", err)
			}
			out, kill := callProcess(p, tt.args)
			got := make([]data.JSON, 0)
			for r := range out {
				got = append(got, r)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("p.ProcessData() = %v , want = %v", got, tt.want)
			}
			select {
			case f := <-kill:
				if f != nil {
					t.Errorf("p.ProcessData() fail %s", f)
				}
			default:
			}
		})
	}
}
