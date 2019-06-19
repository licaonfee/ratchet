package processors_test

import (
	"reflect"
	"testing"

	"github.com/licaonfee/ratchet/data"
	"github.com/licaonfee/ratchet/processors"
)

func TestRegexpMatcher_ProcessData(t *testing.T) {
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
		{name: "Match Only numbers",
			opts: []processors.Option{processors.WithPattern("[0-9]+")},
			args: args{[]data.JSON{data.JSON("123"), data.JSON("foo"), data.JSON("456")},
				make(chan data.JSON, 2), make(chan error, 1)},
			want: []data.JSON{data.JSON("123"), data.JSON("456")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := processors.NewRegexpMatcher(tt.opts...)
			if err != nil {
				t.Error(err)
			}
			for _, d := range tt.args.d {
				r.ProcessData(d, tt.args.outputChan, tt.args.killChan)
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
