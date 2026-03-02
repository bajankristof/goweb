package config

import "time"

type Duration time.Duration

func (d *Duration) UnmarshalText(text []byte) error {
	pd, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}

	*d = Duration(pd)
	return nil
}
