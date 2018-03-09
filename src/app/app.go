package main

import (
	"clwrapper"
	"encoding/json"
	"fmt"
	"model"
	"net/http"
	"os"
)

var PORT int = 3000
var m *model.Model

func InfoHandeler(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		DataSource string
		PointCnt   int
	}{
		os.Args[1],
		0,
	}
	if m.Valid {
		resp.PointCnt = len(m.Pts)
	}
	close_request_json(resp, w)
}

func set_resp_headers(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func close_request_json(vals interface{}, w http.ResponseWriter) {
	set_resp_headers(w)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vals)
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage:")
		fmt.Println("\tlasserver datasource-path")
		os.Exit(0)
	}

	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		fmt.Println("please provide path to valid file")
		os.Exit(1)
	}

	// TODO: check for opencl
	if !clwrapper.HasSupport() {
		fmt.Println("Having open cl support and good performance graphics card is necessary for this application")
		os.Exit(1)
	}

	// TODO: check the structure of the file, in new class
	//  and build structure for ....
	m = model.CreateModel(os.Args[1])

	fmt.Println("Starting web server on port", PORT)
	http.HandleFunc("/info", InfoHandeler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
