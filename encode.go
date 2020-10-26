package json

import (
	"bytes"
	"encoding"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// 接口，用于把对象序列化为有效json
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

var (
	marshalerType     = reflect.TypeOf((*Marshaler)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

// encodeState把json encode成bytes.Buffer
type encodeState struct {
	bytes.Buffer
	scratch [64]byte // 避免可能导致堆栈溢出的循环
	ptrLevel uint
	ptrSeen  map[interface{}]struct{}
}

// encOpts实现编码过程中的配置
type encodeConfig struct {
	quoted bool //如果为true，int float编码之后会加上双引号
	escapeHTML bool //默认为true，会转义 <、>、& 符号，如果不想转义，可以自定义Encoder对象
}

type encoderFunc func(e *encodeState, v reflect.Value, opts encodeConfig)
var encoderCache sync.Map // map[reflect.Type]encoderFunc
type floatEncoder int   // 位数

// 编码浮点数字段
func (bits floatEncoder) encode(e *encodeState, v reflect.Value, opts encodeConfig) {
	f := v.Float()  // 获取浮点数的值到f

	if math.IsInf(f, 0) || math.IsNaN(f) {
		panic("unsupportedValueType")   // 数据类型不支持
	}

	b := e.scratch[:0]  // 空切片
	a := math.Abs(f)  // 获取浮点数的绝对值
	fmt := byte('f')

	if a!=0 && bits == 64 && (a < 1e-6 || a >= 1e21) || bits == 32 && (float32(a) < 1e-6 || float32(a) >= 1e21) {
		fmt = 'e'
	}

	b = strconv.AppendFloat(b, f, fmt, -1, int(bits))
	if fmt == 'e' {
		if n := len(b); n >= 4 && b[n-4] == 'e' && b[n-3] == '-' && b[n-2] == '0' {
			b[n-2] = b[n-1]
			b = b[:n-1]
		}
	}

	if opts.quoted {
		e.WriteByte('"')
	}
	e.Write(b)
	if opts.quoted {
		e.WriteByte('"')
	}
}

var float32Encoder = (floatEncoder(32)).encode
var float64Encoder = (floatEncoder(64)).encode

// distribute递归地为所有类型找到对应的编码方法(函数)
func distribute(t reflect.Type, allowAddr bool) encoderFunc {
	// 如果当前值是值类型、可以取地址并且值类型的指针类型实现了
	if t.Kind() != reflect.Ptr && allowAddr && reflect.PtrTo(t).Implements(marshalerType) {
		// newCondAddrEncoder 获取条件编码器
		// addrMarshalerEncoder如果失败，distribute重新选择新的编码器
		return newCondAddrEncoder(addrMarshalerEncoder, distribute(t, false))
	}
	// 如果该值类型实现了Marshaler接口，直接使用marshalerEncoder对该值进行序列化
	if t.Implements(marshalerType) {
		return marshalerEncoder
	}
	if t.Kind() != reflect.Ptr && allowAddr && reflect.PtrTo(t).Implements(textMarshalerType) {
		return newCondAddrEncoder(addrTextMarshalerEncoder, distribute(t, false))
	}
	if t.Implements(textMarshalerType) {
		return textMarshalerEncoder
	}

	switch t.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncoder
	case reflect.Float32:
		return float32Encoder
	case reflect.Float64:
		return float64Encoder
	case reflect.String:
		return stringEncoder
	case reflect.Interface:
		return interfaceEncoder
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	case reflect.Ptr:
		return newPtrEncoder(t)
	default:
		return unsupportedTypeEncoder
	}
}

func unsupportedTypeEncoder(e *encodeState, v reflect.Value, _ encodeConfig) {
	panic("unsupportedType")
}

type ptrEncoder struct {
	elemEnc encoderFunc
}

const startDetectingCyclesAfter = 1000

func (pe ptrEncoder) encode(e *encodeState, v reflect.Value, opts encodeConfig) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}

	if e.ptrLevel++; e.ptrLevel > startDetectingCyclesAfter {   // 查询
		ptr := v.Interface()
		if _, ok := e.ptrSeen[ptr]; ok {
			panic("unsupportedValue")
		}
		e.ptrSeen[ptr] = struct{}{}
		defer delete(e.ptrSeen, ptr)
	}
	pe.elemEnc(e, v.Elem(), opts)
	e.ptrLevel--
}

func newPtrEncoder(t reflect.Type) encoderFunc {
	enc := ptrEncoder{typeEncoder(t.Elem())}
	return enc.encode
}

type arrayEncoder struct {
	elemEnc encoderFunc
}

func (ae arrayEncoder) encode(e *encodeState, v reflect.Value, opts encodeConfig) {
	e.WriteByte('[')
	n := v.Len()
	for i := 0; i < n; i++ {
		if i > 0 {
			e.WriteByte(',')
		}
		ae.elemEnc(e, v.Index(i), opts)
	}
	e.WriteByte(']')
}

