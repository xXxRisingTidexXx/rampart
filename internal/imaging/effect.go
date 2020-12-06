package imaging

import (
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/disintegration/gift"
	"image"
	"image/color"
)

var Effects = []*Effect{
	NewEffect("flip", gift.FlipHorizontal()),
	NewEffect("rotate_ccw1", rotate(1)...),
	NewEffect("rotate_cw1", rotate(-1)...),
	NewEffect("rotate_ccw2", rotate(2)...),
	NewEffect("rotate_cw2", rotate(-2)...),
	NewEffect("rotate_ccw3", rotate(3)...),
	NewEffect("rotate_cw3", rotate(-3)...),
	NewEffect("brightness_up5", gift.Brightness(5)),
	NewEffect("brightness_down5", gift.Brightness(-5)),
	NewEffect("brightness_up10", gift.Brightness(10)),
	NewEffect("brightness_down10", gift.Brightness(-10)),
	NewEffect("balance_up", gift.ColorBalance(-5, 5, 10)),
	NewEffect("balance_down", gift.ColorBalance(-10, 10, 0)),
	NewEffect("contrast_up20", gift.Contrast(20)),
	NewEffect("contrast_down20", gift.Contrast(-20)),
	NewEffect("contrast_up30", gift.Contrast(30)),
	NewEffect("contrast_down30", gift.Contrast(-30)),
	NewEffect("gamma_light", gift.Gamma(1.3)),
	NewEffect("gamma_dark", gift.Gamma(0.7)),
	NewEffect("blur", gift.GaussianBlur(1.25)),
	NewEffect("hue_ccw", gift.Hue(5)),
	NewEffect("hue_cw", gift.Hue(-5)),
	NewEffect("max", gift.Maximum(3, false)),
	NewEffect("min", gift.Minimum(3, false)),
	NewEffect("median", gift.Median(3, false)),
	NewEffect("mean", gift.Mean(3, false)),
	NewEffect("saturation_up", gift.Saturation(80)),
	NewEffect("saturation_down", gift.Saturation(-20)),
	NewEffect("sepia", gift.Sepia(20)),
	NewEffect("sigmoid", gift.Sigmoid(0.6, 3)),
	NewEffect(
		"scale",
		gift.CropToSize(600, 450, gift.CenterAnchor),
		gift.Resize(620, 460, gift.LanczosResampling),
	),
}

func NewEffect(name string, filters ...gift.Filter) *Effect {
	return &Effect{name, gift.New(filters...)}
}

type Effect struct {
	name string
	gift *gift.GIFT
}

func (effect *Effect) Name() string {
	return effect.name
}

func (effect *Effect) Apply(source image.Image) ([]byte, error) {
	target := image.NewRGBA(effect.gift.Bounds(source.Bounds()))
	effect.gift.Draw(target, source)
	var buffer bytes.Buffer
	if err := webp.Encode(&buffer, target, nil); err != nil {
		return nil, fmt.Errorf(
			"imaging: effect %s failed to encode the target, %v",
			effect.name,
			err,
		)
	}
	return buffer.Bytes(), nil
}

func rotate(angle float32) []gift.Filter {
	return []gift.Filter{
		gift.Rotate(angle, color.White, gift.CubicInterpolation),
		gift.CropToSize(620, 460, gift.CenterAnchor),
	}
}
