package model

import (
	"clwrapper"
	"fmt"
	"github.com/jblindsay/lidario"
	"reflect"
)

type point struct {
	x     float32
	y     float32
	z     float32
	class int
}

type Model struct {
	Numpts int
	Valid  bool
	Pts    []point
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
		fmt.Printf("Processing input dataset containing %d points\n", m.Numpts)

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
		m.Pts = make([]point, m.Numpts)
		valid := true

		fmt.Println("Reading points:\n")
		valfnc := func(i int) (float32, float32, float32) {
			var x, y, z float32
			if x1, y1, z1, err := lf.GetXYZ(i); err == nil {
				// let say that domain will be 0 .. 1?
				// SOMETHING THAT NEEDS A CHECK!
				x = float32(x1 / domains["X"][1])
				y = float32(y1 / domains["Y"][1])
				z = float32(z1 / domains["Z"][1])
				m.Pts[i].x = x
				m.Pts[i].y = y
				m.Pts[i].z = z
			} else {
				valid = false
			}
			if i%500 == 0 {
				fmt.Printf("\r\t %f%%", (float64(i)/float64(m.Numpts))*100)
			}
			return x, y, z
		}
		// TODO: relative path?!
		p := clwrapper.Program{"src/clwrapper/euclid_dist.cl", "euclid_dist", valfnc}
		valid = valid && clwrapper.RunProgram(&p, 32)

		fmt.Printf("\r\tdone \n")

		if valid {
			m.Valid = true
			// m.calculate_dist()
		}

	}

	return m
}

func (m *Model) calculate_dist() {

	// fmt.Println("Classifying point cloud")

	// acctually we don't need such huge matrix, we only need final number of neighbors!?
	// dist := make([]uint8, n*n)
	// for now let just use cpu and see performance!
	// r := [2]float64{math.MaxFloat64, 0}
	// n := 10
	// for i := 0; i < n; i++ {
	// 	for j := n; j > i; j-- {
	// 		x := square(m.Pts[i].x - m.Pts[j].x)
	// 		y := square(m.Pts[i].y - m.Pts[j].y)
	// 		z := square(m.Pts[i].z - m.Pts[j].z)
	// 		d := math.Sqrt(x + y + z)
	// 		// fmt.Println(d)
	// 		r[0] = math.Min(r[0], d)
	// 		r[1] = math.Max(r[1], d)
	// 	}
	// }
	// fmt.Println(r)
	// TODO n by n matrix, maybe even using opencl
	// calculate simple euclidian distances!
}
