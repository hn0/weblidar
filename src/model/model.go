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
	if m.Numpts > 0 {
		fmt.Printf("Processing input dataset containing %i points\n", m.Numpts)

		domains := make(map[string][2]float64)
		directions := []string{"X", "Y", "Z"}
		r := reflect.ValueOf(lf.Header)

		for _, k := range directions {
			vals := [2]float64{0, 0}
			for i, v := range []string{"Min", "Max"} {
				f := reflect.Indirect(r).FieldByName(v + k)
				vals[i] = float64(f.Float())
			}
			domains[k] = vals
		}

		// read and normalize all the points
		points := make([]point, m.Numpts)
		valid := true
		fmt.Println("Reading points:\n")

		for i := 0; i < m.Numpts; i++ {
			if x, y, z, err := lf.GetXYZ(1000); err == nil {
				// let say that domain will be 0 .. 1?
				// SOMETHING THAT NEEDS A CHECK!
				points[i].x = x / domains["X"][1]
				points[i].y = y / domains["Y"][1]
				points[i].z = z / domains["Z"][1]
			} else {
				valid = false
			}
			if i%500 == 0 {
				fmt.Printf("\r\t %f%%", (float64(i)/float64(m.Numpts))*100)
			}
		}
		fmt.Printf("\r\tdone \n")

		if valid {
			m.Pts = &points
			m.Valid = true
			m.calculate_dist()
		}

	}

	return m
}

func (m *Model) calculate_dist() {

	fmt.Println("Classifying point cloud")
	// TODO n by n matrix, maybe even using opencl
	// calculate simple euclidian distances!
}
