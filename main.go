package main

import (
	"log"
	"os"
	"github.com/julienschmidt/httprouter"

	"github.com/nikelmwann/dokku-api/socket"
)

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/test/:string", Test)
	router.GET("/apps", AppIndex)
	router.GET("/apps/:name", AppShow)

	if _, err := os.Stat("/tmp/dokku-api/api.sock"); os.IsNotExist(err) {
		os.MkdirAll("/tmp/dokku-api", 0777)
	}

	log.Fatal(socket.ListenAndServe("/tmp/dokku-api/api.sock", os.FileMode(0666), router))
}
