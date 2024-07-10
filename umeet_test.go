package main

import (
	"fmt"
	"testing"
)

func TestSomething(t *testing.T) {
	a := make([][]int, 2)
	for i := 0; i < 2; i++ {
		a[i] = append(a[i], i)
	}
	fmt.Printf("a: %v\n", a)
	 
	// fmt.Println(unsafe.Sizeof(l))
	// type ooo struct {
	// 	Name string
	// 	Age  int
	// }
	// o := ooo{
	// 	Name: "CC",
	// 	Age:  12,
	// }
	// var d ooo
	// jsbyte, _ := json.Marshal(o)
	// var e []byte
	// e = []byte(string(jsbyte))
	// json.Unmarshal(e, &d)
	// fmt.Println(string(jsbyte))
	// fmt.Println(d)

	// time.AfterFunc(time.Second, func() {
	// 	fmt.Println("1")
	// })

}

func Benchmark_TestSomething(B *testing.B) {

}
