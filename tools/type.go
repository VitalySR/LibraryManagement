package tools

import (
	"encoding/json"
	"strings"
	"time"
)

type CustomDate struct {
	time.Time
}

const layout = "2006-01-02"

func (c *CustomDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`) // remove quotes
	if s == "null" {
		return nil
	}

	c.Time, err = time.Parse(layout, s)
	return err
}

func (c *CustomDate) MarshalJSON() ([]byte, error) {
	if c.Time.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(c.Time.Format(layout))
}

func (c *CustomDate) Scan(value interface{}) (err error) {
	if value != nil {
		c.Time = value.(time.Time)
	}
	return nil
}
