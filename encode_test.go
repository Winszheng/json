package json

import (
	json2 "encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type Message struct {
	Name string
	Body string
	Time int64
}

type Product struct {
	Name      string  `json:"name"`
	ProductID int64   `json:"-"` // 表示不进行序列化
	Number    int     `json:"number"`
	Price     float64 `json:"price"`
	IsOnSale  bool    `json:"is_on_sale,string"`
}

func TestJsonMarshal(t *testing.T) {
	type args struct {
		v interface{}
	}
	// initial cases
	// case 2
	m2 := "youth is wasted on the young"
	c2, _ := json2.Marshal(m2)

	// case 3
	m3 := Message{"Alice", "youth is wasted on the young", 1294706395881547000}
	c3, _ := json2.Marshal(m3)

	//case 4
	m4 := Product{Name:"Peter", ProductID:2, Number:10, Price:12.2, IsOnSale:true}
	c4, _ := json2.Marshal(m4)

	//case 5
	m5 := 0
	c5, _:= json2.Marshal(m5)

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		//1. normal
		{"normal case 1",args{"A"},[]byte{34, 65, 34}, false},
		///{"normal case 2",args{Message{"Alice", "youth is wasted on the young", 1294706395881547000}},[]byte{2}, false},
		//2. normal
		{"normal case 2", args{m2}, c2, false},
		//3. normal
		{"normal case 3", args{m3}, c3, false},
		//4. edge case
		{"normal case 4", args{m4}, c4, false},
		// 5
		{"normal case 5", args{m5}, c5, false},
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

func BenchmarkEasy(b *testing.B) {
	num:=10
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%d",num)
	}
}

func BenchmarkJsonMarshal(b *testing.B) {
	m3 := Message{"Alice", "youth is wasted on the young", 1294706395881547000}
	for i := 0; i < b.N; i++ {

		JsonMarshal(m3)
	}
}

func Test_equalFoldRight(t *testing.T) {
	type args struct {
		s []byte
		t []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"case 1", args{[]byte{1}, []byte{2}}, false},

		{"case 2", args{[]byte{38, 90, 102, 30, 40}, []byte{38, 90, 102, 30, 42}}, false},

		{"case 3", args{[]byte{1}, []byte{1}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equalFoldRight(tt.args.s, tt.args.t); got != tt.want {
				t.Errorf("equalFoldRight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isEmptyValue(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	var a1 int
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"case1", args{reflect.ValueOf(a1)}, true},

		{"case2", args{reflect.ValueOf("")}, true},

		{"case3", args{reflect.ValueOf(nil)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEmptyValue(tt.args.v); got != tt.want {
				t.Errorf("isEmptyValue() = %v, want %v", got, tt.want)
			}
		})
	}
}


func Test_isSpace(t *testing.T) {
	type args struct {
		c byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"case1", args{'c'}, false},

		{"case1", args{' '}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSpace(tt.args.c); got != tt.want {
				t.Errorf("isSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}


func Test_isValidNumber(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"case1", args{"ss"}, false},

		{"case2", args{"8"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidNumber(tt.args.s); got != tt.want {
				t.Errorf("isValidNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}



func Test_isValidTag(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"case1", args{"四月是你的谎言"},true},

		{"case1", args{"%*%……￥……E……F*FJBHOUGOUUBTR&%￥%#……KGIFIYC"},false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidTag(tt.args.s); got != tt.want {
				t.Errorf("isValidTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
