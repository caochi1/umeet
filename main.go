package main

import "fmt"

func main() {

	models.Init(false)
	cache.Init()
	router.Router()

}
