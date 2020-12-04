package main

import (
	"github.com/disintegration/gift"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"image"
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
	g := gift.New(gift.Transpose())
	target := image.NewRGBA(g.Bounds(source.Bounds()))
	g.Draw(target, source)
	file, err = os.Create("images/target.webp")
	if err != nil {
		entry.Fatalf("main: augmentation failed to create the target, %v", err)
	}
	if err := file.Close(); err != nil {
		entry.Fatalf("main: augmentation failed to close the target, %v", err)
	}
}
