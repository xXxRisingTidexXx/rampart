package mining

import (
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"net/http"
)

// TODO: relative city center distance feature (with city diameter).
func NewGaugingAmplifier(config config.GaugingAmplifier) Amplifier {
	return &gaugingAmplifier{
		&http.Client{Timeout: config.Timeout},
		config.Host,
		config.InterpreterFormat,
		config.UserAgent,
		config.SubwayCities,
		-1,
		config.SSFSearchRadius,
		config.SSFMinDistance,
		config.SSFModifier,
		config.IZFSearchRadius,
		config.IZFMinArea,
		config.IZFMinDistance,
		config.IZFModifier,
		config.GZFSearchRadius,
		config.GZFMinArea,
		config.GZFMinDistance,
		config.GZFModifier,
	}
}

type gaugingAmplifier struct {
	client            *http.Client
	host              string
	interpreterFormat string
	userAgent         string
	subwayCities      misc.Set
	unknownDistance   float64
	ssfSearchRadius   float64
	ssfMinDistance    float64
	ssfModifier       float64
	izfSearchRadius   float64
	izfMinArea        float64
	izfMinDistance    float64
	izfModifier       float64
	gzfSearchRadius   float64
	gzfMinArea        float64
	gzfMinDistance    float64
	gzfModifier       float64
}

func (a *gaugingAmplifier) AmplifyFlat(flat Flat) (Flat, error) {
	return flat, nil
}
