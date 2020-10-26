package json

import (
	"bytes"
	json2 "encoding/json"
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
	m2 := "youth is wasted on the young"
	c2, _ := json2.Marshal(m2)
	m3 := Message{"Alice", "youth is wasted on the young", 1294706395881547000}
	c3, _ := json2.Marshal(m3)
	m4 := Product{Name:"Peter", ProductID:2, Number:10, Price:12.2, IsOnSale:true}
	c4, _ := json2.Marshal(m4)
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


func Test_addrMarshalerEncoder(t *testing.T) {
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_addrTextMarshalerEncoder(t *testing.T) {
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_arrayEncoder_encode(t *testing.T) {
	type fields struct {
		elemEnc encoderFunc
	}
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ae := arrayEncoder{
				elemEnc: tt.fields.elemEnc,
			}
		})
	}
}

func Test_asciiEqualFold(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := asciiEqualFold(tt.args.s, tt.args.t); got != tt.want {
				t.Errorf("asciiEqualFold() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_boolEncoder(t *testing.T) {
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_byIndex_Len(t *testing.T) {
	tests := []struct {
		name string
		x    byIndex
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.x.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_byIndex_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		x    byIndex
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.x.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_byIndex_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		x    byIndex
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_cachedTypeFields(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
		want structFields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cachedTypeFields(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cachedTypeFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compact(t *testing.T) {
	type args struct {
		dst    *bytes.Buffer
		src    []byte
		escape bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := compact(tt.args.dst, tt.args.src, tt.args.escape); (err != nil) != tt.wantErr {
				t.Errorf("compact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_condAddrEncoder_encode(t *testing.T) {
	type fields struct {
		canAddrEnc encoderFunc
		elseEnc    encoderFunc
	}
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce := condAddrEncoder{
				canAddrEnc: tt.fields.canAddrEnc,
				elseEnc:    tt.fields.elseEnc,
			}
		})
	}
}

func Test_distribute(t *testing.T) {
	type args struct {
		t         reflect.Type
		allowAddr bool
	}
	tests := []struct {
		name string
		args args
		want encoderFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := distribute(tt.args.t, tt.args.allowAddr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("distribute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dominantField(t *testing.T) {
	type args struct {
		fields []field
	}
	tests := []struct {
		name  string
		args  args
		want  field
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := dominantField(tt.args.fields)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dominantField() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("dominantField() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_encodeByteSlice(t *testing.T) {
	type args struct {
		e   *encodeState
		v   reflect.Value
		in2 encodeConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_encodeState_marshal(t *testing.T) {
	type fields struct {
		Buffer   bytes.Buffer
		scratch  [64]byte
		ptrLevel uint
		ptrSeen  map[interface{}]struct{}
	}
	type args struct {
		v    interface{}
		opts encodeConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &encodeState{
				Buffer:   tt.fields.Buffer,
				scratch:  tt.fields.scratch,
				ptrLevel: tt.fields.ptrLevel,
				ptrSeen:  tt.fields.ptrSeen,
			}
			if err := e.marshal(tt.args.v, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("marshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_encodeState_reflectValue(t *testing.T) {
	type fields struct {
		Buffer   bytes.Buffer
		scratch  [64]byte
		ptrLevel uint
		ptrSeen  map[interface{}]struct{}
	}
	type args struct {
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &encodeState{
				Buffer:   tt.fields.Buffer,
				scratch:  tt.fields.scratch,
				ptrLevel: tt.fields.ptrLevel,
				ptrSeen:  tt.fields.ptrSeen,
			}
		})
	}
}

func Test_encodeState_string(t *testing.T) {
	type fields struct {
		Buffer   bytes.Buffer
		scratch  [64]byte
		ptrLevel uint
		ptrSeen  map[interface{}]struct{}
	}
	type args struct {
		s          string
		escapeHTML bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &encodeState{
				Buffer:   tt.fields.Buffer,
				scratch:  tt.fields.scratch,
				ptrLevel: tt.fields.ptrLevel,
				ptrSeen:  tt.fields.ptrSeen,
			}
		})
	}
}

func Test_encodeState_stringBytes(t *testing.T) {
	type fields struct {
		Buffer   bytes.Buffer
		scratch  [64]byte
		ptrLevel uint
		ptrSeen  map[interface{}]struct{}
	}
	type args struct {
		s          []byte
		escapeHTML bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &encodeState{
				Buffer:   tt.fields.Buffer,
				scratch:  tt.fields.scratch,
				ptrLevel: tt.fields.ptrLevel,
				ptrSeen:  tt.fields.ptrSeen,
			}
		})
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equalFoldRight(tt.args.s, tt.args.t); got != tt.want {
				t.Errorf("equalFoldRight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_floatEncoder_encode(t *testing.T) {
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name string
		bits floatEncoder
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_foldFunc(t *testing.T) {
	type args struct {
		s []byte
	}
	tests := []struct {
		name string
		args args
		want func(s, t []byte) bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := foldFunc(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("foldFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_intEncoder(t *testing.T) {
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_interfaceEncoder(t *testing.T) {
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_invalidValueEncoder(t *testing.T) {
	type args struct {
		e   *encodeState
		v   reflect.Value
		in2 encodeConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_isEmptyValue(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidTag(tt.args.s); got != tt.want {
				t.Errorf("isValidTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mapEncoder_encode(t *testing.T) {
	type fields struct {
		elemEnc encoderFunc
	}
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			me := mapEncoder{
				elemEnc: tt.fields.elemEnc,
			}
		})
	}
}

func Test_marshalerEncoder(t *testing.T) {
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_newArrayEncoder(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
		want encoderFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newArrayEncoder(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newArrayEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newCondAddrEncoder(t *testing.T) {
	type args struct {
		canAddrEnc encoderFunc
		elseEnc    encoderFunc
	}
	tests := []struct {
		name string
		args args
		want encoderFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCondAddrEncoder(tt.args.canAddrEnc, tt.args.elseEnc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCondAddrEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newMapEncoder(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
		want encoderFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newMapEncoder(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newMapEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newPtrEncoder(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
		want encoderFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newPtrEncoder(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPtrEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newSliceEncoder(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
		want encoderFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSliceEncoder(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSliceEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newStructEncoder(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
		want encoderFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newStructEncoder(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newStructEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTag(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseTag(tt.args.tag)
			if got != tt.want {
				t.Errorf("parseTag() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseTag() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_ptrEncoder_encode(t *testing.T) {
	type fields struct {
		elemEnc encoderFunc
	}
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pe := ptrEncoder{
				elemEnc: tt.fields.elemEnc,
			}
		})
	}
}

func Test_reflectByString_resolve(t *testing.T) {
	type fields struct {
		v reflect.Value
		s string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &reflectByString{
				v: tt.fields.v,
				s: tt.fields.s,
			}
			if err := w.resolve(); (err != nil) != tt.wantErr {
				t.Errorf("resolve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_scanner_eof(t *testing.T) {
	type fields struct {
		step       func(*scanner, byte) int
		endTop     bool
		parseState []int
		err        error
		bytes      int64
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &scanner{
				step:       tt.fields.step,
				endTop:     tt.fields.endTop,
				parseState: tt.fields.parseState,
				err:        tt.fields.err,
				bytes:      tt.fields.bytes,
			}
			if got := s.eof(); got != tt.want {
				t.Errorf("eof() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scanner_error(t *testing.T) {
	type fields struct {
		step       func(*scanner, byte) int
		endTop     bool
		parseState []int
		err        error
		bytes      int64
	}
	type args struct {
		c       byte
		context string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &scanner{
				step:       tt.fields.step,
				endTop:     tt.fields.endTop,
				parseState: tt.fields.parseState,
				err:        tt.fields.err,
				bytes:      tt.fields.bytes,
			}
			if got := s.error(tt.args.c, tt.args.context); got != tt.want {
				t.Errorf("error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scanner_popParseState(t *testing.T) {
	type fields struct {
		step       func(*scanner, byte) int
		endTop     bool
		parseState []int
		err        error
		bytes      int64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &scanner{
				step:       tt.fields.step,
				endTop:     tt.fields.endTop,
				parseState: tt.fields.parseState,
				err:        tt.fields.err,
				bytes:      tt.fields.bytes,
			}
		})
	}
}

func Test_scanner_pushParseState(t *testing.T) {
	type fields struct {
		step       func(*scanner, byte) int
		endTop     bool
		parseState []int
		err        error
		bytes      int64
	}
	type args struct {
		c             byte
		newParseState int
		successState  int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &scanner{
				step:       tt.fields.step,
				endTop:     tt.fields.endTop,
				parseState: tt.fields.parseState,
				err:        tt.fields.err,
				bytes:      tt.fields.bytes,
			}
			if got := s.pushParseState(tt.args.c, tt.args.newParseState, tt.args.successState); got != tt.want {
				t.Errorf("pushParseState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_simple0(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := simple0(tt.args.s, tt.args.t); got != tt.want {
				t.Errorf("simple0() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sliceEncoder_encode(t *testing.T) {
	type fields struct {
		arrayEnc encoderFunc
	}
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			se := sliceEncoder{
				arrayEnc: tt.fields.arrayEnc,
			}
		})
	}
}

func Test_state0(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := state0(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("state0() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_state1(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := state1(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("state1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateBeginString(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateBeginString(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateBeginString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateBeginStringOrEmpty(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateBeginStringOrEmpty(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateBeginStringOrEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateBeginValue(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateBeginValue(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateBeginValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateBeginValueOrEmpty(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateBeginValueOrEmpty(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateBeginValueOrEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateDot(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateDot(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateDot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateDot0(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateDot0(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateDot0() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateE(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateE(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateE() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateE0(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateE0(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateE0() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateESign(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateESign(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateESign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateEndTop(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateEndTop(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateEndTop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateEndValue(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateEndValue(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateEndValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateError(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateError(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateF(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateF(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateFa(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateFa(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateFa() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateFal(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateFal(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateFal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateFals(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateFals(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateFals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateInString(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateInString(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateInString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateInStringEsc(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateInStringEsc(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateInStringEsc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateInStringEscU(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateInStringEscU(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateInStringEscU() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateInStringEscU1(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateInStringEscU1(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateInStringEscU1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateInStringEscU12(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateInStringEscU12(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateInStringEscU12() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateInStringEscU123(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateInStringEscU123(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateInStringEscU123() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateN(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateN(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateNeg(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateNeg(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateNeg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateNu(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateNu(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateNu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateNul(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateNul(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateNul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateT(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateT(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateTr(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateTr(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateTr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stateTru(t *testing.T) {
	type args struct {
		s *scanner
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stateTru(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("stateTru() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringEncoder(t *testing.T) {
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_structEncoder_encode(t *testing.T) {
	type fields struct {
		fields structFields
	}
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			se := structEncoder{
				fields: tt.fields.fields,
			}
		})
	}
}

func Test_textMarshalerEncoder(t *testing.T) {
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_typeByIndex(t *testing.T) {
	type args struct {
		t     reflect.Type
		index []int
	}
	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := typeByIndex(tt.args.t, tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("typeByIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typeEncoder(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
		want encoderFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := typeEncoder(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("typeEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typeFields(t *testing.T) {
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name string
		args args
		want structFields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := typeFields(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("typeFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uintEncoder(t *testing.T) {
	type args struct {
		e    *encodeState
		v    reflect.Value
		opts encodeConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_unsupportedTypeEncoder(t *testing.T) {
	type args struct {
		e   *encodeState
		v   reflect.Value
		in2 encodeConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_valueEncoder(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want encoderFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := valueEncoder(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("valueEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

