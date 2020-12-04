package main

import (
	"github.com/disintegration/gift"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	_ "golang.org/x/image/webp"
	"image"
	"os"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "augmentation")
	file, err := os.Open(misc.ResolvePath("images/demo.webp"))
	if err != nil {
		entry.Fatalf("main: augmentation failed to open the file, %v", err)
	}
	_, _, err = image.Decode(file)
	if err != nil {
		_ = file.Close()
		entry.Fatalf("main: augmentation failed to decode the image, %v", err)
	}
	_ = gift.New()  // TODO
	if err := file.Close(); err != nil {
		entry.Fatalf("main: augmentation failed to close the file, %v", err)
	}
}
