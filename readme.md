# 一、项目结构

生成的中文api文档：https://godoc.org/github.com/Winszheng/json

# 二、使用案例

# 三、关于测试

### 1 单元测试

不必要对文件中的所有函数和结构体做单元测试(因为有些函数或者结构体只有寥寥几行)。

针对程序包中需要测试的函数，编写表格驱动的单元测试，测试文件为`encode_test.go`.

#### 1.1 实行单元测试的函数

#### 1.2 具体说明

以func JsonMarshal(v interface{}) ([]byte, error)为例，说明如何进行单元测试。

##### 原函数

```
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
```

##### 单元测试的内容

```
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
```

#### 查看覆盖率



### 基准测试

# 四、设计说明

