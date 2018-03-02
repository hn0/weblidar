package model

import (
	"fmt"
	"github.com/jblindsay/lidario"
)

type Model struct {
	Numpts int
	Valid  bool
}

func CreateModel(path string) *Model {
	m := new(Model)
	fmt.Println( path )

	lf, err := lidario.NewLasFile(path, "r")
    if err != nil {
        fmt.Println(err)
    }
    defer lf.Close()

    m.Numpts = lf.Header.NumberPoints
    for i := 0; i < m.Numpts; i++ {
    	if pt, err := lf.LasPoint(i); err == nil {
    		fmt.Println("Point:", pt)
    	}
    }
    

	return m
}