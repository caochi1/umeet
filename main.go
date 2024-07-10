package main

import "fmt"

func main() {
	type ao struct{
		name string
	}
	a := []ao{ao{"asd"}, ao{"dsa"}}
	b := []ao{}
	for _, v := range a {
		b = append(b, v)
	}
	fmt.Println(b)
	// models.Init(false)
	// cache.Init()
	// router.Router()

}
