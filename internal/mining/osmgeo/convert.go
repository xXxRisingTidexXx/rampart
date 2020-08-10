package osmgeo

import (
	"encoding/xml"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/osm"
)

func Convert(bytes []byte) ([]orb.Geometry, error) {
	gosm := osm.OSM{}
	if err := xml.Unmarshal(bytes, &gosm); err != nil {
		return nil, fmt.Errorf("osmgeo: failed to unmarshal the osm, %v", err)
	}
	return make([]orb.Geometry, 0), nil
}
