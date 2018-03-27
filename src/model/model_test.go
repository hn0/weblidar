package model

import(
	"testing"
	"os"
)


var m *Model

func TestGrid(t *testing.T) {
	create_sortgrid(10)
}

func TestFile(t *testing.T) {

	// var lf *lidario.LasFile

	fname := "../../data/sample.las"
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		t.Error("Sample file path is not valid, cannot perform the test")
		t.Fail()
	} else {
		m = CreateModel( fname )
		if ! m.Valid {
			t.Error("Did not get valid model back!")
			t.Fail()
		}
		t.Fail()
	}

}