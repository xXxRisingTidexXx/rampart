package main

import (
	"github.com/chai2010/webp"
	"github.com/disintegration/gift"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"image"
	"image/color"
	"os"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "augmentation")
	file, err := os.Open(misc.ResolvePath("images/source.webp"))
	if err != nil {
		entry.Fatalf("main: augmentation failed to open the source, %v", err)
	}
	source, _, err := image.Decode(file)
	if err != nil {
		_ = file.Close()
		entry.Fatalf("main: augmentation failed to decode the source, %v", err)
	}
	if err := file.Close(); err != nil {
		entry.Fatalf("main: augmentation failed to close the source, %v", err)
	}
	gifts := map[string]*gift.GIFT{
		"flip":              gift.New(gift.FlipHorizontal()),
		"rotate_ccw1":       rotate(1),
		"rotate_cw1":        rotate(-1),
		"rotate_ccw2":       rotate(2),
		"rotate_cw2":        rotate(-2),
		"rotate_ccw3":       rotate(3),
		"rotate_cw3":        rotate(-3),
		"brightness_up5":    gift.New(gift.Brightness(5)),
		"brightness_down5":  gift.New(gift.Brightness(-5)),
		"brightness_up10":   gift.New(gift.Brightness(10)),
		"brightness_down10": gift.New(gift.Brightness(-10)),
		"balance_up":        gift.New(gift.ColorBalance(-5, 5, 10)),
		"balance_down":      gift.New(gift.ColorBalance(-10, 10, 0)),
		"contrast_up20":     gift.New(gift.Contrast(20)),
		"contrast_down20":   gift.New(gift.Contrast(-20)),
		"contrast_up30":     gift.New(gift.Contrast(30)),
		"contrast_down30":   gift.New(gift.Contrast(-30)),
		"gamma_light":       gift.New(gift.Gamma(1.3)),
		"gamma_dark":        gift.New(gift.Gamma(0.7)),
		"blur":              gift.New(gift.GaussianBlur(1.25)),
		"hue_ccw":           gift.New(gift.Hue(5)),
		"hue_cw":            gift.New(gift.Hue(-5)),
		"max":               gift.New(gift.Maximum(3, false)),
		"min":               gift.New(gift.Minimum(3, false)),
		"median":            gift.New(gift.Median(3, false)),
		"mean":              gift.New(gift.Mean(3, false)),
		"saturation_up":     gift.New(gift.Saturation(80)),
		"saturation_down":   gift.New(gift.Saturation(-20)),
		"sepia":             gift.New(gift.Sepia(20)),
		"sigmoid":           gift.New(gift.Sigmoid(0.6, 3)),
		"scale": gift.New(
			gift.CropToSize(600, 450, gift.CenterAnchor),
			gift.Resize(620, 460, gift.LanczosResampling),
		),
	}
	for name, g := range gifts {
		entry = entry.WithField("gift", name)
		target := image.NewRGBA(g.Bounds(source.Bounds()))
		g.Draw(target, source)
		file, err = os.Create(misc.ResolvePath("images/target_" + name + ".webp"))
		if err != nil {
			entry.Fatalf("main: augmentation failed to create the target, %v", err)
		}
		if err := webp.Encode(file, target, &webp.Options{Lossless: true}); err != nil {
			_ = file.Close()
			entry.Fatalf("main: augmentation failed to encode the target, %v", err)
		}
		if err := file.Close(); err != nil {
			entry.Fatalf("main: augmentation failed to close the target, %v", err)
		}
	}
}

func rotate(angle float32) *gift.GIFT {
	return gift.New(
		gift.Rotate(angle, color.White, gift.CubicInterpolation),
		gift.CropToSize(620, 460, gift.CenterAnchor),
	)
}
