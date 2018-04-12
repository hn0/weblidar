package model

import (
	"clwrapper"
	"fmt"
	"github.com/jblindsay/lidario"
	"math"
	"reflect"
	"time"
)

type point struct {
	x float32
	y float32
	z float32
}

type Model struct {
	Numpts int
	Valid  bool
	Pts    map[uint16][]point
	Sizes  []uint16
}

var sqrt3 float32 = 1.73205080757

func CreateModel(path string) *Model {
	m := new(Model)

	lf, err := lidario.NewLasFile(path, "r")
	if err != nil {
		fmt.Println(err)
	}
	defer lf.Close()

	m.Numpts = lf.Header.NumberPoints
	// m.Numpts = 50000
	if m.Numpts > 0 {
		fmt.Printf("Processing input dataset containing %d points\n", m.Numpts)

		can_res := 2000
		sortgrd := create_sortgrid(can_res)

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
		var start time.Time
		m.Pts = make(map[uint16][]point)
		m.Sizes = []uint16{0}
		valid := true

		fmt.Println("Reading points:\n")
		valfnc := func(i int) (float32, float32, float32) {
			var x, y, z float32
			if x1, y1, z1, err := lf.GetXYZ(i); err == nil {
				// range in webgl -1 .. 1
				x = float32((x1/domains["X"][1])*2 - 1)
				y = float32((y1/domains["Y"][1])*2 - 1)
				z = float32((z1/domains["Z"][1])*2 - 1)
			} else {
				valid = false
			}
			if i == m.Numpts-1 {
				fmt.Printf("\r\t 100%% Done.\n")
				fmt.Println("Processing the data ...")
				start = time.Now()
			} else if i%500 == 0 {
				fmt.Printf("\r\t %f%%", (float64(i)/float64(m.Numpts))*100)
			}
			return x, y, z
		}

		resfnc := func(i int, xyz [3]float32, dist [2]float32) {
			// order of reading is not the best, categorization is next
			if i == 1 {
				fmt.Printf("\rDone. Processing took: %s\n", time.Since(start))
				fmt.Println("Reading results")
			} else if i == m.Numpts-1 {
				fmt.Printf("\r\t 100%% Done.\n")
			} else if i%500 == 0 {
				fmt.Printf("\r\t %f%%n", (float64(i)/float64(m.Numpts))*100)
			}

			// fmt.Println( sortgrd[idist] )

			// domain of the dist 0 .. sqrt(3)
			idist := int(math.Floor(float64(len(sortgrd)) * float64(dist[0]/sqrt3)))
			if idist > len(sortgrd) - 1{
				idist = len(sortgrd) - 1
			}
			// angle is in the range 0 .. 2Pi
			angdist := int(math.Floor(float64(len(sortgrd[0])) * float64(dist[1]/(2*math.Pi))))
			if angdist > len(sortgrd[0]) - 1{
				angdist = len(sortgrd[0]) - 1
			}
			// fmt.Println(sortgrd[idist][angdist], idist, angdist)

			pt := point{
				xyz[0],
				xyz[1],
				xyz[2],
			}

			m.Pts[sortgrd[idist][angdist]] = append(m.Pts[sortgrd[idist][angdist]], pt)

			if int(sortgrd[idist][angdist]) >= len(m.Sizes) {
				m.Sizes = append(m.Sizes, 0)
			}
			m.Sizes[sortgrd[idist][angdist]] += 1

			sortgrd[idist][angdist] += 1
		}

		// TODO: relative path?!
		p := clwrapper.Program{"src/clwrapper/euclid_dist.cl", "euclid_dist", valfnc, resfnc}
		valid = valid && clwrapper.RunProgram(&p, m.Numpts)

		if valid {
			m.Valid = true
			// fmt.Println(m.Pts)
			// fmt.Println(m.Sizes)
		}

	}

	return m
}

func create_sortgrid(size int) [][]uint16 {
	// use 2d !?
	grid := make([][]uint16, size)
	for i := range grid {
		grid[i] = make([]uint16, size)
	}
	return grid
}

func (p *point) GetX() float32 {
	return p.x
}

func (p *point) GetY() float32 {
	return p.y
}

func (p *point) GetZ() float32 {
	return p.z
}
