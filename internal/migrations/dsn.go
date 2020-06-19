package migrations

import (
	"fmt"
	gourl "net/url"
	"os"
)

func getDSN(queryParams map[string]string) (string, error) {
	dsn := os.Getenv("RAMPART_DSN")
	if dsn == "" {
		return "", fmt.Errorf("migrations: dsn isn't configured")
	}
	url, err := gourl.Parse(dsn)
	if err != nil {
		return "", fmt.Errorf("migrations: invalid dsn %s, %v", dsn, err)
	}
	values, err := gourl.ParseQuery(url.RawQuery)
	if err != nil {
		return "", fmt.Errorf("migrations: invalid query params of %s, %v", dsn, err)
	}
	for key, value := range queryParams {
		values.Set(key, value)
	}
	url.RawQuery = values.Encode()
	return url.String(), nil
}
