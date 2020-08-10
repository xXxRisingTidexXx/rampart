package osmgeo_test

import (
	"github.com/xXxRisingTidexXx/rampart/internal/mining/osmgeo"
	"io/ioutil"
	"testing"
)

func TestConvertEmptyInput(t *testing.T) {
	geometries, err := osmgeo.Convert([]byte(""))
	if err == nil {
		t.Fatal("osmgeo_test: got nil error")
	}
	if err.Error() != "osmgeo: failed to unmarshal the osm, EOF" {
		t.Errorf("osmgeo_test: got invalid error, %v", err)
	}
	if len(geometries) != 0 {
		t.Errorf("osmgeo_test: got non-empty geometries, %v", geometries)
	}
}

func TestConvertUnclosedTag(t *testing.T) {
	geometries, err := osmgeo.Convert(readFile(t, "testdata/unclosed_tag.osm"))
	if err == nil {
		t.Fatal("osmgeo_test: got nil error")
	}
	if err.Error() != "osmgeo: failed to unmarshal the osm, XML syntax error on line 6: unexpected EOF" {
		t.Errorf("osmgeo_test: got invalid error, %v", err)
	}
	if len(geometries) != 0 {
		t.Errorf("osmgeo_test: got non-empty geometries, %v", geometries)
	}
}

func readFile(t *testing.T, filePath string) []byte {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("osmgeo_test: failed to read the file, %v", err)
	}
	return bytes
}
