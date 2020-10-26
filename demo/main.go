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
	b1, _ := json.JsonMarshal(m1)

	var cc Message
	json2.Unmarshal(b1, &cc)
	fmt.Printf("1.经过序列化后再次转换出的结构体:\n%v\n", cc)

	c1, _ := json2.Marshal(m1)

	fmt.Println("\n2.因为两者序列化的效果一样，所以bytes.Compare的结果是：",bytes.Compare(b1, c1))

}

