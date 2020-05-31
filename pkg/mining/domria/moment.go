package domria

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type moment time.Time

func (m *moment) UnmarshalJSON(bytes []byte) error {
	s := string(bytes)
	dateTiming := strings.Split(s, " ")
	if len(dateTiming) != 2 {
		return fmt.Errorf("domria: moment can't split date & timing, %s", s)
	}
	date := strings.Split(dateTiming[0], "-")
	if len(date) != 3 {
		return fmt.Errorf("domria: moment cannot split date, %s", s)
	}
	timing := strings.Split(dateTiming[1], ":")
	if len(timing) != 3 {
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
	*m = moment(time.Date(year, time.Month(month), day, hours, minutes, seconds, 0, time.Local))
	return nil
}
