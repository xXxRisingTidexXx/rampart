package imaging

import (
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/disintegration/gift"
	"image"
	"image/color"
)

var Effects = []Effect{
	NewEffect("flip", gift.FlipHorizontal()),
	NewEffect(
		"brightness_up_rotate_ccw_crop",
		gift.Brightness(10),
		gift.Rotate(5, color.White, gift.CubicInterpolation),
		gift.CropToSize(620, 460, gift.CenterAnchor),
	),
	NewEffect(
		"brightness_down_rotate_cw_crop",
		gift.Brightness(-20),
		gift.Rotate(-3, color.White, gift.CubicInterpolation),
		gift.CropToSize(620, 460, gift.CenterAnchor),
	),
	NewEffect("balance_up_saturation_up", gift.ColorBalance(-5, 5, 10), gift.Saturation(80)),
	NewEffect(
		"balance_down_hue_cw_rotate_ccw_crop",
		gift.ColorBalance(-10, 10, 0),
		gift.Hue(-5),
		gift.Rotate(3, color.White, gift.CubicInterpolation),
		gift.CropToSize(620, 460, gift.CenterAnchor),
	),
	NewEffect("contrast_up_hue_ccw", gift.Contrast(30), gift.Hue(5)),
	NewEffect(
		"gamma_light_crop_resize",
		gift.Gamma(1.3),
		gift.CropToSize(589, 437, gift.TopRightAnchor),
		gift.Resize(620, 460, gift.CubicResampling),
	),
	NewEffect("gamma_dark_flip", gift.Gamma(0.7), gift.FlipHorizontal()),
	NewEffect("blur", gift.GaussianBlur(1.25)),
	NewEffect(
		"max_rotate_cw_crop",
		gift.Maximum(3, false),
		gift.Rotate(2, color.White, gift.CubicInterpolation),
		gift.CropToSize(620, 460, gift.CenterAnchor),
	),
	NewEffect("min_flip", gift.Minimum(3, true), gift.FlipHorizontal()),
	NewEffect("median", gift.Median(3, false)),
	NewEffect("mean", gift.Mean(5, true)),
	NewEffect("sigmoid_sepia", gift.Sigmoid(0.6, 3), gift.Sepia(20)),
	NewEffect(
		"crop_resize_saturation_down",
		gift.CropToSize(600, 450, gift.CenterAnchor),
		gift.Resize(620, 460, gift.LanczosResampling),
		gift.Saturation(-15),
	),
}

func NewEffect(name string, filters ...gift.Filter) Effect {
	return Effect{name, gift.New(filters...)}
}

type Effect struct {
	name string
	gift *gift.GIFT
}

func (effect Effect) Name() string {
	return effect.name
}

func (effect Effect) Apply(source image.Image) ([]byte, error) {
	target := image.NewRGBA(effect.gift.Bounds(source.Bounds()))
	effect.gift.Draw(target, source)
	var buffer bytes.Buffer
	if err := webp.Encode(&buffer, target, nil); err != nil {
		return nil, fmt.Errorf("imaging: effect failed to encode the target, %v", err)
	}
	return buffer.Bytes(), nil
}
