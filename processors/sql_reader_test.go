package processors_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/licaonfee/ratchet/data"
	"github.com/licaonfee/ratchet/processors"
)

func TestSQLReader_ProcessData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	mock.ExpectQuery("SELECT row1, row2 FROM test1;").
		WillReturnRows(sqlmock.NewRows([]string{"row1", "row2"}).
			FromCSVString("1,1").FromCSVString("2,2"))

	mock.ExpectQuery("SELECT row1 FROM test2;").
		WillReturnRows(sqlmock.NewRows([]string{"row1"}).
			FromCSVString("1").FromCSVString("2").FromCSVString("3"))

	mock.ExpectQuery("SELECT row1 FROM test3;").
		WillReturnRows(sqlmock.NewRows([]string{"row1"}).
			FromCSVString("1").FromCSVString("2").FromCSVString("3").FromCSVString("4"))
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
		{name: "Simple query: test1",
			opts: []processors.Option{processors.WithQuery("SELECT row1, row2 FROM test1;"),
				processors.WithDB(db)},
			args: args{[]data.JSON{data.JSON("GO")}, make(chan data.JSON, 5), make(chan error, 1)},
			want: []data.JSON{data.JSON(`[{"row1":"1","row2":"1"},{"row1":"2","row2":"2"}]`)}},

		{name: "Query with generator: test2",
			opts: []processors.Option{processors.WithSQLGenerator(
				func(d data.JSON) (string, error) {
					return fmt.Sprintf("SELECT row1 FROM %s;", string(d)), nil
				}),
				processors.WithDB(db)},
			args: args{[]data.JSON{data.JSON("test2")}, make(chan data.JSON, 5), make(chan error, 1)},
			want: []data.JSON{data.JSON(`[{"row1":"1"},{"row1":"2"},{"row1":"3"}]`)}},

		{name: "Simple query multi batch: test3",
			opts: []processors.Option{processors.WithQuery("SELECT row1 FROM test3;"),
				processors.WithBatchSize(1),
				processors.WithDB(db)},
			args: args{[]data.JSON{data.JSON("GO")}, make(chan data.JSON, 5), make(chan error, 1)},
			want: []data.JSON{data.JSON(`[{"row1":"1"}]`), data.JSON(`[{"row1":"2"}]`),
				data.JSON(`[{"row1":"3"}]`), data.JSON(`[{"row1":"4"}]`)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := processors.NewSQLReader(tt.opts...)
			if err != nil {
				t.Error(err)
			}
			for _, d := range tt.args.d {
				r.ProcessData(d, tt.args.outputChan, tt.args.killChan)
			}
			got := make([]data.JSON, 0)
			for i := 0; i < len(tt.want); i++ {
				r := <-tt.args.outputChan

				got = append(got, r)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessData() = %s , want = %s", got, tt.want)
			}

		})
	}
}
