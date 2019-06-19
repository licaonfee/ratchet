package processors_test

import (
	"reflect"
	"testing"

	"github.com/licaonfee/ratchet/data"
	"github.com/licaonfee/ratchet/processors"
)

func unchanged(d data.JSON) data.JSON {
	return d
}

func TestFuncTransformer_ProcessData(t *testing.T) {

	type args struct {
		d          []data.JSON
		outputChan chan data.JSON
		killChan   chan error
	}
	tests := []struct {
		name string
		opts []processors.Option
		args args
		want []data.JSON
	}{
		{name: "return unchanged",
			opts: []processors.Option{processors.WithFunc(unchanged)},
			args: args{[]data.JSON{data.JSON("foo"), data.JSON("bar"), data.JSON("baz")},
				make(chan data.JSON, 3), make(chan error, 1)},
			want: []data.JSON{data.JSON("foo"), data.JSON("bar"), data.JSON("baz")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := processors.NewFuncTransformer(tt.opts...)
			if err != nil {
				t.Error(err)
			}
			for _, d := range tt.args.d {
				p.ProcessData(d, tt.args.outputChan, tt.args.killChan)
			}
			got := make([]data.JSON, 0)
			for i := 0; i < len(tt.want); i++ {
				q := <-tt.args.outputChan
				got = append(got, q)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessData() = %v , want = %v", got, tt.want)
			}
		})
	}
}
