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
var CLIENT_BUFFER = 1500
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
	if start > end {
		start = end
	}
	// fmt.Println(start, end)
	// fmt.Println(m.Sizes)

	var data []byte
	data = make([]byte, 4+(16*(end-start)))
	binary.LittleEndian.PutUint32(data[0:4], uint32(end-start))
	
	cpt := 0
	off := 0
	for i, n := range m.Sizes {
		if cpt >= end {
			// fmt.Println("loop end!")
			break
		}

		// fmt.Println(start, cpt, n)
		// fmt.Println(start, cpt, cpt <= start, start < cpt+int(n))
		if (cpt <= start && start < cpt+int(n)) || (cpt > start && cpt < end) {
			for j := minpt( start, cpt ); j+cpt < end && j < int(n); j++ {

				binary.LittleEndian.PutUint32(data[off+4:off+8],   math.Float32bits(m.Pts[uint16(i)][j].GetX()))
				binary.LittleEndian.PutUint32(data[off+8:off+12],  math.Float32bits(m.Pts[uint16(i)][j].GetY()))
				binary.LittleEndian.PutUint32(data[off+12:off+16], math.Float32bits(m.Pts[uint16(i)][j].GetZ()))

				off += 16

				// fmt.Println("Take", i, "point", j, n, cpt, start, end)
				// fmt.Println( math.Float32bits( m.Pts[uint16(i)][j].GetX() ) )
			}
		}
		cpt += int(n)
	}

	// fmt.Println("-----------request end-----------------")

	set_resp_headers(w)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	binary.Write(w, binary.LittleEndian, data)

	// fmt.Printf("No pts: %d\n", sz)
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


func minpt(a int,b int) int{
	if a - b > 0 {
		return a - b
	}
	return 0
}