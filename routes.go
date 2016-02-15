package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func Index(c *echo.Context) error {
	return c.String(http.StatusOK, "HTTP working over a socket!")
}

// func Apps(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	apps, err := resources.GetApps()
// 	if err != nil {
// 		fmt.Fprintln(w, "Error requesting apps list: ", err)
// 		return
// 	}

// 	data, err := json.MarshalIndent(apps, "", "  ")
// 	if err != nil {
// 		fmt.Fprintln(w, "Error serializing apps list: ", err)
// 		return
// 	}

// 	w.Write(data)
// }

// func App(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	name := params.ByName("name")

// 	app, err := resources.GetApp(name)
// 	if err != nil {
// 		fmt.Fprintf(w, "Error requesting app %s: %s\n", name, err)
// 		return
// 	}

// 	data, err := json.MarshalIndent(app, "", "  ")
// 	if err != nil {
// 		fmt.Fprintf(w, "Error serializing app %s: %s\n", name, err)
// 		return
// 	}

// 	w.Write(data)
// }
