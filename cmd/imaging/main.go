package main

import (
	gobytes "bytes"
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/disintegration/gift"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

// TODO: change thread number to worker number & use cpu core amount.
func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "imaging")
	c, err := config.NewConfig()
	if err != nil {
		entry.Fatal(err)
	}
	client := &http.Client{Timeout: c.Imaging.Timeout}
	records := make(chan []string, c.Imaging.ThreadNumber)
	group := &sync.WaitGroup{}
	group.Add(c.Imaging.ThreadNumber)
	for i := 0; i < c.Imaging.ThreadNumber; i++ {
		go work(client, records, group, c.Imaging, entry)
	}
	file, err := os.Open(c.Imaging.InputPath)
	if err != nil {
		entry.Fatalf("main: imaging failed to open the input file, %v", err)
	}
	err = run(file, records)
	group.Wait()
	if err != nil {
		_ = file.Close()
		entry.Fatal(err)
	}
	if err := file.Close(); err != nil {
		entry.Fatalf("main: imaging failed to close the input file, %v", err)
	}
}

func work(
	client *http.Client,
	records <-chan []string,
	group *sync.WaitGroup,
	config config.Imaging,
	logger log.FieldLogger,
) {
	for record := range records {
		for retry, err := 1, io.EOF; retry <= config.RetryLimit && err != nil; retry++ {
			if err = load(client, config, record); err != nil {
				logger.WithFields(log.Fields{"url": record[0], "retry": retry}).Error(err)
			}
		}
	}
	group.Done()
}

func load(client *http.Client, config config.Imaging, record []string) error {
	request, err := http.NewRequest(http.MethodGet, record[0], nil)
	if err != nil {
		return fmt.Errorf("main: imaging failed to make a request, %v", err)
	}
	config.Headers.Inject(request)
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("main: imaging failed to send a request, %v", err)
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return fmt.Errorf("main: imaging got a non-ok status %s", response.Status)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return fmt.Errorf("main: imaging failed to read the response body, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return fmt.Errorf("main: imaging failed to close the response body, %v", err)
	}
	hash := sha1.Sum([]byte(record[0]))
	err = ioutil.WriteFile(
		misc.ResolvePath(fmt.Sprintf(config.OutputFormat, hash, "origin", record[1])),
		bytes,
		0644,
	)
	if err != nil {
		return fmt.Errorf("main: imaging failed to write the file, %v", err)
	}
	gifts := map[string]*gift.GIFT{
		"flip":              gift.New(gift.FlipHorizontal()),
		//"rotate_ccw1":       rotate(1),
		//"rotate_cw1":        rotate(-1),
		//"rotate_ccw2":       rotate(2),
		//"rotate_cw2":        rotate(-2),
		//"rotate_ccw3":       rotate(3),
		//"rotate_cw3":        rotate(-3),
		//"brightness_up5":    gift.New(gift.Brightness(5)),
		//"brightness_down5":  gift.New(gift.Brightness(-5)),
		//"brightness_up10":   gift.New(gift.Brightness(10)),
		//"brightness_down10": gift.New(gift.Brightness(-10)),
		//"balance_up":        gift.New(gift.ColorBalance(-5, 5, 10)),
		//"balance_down":      gift.New(gift.ColorBalance(-10, 10, 0)),
		//"contrast_up20":     gift.New(gift.Contrast(20)),
		//"contrast_down20":   gift.New(gift.Contrast(-20)),
		//"contrast_up30":     gift.New(gift.Contrast(30)),
		//"contrast_down30":   gift.New(gift.Contrast(-30)),
		"gamma_light":       gift.New(gift.Gamma(1.3)),
		"gamma_dark":        gift.New(gift.Gamma(0.7)),
		"blur":              gift.New(gift.GaussianBlur(1.25)),
		//"hue_ccw":           gift.New(gift.Hue(5)),
		//"hue_cw":            gift.New(gift.Hue(-5)),
		//"max":               gift.New(gift.Maximum(3, false)),
		//"min":               gift.New(gift.Minimum(3, false)),
		//"median":            gift.New(gift.Median(3, false)),
		//"mean":              gift.New(gift.Mean(3, false)),
		//"saturation_up":     gift.New(gift.Saturation(80)),
		//"saturation_down":   gift.New(gift.Saturation(-20)),
		//"sepia":             gift.New(gift.Sepia(20)),
		//"sigmoid":           gift.New(gift.Sigmoid(0.6, 3)),
		//"scale": gift.New(
		//	gift.CropToSize(600, 450, gift.CenterAnchor),
		//	gift.Resize(620, 460, gift.LanczosResampling),
		//),
	}
	source, _, err := image.Decode(gobytes.NewBuffer(bytes))
	if err != nil {
		return fmt.Errorf("main: imaging failed to decode the source, %v", err)
	}
	for name, g := range gifts {
		target := image.NewRGBA(g.Bounds(source.Bounds()))
		g.Draw(target, source)
		file, err := os.Create(
			misc.ResolvePath(fmt.Sprintf(config.OutputFormat, hash, name, record[1])),
		)
		if err != nil {
			return fmt.Errorf("main: imaging failed to create the target, %v", err)
		}
		if err := webp.Encode(file, target, nil); err != nil {
			_ = file.Close()
			return fmt.Errorf("main: imaging failed to encode the target, %v", err)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("main: imaging failed to close the target, %v", err)
		}
	}
	return nil
}

func rotate(angle float32) *gift.GIFT {
	return gift.New(
		gift.Rotate(angle, color.White, gift.CubicInterpolation),
		gift.CropToSize(620, 460, gift.CenterAnchor),
	)
}

func run(input io.Reader, records chan<- []string) error {
	reader := csv.NewReader(input)
	if _, err := reader.Read(); err == io.EOF {
		return nil
	} else if err != nil {
		return fmt.Errorf("main: imaging failed to read header of the input file, %v", err)
	}
	if reader.FieldsPerRecord != 2 {
		return fmt.Errorf(
			"main: imaging got invalid field number, %d != 2",
			reader.FieldsPerRecord,
		)
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			close(records)
			return nil
		} else if err != nil {
			close(records)
			return fmt.Errorf("main: imaging failed to read a row of the input file, %v", err)
		}
		records <- record
	}
}
