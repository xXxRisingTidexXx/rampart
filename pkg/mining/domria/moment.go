package domria

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type moment time.Time

const (
	momentLength     = 3
	dateTimingLength = 2
	dateLength       = 3
	timingLength     = 3
)

func (m *moment) UnmarshalJSON(bytes []byte) error {
	length := len(bytes)
	if length < momentLength {
		return fmt.Errorf("domria: moment string is too short, %d", length)
	}
	s := string(bytes[1 : length-1])
	dateTiming := strings.Split(s, " ")
	if len(dateTiming) != dateTimingLength {
		return fmt.Errorf("domria: moment can't split date & timing, %s", s)
	}
	date := strings.Split(dateTiming[0], "-")
	if len(date) != dateLength {
		return fmt.Errorf("domria: moment cannot split date, %s", s)
	}
	timing := strings.Split(dateTiming[1], ":")
	if len(timing) != timingLength {
		return fmt.Errorf("domria: moment cannot split timing, %s", s)
	}
	year, err := strconv.Atoi(date[0])
	if err != nil {
		return err
	}
	month, err := strconv.Atoi(date[1])
	if err != nil {
		return err
	}
	day, err := strconv.Atoi(date[2])
	if err != nil {
		return err
	}
	hours, err := strconv.Atoi(timing[0])
	if err != nil {
		return err
	}
	minutes, err := strconv.Atoi(timing[1])
	if err != nil {
		return err
	}
	seconds, err := strconv.Atoi(timing[2])
	if err != nil {
		return err
	}
	*m = moment(time.Date(year, time.Month(month), day, hours, minutes, seconds, 0, time.Local).UTC())
	return nil
}
