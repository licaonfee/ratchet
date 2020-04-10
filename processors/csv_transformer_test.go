package processors_test

import (
	"reflect"
	"testing"

	"github.com/licaonfee/ratchet/data"
	"github.com/licaonfee/ratchet/processors"
)

func TestCSVTransformer_ProcessData(t *testing.T) {

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
		// FIXME: #9
		// {name: "JSON Array",
		// 	opts: []processors.Option{},
		// 	args: args{[]data.JSON{data.JSON("[1,2,3]"), data.JSON("[4,5,6]"), data.JSON("[7,8,9]")},
		// 		make(chan data.JSON, 3), make(chan error, 1)},
		// 	want: []data.JSON{data.JSON("1,2,3"), data.JSON("4,5,6"), data.JSON("7,8,9")}},

		{name: "JSON Object numeric data",
			opts: []processors.Option{},
			args: args{[]data.JSON{
				data.JSON(`{"foo": 1, "bar": 2, "baz": 3}`),
				data.JSON(`{"foo": 4, "bar": 5, "baz": 6}`),
				data.JSON(`{"foo": 7, "bar": 8, "baz": 9}`),
			},
				make(chan data.JSON, 3), make(chan error, 1)},
			want: []data.JSON{
				data.JSON(`"bar","baz","foo"` + "\n" + `"2.0000","3.0000","1.0000"` + "\n"),
				data.JSON(`"5.0000","6.0000","4.0000"` + "\n"),
				data.JSON(`"8.0000","9.0000","7.0000"` + "\n"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := processors.NewCSVTransformer(tt.opts...)
			if err != nil {
				t.Error(err)
			}
			for _, d := range tt.args.d {
				p.ProcessData(d, tt.args.outputChan, tt.args.killChan)
			}
			got := make([]data.JSON, 0)
			for i := 0; i < len(tt.want); i++ {
				q := <-tt.args.outputChan
				t.Logf("GET: %s", q)
				got = append(got, q)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessData() = %s , want = %s", got, tt.want)
			}
		})
	}
}

func BenchmarkCSVTransformer_ProcessData(b *testing.B) {
	b.StopTimer()
	row := data.JSON(`{"foo": 1, "bar": 2, "baz": 3}`)
	w, err := processors.NewCSVTransformer()
	if err != nil {
		b.FailNow()
	}
	outC := make(chan data.JSON)
	killC := make(chan error)
	go func() {
		for {
			select {
			case <-outC:
			case <-killC:
				b.FailNow()
			}
		}
	}()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w.ProcessData(row, outC, killC)
	}

}
