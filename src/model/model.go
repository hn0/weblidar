package model

import (
	"fmt"
	"github.com/jblindsay/lidario"
	"reflect"
)

type point struct {
	x     float64
	y     float64
	z     float64
	class int
}

type Model struct {
	Numpts int
	Valid  bool
	Pts    *[]point
}

func CreateModel(path string) *Model {
	m := new(Model)
	fmt.Println(path)

	lf, err := lidario.NewLasFile(path, "r")
	if err != nil {
		fmt.Println(err)
	}
	defer lf.Close()

	m.Numpts = lf.Header.NumberPoints
	fmt.Printf("Processing input dataset containing %i points\n", m.Numpts)

	domains := make(map[string][2]float64)
	directions := []string{"X", "Y", "Z"}
	r := reflect.ValueOf(lf.Header)

	for _, k := range directions {
		vals := [2]float64{0, 0}
		for i, v := range []string{"Min", "Max"} {
			f := reflect.Indirect(r).FieldByName(v + k)
			vals[i] = float64(f.Float())
			fmt.Println(f)
		}
		domains[k] = vals
	}
	fmt.Println(domains)

	// read and normalize all the points
	points := make([]point, m.Numpts)
	fmt.Println("Reading points:\n")

	for i := 0; i < m.Numpts; i++ {
		// if pt, err := lf.LasPoint(i); err == nil {
		// fmt.Println("Point:", pt)

		// }
		if i%500 == 0 {
			fmt.Printf("\r\t %d of %d", i, m.Numpts)
		}
	}

	m.Pts = &points

	return m
}
