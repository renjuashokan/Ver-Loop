package main

import (
	"fmt"
	"net/http"
)

type datastore struct {
	HTTPClient http.Client
}

func (instance *datastore) initDatastore() {
	fmt.Println("initiating datastore")
}
