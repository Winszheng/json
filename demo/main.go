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
	//1.
	m1 := Message{"Alice", "youth is wasted on the young", 1294706395881547000}
	m2 :="youth is wasted on the young"
	b1, _ := json.JsonMarshal(m1)
	b2, _ := json.JsonMarshal(m2)

	//bbb, _ := json.JsonMarshal(m)
	//fmt.Printf("%v\n, %v\n", b, bbb)
	c1, _ := json2.Marshal(m1)
	c2, _ := json2.Marshal(m2)
	//m1 := 1
	//b1, _ := json.JsonMarshal(m1)
	//fmt.Println("b1: ", string(b1))
	fmt.Println(bytes.Compare(b1, c1), bytes.Compare(b2, c2))
	//fmt.Println(reflect.TypeOf(b1))

	// 2.----
	m3 := Product{Name:"Peter", ProductID:2, Number:10, Price:12.2, IsOnSale:true}
	b3, _ := json.JsonMarshal(m3)
	c3, _ := json2.Marshal(m3)
	fmt.Println(bytes.Compare(b3, c3))

	//3.
	data := []byte(`<div>这个是html标签</div>`)
	var buffer1, buffer2 bytes.Buffer
	json.HTMLEscape(&buffer1,data)
	//fmt.Println(buffer1.String())    //在这里直接compare两个字节流？？也可以吧
	json2.HTMLEscape(&buffer2, data)
	fmt.Println(bytes.Compare(buffer1.Bytes(), buffer2.Bytes()))



}

