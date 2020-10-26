package json

import (
	"reflect"
	"testing"
)

type Message struct {
	Name string
	Body string
	Time int64
}


// 这就是表格驱动的单元测试！1
func TestJsonMarshal(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{"1",args{"Alice"}, []byte{34, 65, 108, 105, 99, 101, 34},false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JsonMarshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonMarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonMarshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}