func newArrayEncoder(t reflect.Type) encoderFunc {
	enc := arrayEncoder{typeEncoder(t.Elem())}
	return enc.encode
}

// 编码字节切片
func encodeByteSlice(e *encodeState, v reflect.Value, _ encodeConfig) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}
	s := v.Bytes()
	e.WriteByte('"')
	encodedLen := base64.StdEncoding.EncodedLen(len(s))
	if encodedLen <= len(e.scratch) {
		dst := e.scratch[:encodedLen]
		base64.StdEncoding.Encode(dst, s)
		e.Write(dst)
	} else if encodedLen <= 1024 {
		dst := make([]byte, encodedLen)
		base64.StdEncoding.Encode(dst, s)
		e.Write(dst)
	} else {
		enc := base64.NewEncoder(base64.StdEncoding, e)
		enc.Write(s)
		enc.Close()
	}
	e.WriteByte('"')
}

// 保证值不为nil
type sliceEncoder struct {
	arrayEnc encoderFunc
}

func (se sliceEncoder) encode(e *encodeState, v reflect.Value, opts encodeConfig) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}
	if e.ptrLevel++; e.ptrLevel > startDetectingCyclesAfter {
		ptr := struct {
			ptr uintptr
			len int
		}{v.Pointer(), v.Len()}
		if _, ok := e.ptrSeen[ptr]; ok {
			panic("unsupportedValue")
		}
		e.ptrSeen[ptr] = struct{}{}
		defer delete(e.ptrSeen, ptr)
	}
	se.arrayEnc(e, v, opts)
	e.ptrLevel--
}

func newSliceEncoder(t reflect.Type) encoderFunc {
	if t.Elem().Kind() == reflect.Uint8 {
		p := reflect.PtrTo(t.Elem())
		if !p.Implements(marshalerType) && !p.Implements(textMarshalerType) {
			return encodeByteSlice
		}
	}
	enc := sliceEncoder{newArrayEncoder(t)}
	return enc.encode
}

type mapEncoder struct {
	elemEnc encoderFunc
}

func (me mapEncoder) encode(e *encodeState, v reflect.Value, opts encodeConfig) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}
	if e.ptrLevel++; e.ptrLevel > startDetectingCyclesAfter {   // 嵌套
		ptr := v.Pointer()
		if _, ok := e.ptrSeen[ptr]; ok {
			panic("unsupportedValue")
		}
		e.ptrSeen[ptr] = struct{}{}
		defer delete(e.ptrSeen, ptr)
	}
	e.WriteByte('{')

	keys := v.MapKeys()
	sv := make([]reflectByString, len(keys))
	for i, v := range keys {
		sv[i].v = v
		if err := sv[i].resolve(); err != nil {
			panic(fmt.Errorf("json: encoding error for type %q: %q", v.Type().String(), err.Error()))
		}
	}
	sort.Slice(sv, func(i, j int) bool { return sv[i].s < sv[j].s })

	for i, kv := range sv {
		if i > 0 {
			e.WriteByte(',')
		}
		e.string(kv.s, opts.escapeHTML)
		e.WriteByte(':')
		me.elemEnc(e, v.MapIndex(kv.v), opts)
	}
	e.WriteByte('}')
	e.ptrLevel--
}

var htmlSafeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      false,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      false,
	'=':      true,
	'>':      false,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

// 十六進制数
var hex = "0123456789abcdef"

// NOTE: keep in sync with stringBytes below.
func (e *encodeState) string(s string, escapeHTML bool) {
	e.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] || (!escapeHTML && safeSet[b]) {
				i++
				continue
			}
			if start < i {
				e.WriteString(s[start:i])
			}
			e.WriteByte('\\')
			switch b {
			case '\\', '"':
				e.WriteByte(b)
			case '\n':
				e.WriteByte('n')
			case '\r':
				e.WriteByte('r')
			case '\t':
				e.WriteByte('t')
			default:
				e.WriteString(`u00`)
				e.WriteByte(hex[b>>4])
				e.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				e.WriteString(s[start:i])
			}
			e.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				e.WriteString(s[start:i])
			}
			e.WriteString(`\u202`)
			e.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		e.WriteString(s[start:])
	}
	e.WriteByte('"')
}

type reflectByString struct {
	v reflect.Value
	s string
}

func (w *reflectByString) resolve() error {
	if w.v.Kind() == reflect.String {
		w.s = w.v.String()
		return nil
	}
	if tm, ok := w.v.Interface().(encoding.TextMarshaler); ok {
		if w.v.Kind() == reflect.Ptr && w.v.IsNil() {
			return nil
		}
		buf, err := tm.MarshalText()
		w.s = string(buf)
		return err
	}
	switch w.v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		w.s = strconv.FormatInt(w.v.Int(), 10)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		w.s = strconv.FormatUint(w.v.Uint(), 10)
		return nil
	}
	panic("unexpected map key type")
}

