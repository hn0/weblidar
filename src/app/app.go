package main

import (
	"fmt"
	"net/http"
	"os"
	"model"
)

var PORT int = 3000

func RouteHandeler(w http.ResponseWriter, r *http.Request) {

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

	// TODO: check the structure of the file, in new class
	//  and build structure for ....
	m := model.CreateModel(os.Args[1]);
	fmt.Println(m)

	fmt.Println("Starting web server on port", PORT)

	http.HandleFunc("/data", RouteHandeler)
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", PORT), http.FileServer(http.Dir("./static"))))
}
