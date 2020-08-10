package osmgeo_test

import (
	"github.com/xXxRisingTidexXx/rampart/internal/mining/osmgeo"
	"testing"
)

func TestConvertEmptyInput(t *testing.T) {
	geometries, err := osmgeo.Convert([]byte(""))
	if err != nil {
		t.Fatalf("osmgeo_test: got non-nil error, %v", err)
	}
	if len(geometries) != 0 {
		t.Errorf("osmgeo_test: got non-empty geometries, %v", geometries)
	}
}
