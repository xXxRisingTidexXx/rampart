package domria

import (
	"encoding/json"
	"strconv"
)

type coordinate float64

func (c *coordinate) UnmarshalJSON(bytes []byte) error {
	if bytes[0] != '"' {
		return json.Unmarshal(bytes, (*float64)(c))
	}
	if string(bytes) == "\"\"" {
		*c = coordinate(0)
		return nil
	}
	s := ""
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*c = coordinate(f)
	return nil
}
