package mining

import (
	"fmt"
	"strconv"
	"strings"
)

type price float64

func (p *price) UnmarshalJSON(bytes []byte) error {
	length := len(bytes)
	if length < 3 {
		return fmt.Errorf("mining: price string is too short, %d", length)
	}
	f, err := strconv.ParseFloat(strings.ReplaceAll(string(bytes[1:length-1]), " ", ""), 64)
	if err != nil {
		return err
	}
	*p = price(f)
	return nil
}
