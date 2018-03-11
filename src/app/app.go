package main

import (
	"clwrapper"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"model"
	"net/http"
	"os"
	"strconv"
)

var PORT int = 3000
var CLIENT_BUFFER = 5
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

func DataHandeler(w http.ResponseWriter, r *http.Request) {

	//start point
	start := 0
	if keys := r.URL.Query(); len(keys) > 0 {
		if val, exists := keys["itter"]; exists {
			if v, err := strconv.Atoi(val[0]); err == nil {
				start = CLIENT_BUFFER * v
			}
		}
	}
	end := start + CLIENT_BUFFER
	if end > m.Numpts {
		end = m.Numpts
	}
	sz := end - start
	if sz < 0 {
		sz = 0
	}

	var data []byte
	// lets define first byte length of folloup points
	data = make([]byte, 4+(12*sz))
	binary.LittleEndian.PutUint32(data[0:4], uint32(sz))
	for i := 0; i < sz; i++ {
		off := i * 12
		binary.LittleEndian.PutUint32(data[off+4:off+8], math.Float32bits(m.Pts[start+i].GetX()))
		binary.LittleEndian.PutUint32(data[off+8:off+12], math.Float32bits(m.Pts[start+i].GetY()))
		binary.LittleEndian.PutUint32(data[off+12:off+16], math.Float32bits(m.Pts[start+i].GetZ()))

		fmt.Println(m.Pts[i].GetX(), m.Pts[i].GetY(), m.Pts[i].GetZ())
	}

	set_resp_headers(w)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	binary.Write(w, binary.LittleEndian, data)

	fmt.Printf("No pts: %d\n", sz)
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
	http.HandleFunc("/points/", DataHandeler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
