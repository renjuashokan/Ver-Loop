package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello world")
	db := datastore{}
	db.initDatastore()
	db.insertTest()
}
