package osmgeo_test

import (
	"github.com/paulmach/orb"
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

func TestConvertEmptyResponse(t *testing.T) {
	testConvert(t, "testdata/empty_response.osm", make([]orb.Geometry, 0))
}

func testConvert(t *testing.T, filePath string, expected []orb.Geometry) {
	actual, err := osmgeo.Convert(readFile(t, filePath))
	if err != nil {
		t.Fatalf("osmgeo_test: got non-nil error, %v", err)
	}
	if got, wanted := len(actual), len(expected); got != wanted {
		t.Fatalf("osmgeo_test: got invalid geometry lengths, %d != %d", got, wanted)
	}
	for i := range actual {
		if !orb.Equal(actual[i], expected[i]) {
			t.Errorf("osmgeo_test: got invalid geometries at %d, %v != %v", i, actual[i], expected[i])
		}
	}
}

func TestConvertSingleNode(t *testing.T) {
	testConvert(t, "testdata/single_node.osm", []orb.Geometry{orb.Point{30.5244014, 50.4495672}})
}
