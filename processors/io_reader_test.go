package processors_test

import (
	"bytes"
	"compress/gzip"
	"reflect"
	"testing"

	"github.com/licaonfee/ratchet/data"
	"github.com/licaonfee/ratchet/processors"
)

func gzipData(d string) []byte {
	buff := &bytes.Buffer{}
	w := gzip.NewWriter(buff)
	if _, err := w.Write([]byte(d)); err != nil {
		panic(err)
	}
	w.Flush()
	return buff.Bytes()
}
func TestIoReader_ProcessData(t *testing.T) {

	type args struct {
		outputChan chan data.JSON
		killChan   chan error
	}
	tests := []struct {
		name string
		opts []processors.Option
		args args
		want []data.JSON
	}{
		{name: "LineByLine",
			opts: []processors.Option{processors.LineByLine(),
				processors.WithReader(bytes.NewReader([]byte("1\n2\n3")))},
			args: args{make(chan data.JSON, 5), make(chan error, 1)},
			want: []data.JSON{data.JSON("1"), data.JSON("2"), data.JSON("3")}},
		{name: "Raw read with buffer size == data size",
			opts: []processors.Option{processors.WithReader(bytes.NewReader([]byte("12345"))),
				processors.WithBufferSize(5)},
			args: args{make(chan data.JSON, 5), make(chan error, 1)},
			want: []data.JSON{data.JSON("12345")}},
		{name: "Raw read with buffer < data size",
			opts: []processors.Option{processors.WithReader(bytes.NewReader([]byte("1234512345"))),
				processors.WithBufferSize(5)},
			args: args{make(chan data.JSON, 5), make(chan error, 1)},
			want: []data.JSON{data.JSON("12345"), data.JSON("12345")}},
		{name: "Read Gzipped LineByLine",
			opts: []processors.Option{processors.LineByLine(),
				processors.Gzipped(),
				processors.WithReader(bytes.NewReader(gzipData("12345\n12345\n12345")))},
			args: args{make(chan data.JSON, 5), make(chan error, 1)},
			want: []data.JSON{data.JSON("12345"), data.JSON("12345"), data.JSON("12345")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := processors.NewIOReader(tt.opts...)
			if err != nil {
				t.Error(err)
			}
			r.ProcessData(nil, tt.args.outputChan, tt.args.killChan)
			got := make([]data.JSON, 0)
			for i := 0; i < len(tt.want); i++ {
				r := <-tt.args.outputChan
				got = append(got, r)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessData() = %v , want = %v", got, tt.want)
			}
		})
	}
}
