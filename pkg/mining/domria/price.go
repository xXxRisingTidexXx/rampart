package domria

import (
	"fmt"
	"strconv"
	"strings"
)

type price float64

const priceLength = 3

func (p *price) UnmarshalJSON(bytes []byte) error {
	length := len(bytes)
	if length < priceLength {
		return fmt.Errorf("domria: price string is too short, %d", length)
	}
	f, err := strconv.ParseFloat(strings.ReplaceAll(string(bytes[1:length-1]), " ", ""), 64)
	if err != nil {
		return err
	}
	*p = price(f)
	return nil
}
