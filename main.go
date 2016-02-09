package main

import (
	"log"
	"fmt"
	"os"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("Index!")
	fmt.Fprintln(w, "HTTP working over a socket!")
}

func Test(w http.ResponseWriter, r * http.Request, params httprouter.Params) {
	fmt.Fprintln(w, params.ByName("string"))
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/test/:string", Test)

	if _, err := os.Stat("/tmp/dokku-api/api.sock"); os.IsNotExist(err) {
		os.MkdirAll("/tmp/dokku-api", 0777)
	}

	log.Fatal(ListenAndServeUnix("/tmp/dokku-api/api.sock", os.FileMode(0666), router))
}
