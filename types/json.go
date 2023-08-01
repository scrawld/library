package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var (
	dateT = reflect.TypeOf((*Date)(nil))
	timeT = reflect.TypeOf((*Time)(nil))
)

// Date marshal/unmarshal to a fixed-format JSON string.
type Date time.Time

// String returns a fixed format string for time.
func (t Date) String() string { return time.Time(t).Format("2006-01-02") }

// ToTime return time.Time type
func (t Date) ToTime() time.Time { return time.Time(t) }

// MarshalText implements encoding.TextMarshaler
func (t Date) MarshalText() ([]byte, error) {
	stamp := time.Time(t).Format("2006-01-02")
	return []byte(stamp), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *Date) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return errNonString(dateT)
	}
	return wrapTypeError(t.UnmarshalText(input[1:len(input)-1]), dateT)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *Date) UnmarshalText(text []byte) error {
	//tim, err := time.Parse("2006-01-02", string(text))
	tim, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(string(text)), time.Local)
	if err != nil {
		return fmt.Errorf("unmarshal time error: %s", err)
	}
	*t = Date(tim)
	return nil
}

// Time marshal/unmarshal to a fixed-format JSON string.
type Time time.Time

// String returns a fixed format string for time.
func (t Time) String() string { return time.Time(t).Format("2006-01-02 15:04:05") }

// ToTime return time.Time type
func (t Time) ToTime() time.Time { return time.Time(t) }

// MarshalText implements encoding.TextMarshaler
func (t Time) MarshalText() ([]byte, error) {
	stamp := time.Time(t).Format("2006-01-02 15:04:05")
	return []byte(stamp), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *Time) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return errNonString(timeT)
	}
	return wrapTypeError(t.UnmarshalText(input[1:len(input)-1]), timeT)
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (t *Time) UnmarshalText(text []byte) error {
	//tim, err := time.Parse("2006-01-02 15:04:05", string(text))
	tim, err := time.ParseInLocation("2006-01-02 15:04:05", strings.TrimSpace(string(text)), time.Local)
	if err != nil {
		return fmt.Errorf("unmarshal time error: %s", err)
	}
	*t = Time(tim)
	return nil
}

func isString(input []byte) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

func wrapTypeError(err error, typ reflect.Type) error {
	if err != nil {
		return &json.UnmarshalTypeError{Value: err.Error(), Type: typ}
	}
	return err
}

func errNonString(typ reflect.Type) error {
	return &json.UnmarshalTypeError{Value: "non-string", Type: typ}
}
