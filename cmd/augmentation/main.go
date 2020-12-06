package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/imaging"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"image"
	"io/ioutil"
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
	for _, effect := range imaging.Effects {
		entry = entry.WithField("effect", effect.Name())
		bytes, err := effect.Apply(source)
		if err != nil {
			entry.Fatal(err)
		}
		err = ioutil.WriteFile(
			misc.ResolvePath("images/target_"+effect.Name()+".webp"),
			bytes,
			0644,
		)
		if err != nil {
			entry.Fatalf("main: augmentation failed to write the target, %v", err)
		}
	}
}
