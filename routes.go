package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"

	"github.com/nikelmwann/dokku-api/dokku"
	"github.com/nikelmwann/dokku-api/models"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("Index!")
	fmt.Fprintln(w, "HTTP working over a socket!")
}

func Test(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintln(w, params.ByName("string"))
}

func AppIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	args := []string{"apps"}

	output, err := dokku.Exec(args...)
	if err != nil {
		fmt.Fprintln(w, "Error requesting apps list: ", err)
		return
	}

	// Skip first line of output which is simply "=====> My Apps"
	<-output

	var apps models.Apps
	// Remaining output is the list of apps; one app per line
	for line := range output {
		app, err := models.GetApp(line)
		if err != nil {
			fmt.Fprintln(w, "Error requesting apps list: ", err)
			return
		}

		apps = append(apps, app)
	}

	json.NewEncoder(w).Encode(apps)
}

func AppShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")

	app, err := models.GetApp(name)
	if err != nil {
		fmt.Fprintf(w, "Error requesting app %s: %s\n", name, err)
		return
	}

	json.NewEncoder(w).Encode(app)
}
