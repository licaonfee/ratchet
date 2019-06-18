package processors_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/licaonfee/ratchet/data"
	"github.com/licaonfee/ratchet/processors"
)

func TestIoWriter_ProcessData(t *testing.T) {

	type args struct {
		d          []data.JSON
		outputChan chan data.JSON
		killChan   chan error
	}
	sharedBuff := &bytes.Buffer{}
	tests := []struct {
		name string
		opts []processors.Option
		args args
		want []byte
	}{
		{name: "Write raw",
			opts: []processors.Option{processors.WithWriter(sharedBuff)},
			args: args{[]data.JSON{data.JSON("123"), data.JSON("456"), data.JSON("789")}, make(chan data.JSON, 5), make(chan error, 1)},
			want: []byte("123456789"),
		},
		{name: "Add new line",
			opts: []processors.Option{processors.WithWriter(sharedBuff),
				processors.AddNewline()},
			args: args{[]data.JSON{data.JSON("123"), data.JSON("456"), data.JSON("789")}, make(chan data.JSON, 5), make(chan error, 1)},
			want: []byte("123\n456\n789\n")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sharedBuff.Reset()
			w, err := processors.NewIOWriter(tt.opts...)
			if err != nil {
				t.Error(err)
			}
			for _, d := range tt.args.d {
				w.ProcessData(d, tt.args.outputChan, tt.args.killChan)
			}

			if !reflect.DeepEqual(sharedBuff.Bytes(), tt.want) {
				t.Errorf("ProcessData() = %v, want = %v", sharedBuff.Bytes(), tt.want)
			}
		})
	}
}
