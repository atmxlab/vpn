package config

import (
	"encoding/json"
	"fmt"
	"time"
)

type Duration time.Duration

func (d *Duration) ToDuration() time.Duration {
	return time.Duration(*d)
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
	case float64:
		*d = Duration(time.Duration(value) * time.Second)
	default:
		return fmt.Errorf("invalid duration")
	}
	return nil
}
