package glucose

import "errors"

type NewReadingsEvent struct {
}

func (e NewReadingsEvent) Topic() string {
	return "newReading"
}

func (e NewReadingsEvent) Decode(v any) error {
	return errors.New("not implemented")
}