func newMapEncoder(t reflect.Type) encoderFunc {
	switch t.Key().Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	default:
		if !t.Key().Implements(textMarshalerType) {
			return unsupportedTypeEncoder
		}
	}
	me := mapEncoder{typeEncoder(t.Elem())}
	return me.encode
}

type structEncoder struct {
	fields structFields
}

type field struct {
	name      string
	equalFold func(s, t []byte) bool // bytes.EqualFold or equivalent
	nameEscHTML string // `"` + HTMLEscape(name) + `":`
	tag       bool
	index     []int
	typ       reflect.Type
	omitEmpty bool
	quoted    bool
	encoder encoderFunc
}

type structFields struct {
	list      []field
	nameIndex map[string]int
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// 结构体的json编码器
func newStructEncoder(t reflect.Type) encoderFunc {
	se := structEncoder{fields: cachedTypeFields(t)}
	return se.encode
}

// encode遍历结构体中的所有字段，在写入字段名后，调用字段对应类型的编码方法将字段对应的json写入缓冲区
func (se structEncoder) encode(e *encodeState, v reflect.Value, opts encodeConfig) {
	next := byte('{')
	for i := range se.fields.list {
		f := &se.fields.list[i]

		fv := v
		for _, i := range f.index {
			if fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					continue
				}
				fv = fv.Elem()
			}
			fv = fv.Field(i)
		}

		if f.omitEmpty && isEmptyValue(fv) {
			continue
		}
		e.WriteByte(next)
		next = ','
		if opts.escapeHTML {
			e.WriteString(f.nameEscHTML)
		} else {
			e.WriteString(`"` + f.name + `":`)
		}
		opts.quoted = f.quoted
		f.encoder(e, fv, opts)
	}
	if next == '{' {
		e.WriteString("{}")
	} else {
		e.WriteByte('}')
	}
}

var fieldCache sync.Map // map[reflect.Type]structFields

// 依据名字和逗号解析标签
func parseTag(tag string) (string, string) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], string(tag[idx+1:])
	}
	return tag, string("")
}

