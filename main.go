package main

import (
	"log"
	"fmt"
	"os"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "HTTP working over a socket!")
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)

	log.Fatal(ListenAndServeUNIX("/tmp/dokku-api.sock", os.FileMode(0666), router))
}
