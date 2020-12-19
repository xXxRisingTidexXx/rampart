package main

import (
	gobytes "bytes"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/imaging"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"image"
	"io/ioutil"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "travesty")
	bytes, err := ioutil.ReadFile(misc.ResolvePath("images/source.webp"))
	if err != nil {
		entry.Fatalf("main: travesty failed to read the source, %v", err)
	}
	source, _, err := image.Decode(gobytes.NewBuffer(bytes))
	if err != nil {
		entry.Fatalf("main: travesty failed to decode the source, %v", err)
	}
	for _, effect := range imaging.Effects {
		entry = entry.WithField("effect", effect.Name())
		bytes, err = effect.Apply(source)
		if err != nil {
			entry.Fatal(err)
		}
		err = ioutil.WriteFile(
			misc.ResolvePath("images/target_"+effect.Name()+".webp"),
			bytes,
			0644,
		)
		if err != nil {
			entry.Fatalf("main: travesty failed to write the target, %v", err)
		}
	}
}