// Contains 逗号分隔
func Contains(s string, optionName string) bool {
	for s != "" {
		var next string
		index := strings.Index(s, ",")  // 寻找逗号下标
		if index >= 0 {
			s, next = s[:index], s[index+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

// 判断是否有效标签
func isValidTag(s string) bool {
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:;<=>?@[]^_{|}~ ", c):
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}
	return true
}

const (
	caseMask     = ^byte(0x20) // Mask to ignore case in ASCII.
	kelvin       = '\u212a'
	smallLongEss = '\u017f'
)

// 四种情况
func foldFunc(s []byte) func(s, t []byte) bool {
	nonLetter := false
	special := false // special letter
	for _, b := range s {
		if b >= utf8.RuneSelf {
			return bytes.EqualFold
		}
		upper := b & caseMask
		if upper < 'A' || upper > 'Z' {
			nonLetter = true
		} else if upper == 'K' || upper == 'S' {
			special = true
		}
	}
	if special {
		return equalFoldRight
	}
	if nonLetter {
		return asciiEqualFold
	}
	return simple0
}

func simple0(s, t []byte) bool {
	if len(s) != len(t) {
		return false
	}
	for i, b := range s {
		if b&caseMask != t[i]&caseMask {
			return false
		}
	}
	return true
}

func asciiEqualFold(s, t []byte) bool {
	if len(s) != len(t) {
		return false
	}
	for i, sb := range s {
		tb := t[i]
		if sb == tb {
			continue
		}
		if ('a' <= sb && sb <= 'z') || ('A' <= sb && sb <= 'Z') {
			if sb&caseMask != tb&caseMask {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func equalFoldRight(s, t []byte) bool {
	for _, sb := range s {
		if len(t) == 0 {
			return false
		}
		tb := t[0]
		if tb < utf8.RuneSelf {
			if sb != tb {
				sbUpper := sb & caseMask
				if 'A' <= sbUpper && sbUpper <= 'Z' {
					if sbUpper != tb&caseMask {
						return false
					}
				} else {
					return false
				}
			}
			t = t[1:]
			continue
		}

		tr, size := utf8.DecodeRune(t)
		switch sb {
		case 's', 'S':
			if tr != smallLongEss {
				return false
			}
		case 'k', 'K':
			if tr != kelvin {
				return false
			}
		default:
			return false
		}
		t = t[size:]

	}
	if len(t) > 0 {
		return false
	}
	return true
}

// HTMLEscape 实现标准HTML转义
// 将json编码的src中的<、>、&、U+2028 和U+2029字符替换为\u003c、\u003e、\u0026、\u2028、\u2029 转义字符串，
// 以便json编码可以安全的嵌入HTML的<script>标签里
// 网络浏览器不支持在<script>标签中使用标准HTML转义
func HTMLEscape(dst *bytes.Buffer, src []byte) {
	start := 0
	for i, c := range src {
		if c == '<' || c == '>' || c == '&' {
			if start < i {
				dst.Write(src[start:i])
			}
			dst.WriteString(`\u00`)
			dst.WriteByte(hex[c>>4])
			dst.WriteByte(hex[c&0xF])
			start = i + 1
		}

		if c == 0xE2 && i+2 < len(src) && src[i+1] == 0x80 && src[i+2]&^1 == 0xA8 {
			if start < i {
				dst.Write(src[start:i])
			}
			dst.WriteString(`\u202`)
			dst.WriteByte(hex[src[i+2]&0xF])
			start = i + 3
		}
	}
	if start < len(src) {
		dst.Write(src[start:])
	}
}

type byIndex []field

func (x byIndex) Less(i, j int) bool {
	for k, xik := range x[i].index {
		if k >= len(x[j].index) {
			return false
		}
		if xik != x[j].index[k] {
			return xik < x[j].index[k]
		}
	}
	return len(x[i].index) < len(x[j].index)
}

func (x byIndex) Len() int { return len(x) }

func (x byIndex) Swap(i, j int) { x[i], x[j] = x[j], x[i] }


func typeFields(t reflect.Type) structFields {
	cur := []field{}
	next := []field{{typ: t}}
	var count, nextCount map[reflect.Type]int

	visited := map[reflect.Type]bool{}  //已访问

	var fields []field

	var nameEscBuf bytes.Buffer //缓冲区

	for len(next) > 0 {
		cur, next = next, cur[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}
		for _, f := range cur {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			for i := 0; i < f.typ.NumField(); i++ { // Scan f.typ for fields to include.
				sf := f.typ.Field(i)
				isUnexported := sf.PkgPath != ""
				if sf.Anonymous {
					t := sf.Type
					if t.Kind() == reflect.Ptr {
						t = t.Elem()
					}
					if isUnexported && t.Kind() != reflect.Struct {
						continue
					}
				} else if isUnexported {
					continue
				}
				tag := sf.Tag.Get("json")
				if tag == "-" {
					continue
				}
				name, opts := parseTag(tag)
				if !isValidTag(name) {
					name = ""
				}
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Ptr {
					ft = ft.Elem()
				}

				quoted := false
				if Contains(opts, "string") {
					switch ft.Kind() {
					case reflect.Bool,
						reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
						reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
						reflect.Float32, reflect.Float64,
						reflect.String:
						quoted = true
					}
				}

				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					tagged := name != ""
					if name == "" {
						name = sf.Name
					}
					field := field{
						name:      name,
						tag:       tagged,
						index:     index,
						typ:       ft,
						omitEmpty: Contains(opts, "omitempty"),
						quoted:    quoted,
					}

					field.equalFold = foldFunc([]byte(field.name))
					nameEscBuf.Reset()
					nameEscBuf.WriteString(`"`)
					HTMLEscape(&nameEscBuf, []byte(field.name))
					nameEscBuf.WriteString(`":`)
					field.nameEscHTML = nameEscBuf.String()

					fields = append(fields, field)
					if count[f.typ] > 1 {
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, field{name: ft.Name(), index: index, typ: ft})
				}
			}
		}
	}

	sort.Slice(fields, func(i, j int) bool {
		x := fields
		// 按名排序
		if x[i].name != x[j].name {
			return x[i].name < x[j].name
		}
		if len(x[i].index) != len(x[j].index) {
			return len(x[i].index) < len(x[j].index)
		}
		if x[i].tag != x[j].tag {
			return x[i].tag
		}
		return byIndex(x).Less(i, j)
	})

	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance {
		fi := fields[i]
		name := fi.name
		for advance = 1; i+advance < len(fields); advance++ {
			fj := fields[i+advance]
			if fj.name != name {
				break
			}
		}
		if advance == 1 {
			out = append(out, fi)
			continue
		}
		dominant, ok := dominantField(fields[i : i+advance])
		if ok {
			out = append(out, dominant)
		}
	}

	fields = out
	sort.Sort(byIndex(fields))  // 排序

	for i := range fields {
		f := &fields[i]
		f.encoder = typeEncoder(typeByIndex(t, f.index))
	}
	nameIndex := make(map[string]int, len(fields))
	for i, field := range fields {
		nameIndex[field.name] = i
	}
	return structFields{fields, nameIndex}
}

func typeByIndex(t reflect.Type, index []int) reflect.Type {
	for _, i := range index {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		t = t.Field(i).Type
	}
	return t
}

func dominantField(fields []field) (field, bool) {
	if len(fields) > 1 && len(fields[0].index) == len(fields[1].index) && fields[0].tag == fields[1].tag {
		return field{}, false
	}
	return fields[0], true
}

// 使用cache避免repeat
func cachedTypeFields(t reflect.Type) structFields {
	if f, ok := fieldCache.Load(t); ok {
		return f.(structFields)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return f.(structFields)
}

func interfaceEncoder(e *encodeState, v reflect.Value, opts encodeConfig) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}
	e.reflectValue(v.Elem(), opts)
}

type Number string
var numberType = reflect.TypeOf(Number(""))

// 合法性
func isValidNumber(s string) bool {
	if s == "" {
		return false
	}
	// Optional -
	if s[0] == '-' {
		s = s[1:]
		if s == "" {
			return false
		}
	}
	// Digits
	switch {
	default:
		return false

	case s[0] == '0':
		s = s[1:]

	case '0' < s[0] && s[0] <= '9':
		s = s[1:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	if len(s) >= 2 && s[0] == '.' && '0' <= s[1] && s[1] <= '9' {
		s = s[2:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	if len(s) >= 2 && (s[0] == 'e' || s[0] == 'E') {
		s = s[1:]
		if s[0] == '+' || s[0] == '-' {
			s = s[1:]
			if s == "" {
				return false
			}
		}
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	return s == ""
}

func (e *encodeState) stringBytes(s []byte, escapeHTML bool) {
	e.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] || (!escapeHTML && safeSet[b]) {
				i++
				continue
			}
			if start < i {
				e.Write(s[start:i])
			}
			e.WriteByte('\\')
			switch b {
			case '\\', '"':
				e.WriteByte(b)
			case '\n':
				e.WriteByte('n')
			case '\r':
				e.WriteByte('r')
			case '\t':
				e.WriteByte('t')
			default:
				e.WriteString(`u00`)
				e.WriteByte(hex[b>>4])
				e.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRune(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				e.Write(s[start:i])
			}
			e.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}

		if c == '\u2028' || c == '\u2029' {
			if start < i {
				e.Write(s[start:i])
			}
			e.WriteString(`\u202`)
			e.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		e.Write(s[start:])
	}
	e.WriteByte('"')
}

func stringEncoder(e *encodeState, v reflect.Value, opts encodeConfig) {
	if v.Type() == numberType {
		numStr := v.String()

		if numStr == "" {
			numStr = "0"
		}
		if !isValidNumber(numStr) {
			panic(fmt.Errorf("json: invalid number literal %q", numStr))
		}
		if opts.quoted {
			e.WriteByte('"')
		}
		e.WriteString(numStr)
		if opts.quoted {
			e.WriteByte('"')
		}
		return
	}
	if opts.quoted {
		e2 := &encodeState{}
		e2.string(v.String(), opts.escapeHTML)
		e.stringBytes(e2.Bytes(), false)
	} else {
		e.string(v.String(), opts.escapeHTML)
	}
}

func intEncoder(e *encodeState, v reflect.Value, opts encodeConfig) {
	b := strconv.AppendInt(e.scratch[:0], v.Int(), 10)
	if opts.quoted {
		e.WriteByte('"')
	}
	e.Write(b)
	if opts.quoted {
		e.WriteByte('"')
	}
}

func uintEncoder(e *encodeState, v reflect.Value, opts encodeConfig) {
	b := strconv.AppendUint(e.scratch[:0], v.Uint(), 10)
	if opts.quoted {
		e.WriteByte('"')
	}
	e.Write(b)
	if opts.quoted {
		e.WriteByte('"')
	}
}

func textMarshalerEncoder(e *encodeState, v reflect.Value, opts encodeConfig) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		e.WriteString("null")
		return
	}
	m, ok := v.Interface().(encoding.TextMarshaler)
	if !ok {
		e.WriteString("null")
		return
	}
	b, err := m.MarshalText()
	if err != nil {
		panic("MarshalerError")
	}
	e.stringBytes(b, opts.escapeHTML)
}

func addrTextMarshalerEncoder(e *encodeState, v reflect.Value, opts encodeConfig) {
	va := v.Addr()
	if va.IsNil() {
		e.WriteString("null")
		return
	}
	m := va.Interface().(encoding.TextMarshaler)
	b, err := m.MarshalText()
	if err != nil {
		panic("MarshalerError")
	}
	e.stringBytes(b, opts.escapeHTML)
}

func marshalerEncoder(e *encodeState, v reflect.Value, opts encodeConfig) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		e.WriteString("null")
		return
	}
	m, ok := v.Interface().(Marshaler)
	if !ok {
		e.WriteString("null")
		return
	}
	b, err := m.MarshalJSON()
	if err == nil {
		err = compact(&e.Buffer, b, opts.escapeHTML)    // 把json copy到缓冲区，并检查有效性
	}
	if err != nil {
		panic("MarshalerError")
	}
}

func addrMarshalerEncoder(e *encodeState, v reflect.Value, opts encodeConfig) {
	va := v.Addr()
	if va.IsNil() {
		e.WriteString("null")
		return
	}
	m := va.Interface().(Marshaler)
	b, err := m.MarshalJSON()
	if err == nil {
		err = compact(&e.Buffer, b, opts.escapeHTML)    // 把json copy到缓冲区，并检查有效性
	}
	if err != nil {
		panic("MarshalerError")
	}
}

// scanner是一个状态机
// 把byte切片分割成一个个完整的值
type scanner struct {
	step func(*scanner, byte) int
	endTop bool
	parseState []int
	err error
	bytes int64
}

var scannerPool = sync.Pool{
	New: func() interface{} {
		return &scanner{}
	},
}

func isSpace(c byte) bool {
	return c <= ' ' && (c == ' ' || c == '\t' || c == '\r' || c == '\n')
}

// 状态转移函数的返回值 - 这里是复制json包的～
const (
	scanContinue     = iota // uninteresting byte
	scanBeginLiteral        // end implied by next result != scanContinue
	scanBeginObject         // begin object
	scanObjectKey           // just finished object key (string)
	scanObjectValue         // just finished non-last object value
	scanEndObject           // end object (implies scanObjectValue if possible)
	scanBeginArray          // begin array
	scanArrayValue          // just finished array value
	scanEndArray            // end array (implies scanArrayValue if possible)
	scanSkipSpace           // space byte; can skip; known to be last "continue" result

	scanEnd   // top-level value ended *before* this byte; known to be first "stop" result
	scanError // 错误
)

func stateEndTop(s *scanner, c byte) int {
	if !isSpace(c) {
		s.error(c, "after top-level value")
	}
	return scanEnd
}

func stateError(s *scanner, c byte) int {
	return scanError
}

func (s *scanner) error(c byte, context string) int {
	s.step = stateError
	s.err = errors.New("invalid character")
	return scanError
}

const (
	parseKey   = iota // 解析关键字
	parseValue        // 解析value
	parseArray         // 数组
)

func stateInStringEscU123(s *scanner, c byte) int {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInString
		return scanContinue
	}
	
	return s.error(c, "in \\u hexadecimal character escape")    // numbers
}

func stateInStringEscU12(s *scanner, c byte) int {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU123
		return scanContinue
	}
	return s.error(c, "in \\u hexadecimal character escape")
}

func stateInStringEscU1(s *scanner, c byte) int {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU12
		return scanContinue
	}
	return s.error(c, "in \\u hexadecimal character escape")
}

func stateInStringEscU(s *scanner, c byte) int {
	if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
		s.step = stateInStringEscU1
		return scanContinue
	}
	return s.error(c, "in \\u hexadecimal character escape")
}

func stateInStringEsc(s *scanner, c byte) int {
	switch c {
	case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
		s.step = stateInString
		return scanContinue
	case 'u':
		s.step = stateInStringEscU
		return scanContinue
	}
	return s.error(c, "in string escape code")
}

func stateInString(s *scanner, c byte) int {
	if c == '"' {
		s.step = stateEndValue
		return scanContinue
	}
	if c == '\\' {
		s.step = stateInStringEsc
		return scanContinue
	}
	if c < 0x20 {
		return s.error(c, "in string literal")
	}
	return scanContinue
}

func stateBeginString(s *scanner, c byte) int {
	if isSpace(c) {
		return scanSkipSpace
	}
	if c == '"' {
		s.step = stateInString
		return scanBeginLiteral
	}
	return s.error(c, "looking for beginning of object key string")
}

func (s *scanner) popParseState() {
	n := len(s.parseState) - 1
	s.parseState = s.parseState[0:n]
	if n == 0 {
		s.step = stateEndTop
		s.endTop = true
	} else {
		s.step = stateEndValue
	}
}

func stateEndValue(s *scanner, c byte) int {
	n := len(s.parseState)
	if n == 0 {
		// Completed top-level before the current byte.
		s.step = stateEndTop
		s.endTop = true
		return stateEndTop(s, c)
	}
	if isSpace(c) {
		s.step = stateEndValue
		return scanSkipSpace
	}
	ps := s.parseState[n-1]
	switch ps {
	case parseKey:
		if c == ':' {
			s.parseState[n-1] = parseValue
			s.step = stateBeginValue
			return scanObjectKey
		}
		return s.error(c, "after object key")
	case parseValue:
		if c == ',' {
			s.parseState[n-1] = parseKey
			s.step = stateBeginString
			return scanObjectValue
		}
		if c == '}' {
			s.popParseState()
			return scanEndObject
		}
		return s.error(c, "after object key:value pair")
	case parseArray:
		if c == ',' {
			s.step = stateBeginValue
			return scanArrayValue
		}
		if c == ']' {
			s.popParseState()
			return scanEndArray
		}
		return s.error(c, "after array element")
	}
	return s.error(c, "")
}

func stateBeginValueOrEmpty(s *scanner, c byte) int {
	if isSpace(c) {
		return scanSkipSpace
	}
	if c == ']' {
		return stateEndValue(s, c)
	}
	return stateBeginValue(s, c)
}

const maxNestDepth = 10000   // 最大嵌套深度

func (s *scanner) pushParseState(c byte, newParseState int, successState int) int {
	s.parseState = append(s.parseState, newParseState)
	if len(s.parseState) <= maxNestDepth {
		return successState
	}
	return s.error(c, "exceeded max depth")
}

// 状态0
func state0(s *scanner, c byte) int {
	if c == '.' {
		s.step = stateDot
		return scanContinue
	}
	if c == 'e' || c == 'E' {
		s.step = stateE
		return scanContinue
	}
	return stateEndValue(s, c)
}

func stateDot(s *scanner, c byte) int {
	if '0' <= c && c <= '9' {
		s.step = stateDot0
		return scanContinue
	}
	return s.error(c, "after decimal point in numeric literal")
}

func stateDot0(s *scanner, c byte) int {
	if '0' <= c && c <= '9' {
		return scanContinue
	}
	if c == 'e' || c == 'E' {
		s.step = stateE
		return scanContinue
	}
	return stateEndValue(s, c)
}

func stateE(s *scanner, c byte) int {
	if c == '+' || c == '-' {
		s.step = stateESign
		return scanContinue
	}
	return stateESign(s, c)
}

func stateE0(s *scanner, c byte) int {
	if '0' <= c && c <= '9' {
		return scanContinue
	}
	return stateEndValue(s, c)
}

func stateESign(s *scanner, c byte) int {
	if '0' <= c && c <= '9' {
		s.step = stateE0
		return scanContinue
	}
	return s.error(c, "in exponent of numeric literal")
}

func state1(s *scanner, c byte) int {
	if '0' <= c && c <= '9' {
		s.step = state1
		return scanContinue
	}
	return state0(s, c)
}

func stateNeg(s *scanner, c byte) int {
	if c == '0' {
		s.step = state0
		return scanContinue
	}
	if '1' <= c && c <= '9' {
		s.step = state1
		return scanContinue
	}
	return s.error(c, "in numeric literal")
}

func stateBeginStringOrEmpty(s *scanner, c byte) int {
	if isSpace(c) {
		return scanSkipSpace
	}
	if c == '}' {
		n := len(s.parseState)
		s.parseState[n-1] = parseValue
		return stateEndValue(s, c)
	}
	return stateBeginString(s, c)
}

func stateBeginValue(s *scanner, c byte) int {
	if isSpace(c) {
		return scanSkipSpace
	}
	switch c {
	case '{':
		s.step = stateBeginStringOrEmpty
		return s.pushParseState(c, parseKey, scanBeginObject)
	case '[':
		s.step = stateBeginValueOrEmpty
		return s.pushParseState(c, parseArray, scanBeginArray)
	case '"':
		s.step = stateInString
		return scanBeginLiteral
	case '-':
		s.step = stateNeg
		return scanBeginLiteral
	case '0': // beginning of 0.123
		s.step = state0
		return scanBeginLiteral
	case 't': // beginning of true
		s.step = stateT
		return scanBeginLiteral
	case 'f': // beginning of false
		s.step = stateF
		return scanBeginLiteral
	case 'n': // beginning of null
		s.step = stateN
		return scanBeginLiteral
	}
	if '1' <= c && c <= '9' { // beginning of 1234.5
		s.step = state1
		return scanBeginLiteral
	}
	return s.error(c, "looking for beginning of value")
}

func stateF(s *scanner, c byte) int {
	if c == 'a' {
		s.step = stateFa
		return scanContinue
	}
	return s.error(c, "in literal false (expecting 'a')")
}

func stateFa(s *scanner, c byte) int {
	if c == 'l' {
		s.step = stateFal
		return scanContinue
	}
	return s.error(c, "in literal false (expecting 'l')")
}

func stateFal(s *scanner, c byte) int {
	if c == 's' {
		s.step = stateFals
		return scanContinue
	}
	return s.error(c, "in literal false (expecting 's')")
}

func stateFals(s *scanner, c byte) int {
	if c == 'e' {
		s.step = stateEndValue
		return scanContinue
	}
	return s.error(c, "in literal false (expecting 'e')")
}

func stateN(s *scanner, c byte) int {
	if c == 'u' {
		s.step = stateNu
		return scanContinue
	}
	return s.error(c, "in literal null (expecting 'u')")
}

func stateNu(s *scanner, c byte) int {
	if c == 'l' {
		s.step = stateNul
		return scanContinue
	}
	return s.error(c, "in literal null (expecting 'l')")
}

func stateNul(s *scanner, c byte) int {
	if c == 'l' {
		s.step = stateEndValue
		return scanContinue
	}
	return s.error(c, "in literal null (expecting 'l')")
}

func stateT(s *scanner, c byte) int {
	if c == 'r' {
		s.step = stateTr
		return scanContinue
	}
	return s.error(c, "in literal true (expecting 'r')")
}

func stateTr(s *scanner, c byte) int {
	if c == 'u' {
		s.step = stateTru
		return scanContinue
	}
	return s.error(c, "in literal true (expecting 'u')")
}

func stateTru(s *scanner, c byte) int {
	if c == 'e' {
		s.step = stateEndValue
		return scanContinue
	}
	return s.error(c, "in literal true (expecting 'e')")
}

func compact(dst *bytes.Buffer, src []byte, escape bool) error {
	origLen := dst.Len()
	scan := scannerPool.Get().(*scanner)
	scan.bytes = 0
	scan.step = stateBeginValue
	scan.parseState = scan.parseState[0:0]
	scan.err = nil
	scan.endTop = false

	defer func() {
		if len(scan.parseState) > 1024 {
			scan.parseState = nil
		}
		scannerPool.Put(scan)
	}()

	start := 0
	for i, c := range src {
		if escape && (c == '<' || c == '>' || c == '&') {
			if start < i {
				dst.Write(src[start:i])
			}
			dst.WriteString(`\u00`)
			dst.WriteByte(hex[c>>4])
			dst.WriteByte(hex[c&0xF])
			start = i + 1
		}
		// Convert U+2028 and U+2029 (E2 80 A8 and E2 80 A9).
		if escape && c == 0xE2 && i+2 < len(src) && src[i+1] == 0x80 && src[i+2]&^1 == 0xA8 {
			if start < i {
				dst.Write(src[start:i])
			}
			dst.WriteString(`\u202`)
			dst.WriteByte(hex[src[i+2]&0xF])
			start = i + 3
		}
		v := scan.step(scan, c)
		if v >= scanSkipSpace {
			if v == scanError {
				break
			}
			if start < i {
				dst.Write(src[start:i])
			}
			start = i + 1
		}
	}
	if scan.eof() == scanError {
		dst.Truncate(origLen)
		return scan.err
	}

	if start < len(src) {
		dst.Write(src[start:])
	}
	return nil
}

// 是否到最后
func (s *scanner) eof() int {
	if s.err != nil {
		return scanError
	}
	if s.endTop {
		return scanEnd
	}
	s.step(s, ' ')
	if s.endTop {
		return scanEnd
	}
	if s.err == nil {
		s.err = errors.New("unexpected end of JSON input")
	}
	return scanError
}

type condAddrEncoder struct {
	canAddrEnc, elseEnc encoderFunc
}

func (ce condAddrEncoder) encode(e *encodeState, v reflect.Value, opts encodeConfig) {
	if v.CanAddr() {
		ce.canAddrEnc(e, v, opts)
	} else {
		ce.elseEnc(e, v, opts)
	}
}

func newCondAddrEncoder(canAddrEnc, elseEnc encoderFunc) encoderFunc {
	enc := condAddrEncoder{canAddrEnc: canAddrEnc, elseEnc: elseEnc}
	return enc.encode
}

// 布尔值的json编码器，用于查看最终json串的赋值操作
func boolEncoder(e *encodeState, v reflect.Value, opts encodeConfig) {
	if opts.quoted {
		e.WriteByte('"')
	}
	if v.Bool() {
		e.WriteString("true")
	} else {
		e.WriteString("false")
	}
	if opts.quoted {
		e.WriteByte('"')
	}
}

func typeEncoder(t reflect.Type) encoderFunc {
	f := distribute(t, true)
	encoderCache.Store(t, f)
	return f
}

func invalidValueEncoder(e *encodeState, v reflect.Value, _ encodeConfig) {
	e.WriteString("null")
}

func valueEncoder(v reflect.Value) encoderFunc {
	if !v.IsValid() {
		return invalidValueEncoder
	}
	return typeEncoder(v.Type())
}

func (e *encodeState) reflectValue(v reflect.Value, opts encodeConfig) {
	valueEncoder(v)(e, v, opts)
}

func (e *encodeState) marshal(v interface{}, opts encodeConfig)(err error){
	defer func() {  // 错误处理
		if r := recover(); r!= nil {
			panic("error")
		}
	}()
	e.reflectValue(reflect.ValueOf(v), opts)
	return nil
}

// JsonMarshal根据传入的数据(如结构体)生成json数据
func JsonMarshal(v interface{}) ([]byte, error){
	e := &encodeState{}
	//将v递归转换为json串并写入结构e的buffer中
	err := e.marshal(v, encodeConfig{escapeHTML: true}) //escapeHTML默认为true，转义<、>、&符号
	if err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}