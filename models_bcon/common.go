package models_bcon

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

type Duration struct {
	time.Duration
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	var dest any
	if err = json.Unmarshal(b, &dest); err != nil {
		return err
	}

	switch val := dest.(type) {
	case float64:
		*d = Duration{time.Duration(val)}
	case string:
		var tmp time.Duration
		if tmp, err = time.ParseDuration(val); err != nil {
			break
		}
		*d = Duration{tmp}
	default:
		err = errors.New("invalid format Duration")
	}

	return err
}

func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tm string
	if err := unmarshal(&tm); err != nil {
		return err
	}

	td, err := time.ParseDuration(tm)
	if err != nil {
		return fmt.Errorf("failed to parse '%s' to time.Duration: %v", tm, err)
	}

	*d = Duration{td}
	return nil
}

func (d *Duration) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(d.String())
}
