package main

import (
	"bytes"
	json2 "encoding/json"
	"fmt"
	"github.com/Winszheng/json"
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

func main()  {

	m1 := Message{"Alice", "youth is wasted on the young", 1294706395881547000}
	m2 :="youth is wasted on the young"
	b1, _ := json.JsonMarshal(m1)
	b2, _ := json.JsonMarshal(m2)

	var cc Message
	json2.Unmarshal(b1, &cc)
	fmt.Printf("1.经过序列化后再次转换出的结构体:\n%v\n", cc)

	c1, _ := json2.Marshal(m1)
	c2, _ := json2.Marshal(m2)

	fmt.Println("2.因为两者序列化的效果一样，所以bytes.Compare的结果是(后面的每一行输入都是比较字节切片，输出0表示效果相同，符合要求)：")
	fmt.Println("3: ",bytes.Compare(b1, c1))
	fmt.Println("4: ",bytes.Compare(b2, c2))

	m3 := Product{Name:"Peter", ProductID:2, Number:10, Price:12.2, IsOnSale:true}
	b3, _ := json.JsonMarshal(m3)
	c3, _ := json2.Marshal(m3)
	fmt.Println("4: ",bytes.Compare(b3, c3))

	data := []byte(`<div>这个是html标签</div>`)
	var buffer1, buffer2 bytes.Buffer
	json.HTMLEscape(&buffer1,data)

	json2.HTMLEscape(&buffer2, data)
	fmt.Println("5: ",bytes.Compare(buffer1.Bytes(), buffer2.Bytes()))

}