package rawdate

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

var TimeZoneDB = time.Local // this variable defines the timezone in which data represented in the Database

var (
	_ fmt.Stringer             = RawDate{}
	_ encoding.TextMarshaler   = RawDate{}
	_ json.Marshaler           = RawDate{}
	_ fmt.GoStringer           = RawDate{}
	_ driver.Valuer            = RawDate{}
	_ encoding.TextUnmarshaler = (*RawDate)(nil)
	_ json.Unmarshaler         = (*RawDate)(nil)
	_ sql.Scanner              = (*RawDate)(nil)
)

// MarshalText implements the encoding.TextMarshaler
func (d RawDate) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// MarshalJSON implements `json.Marshaler` using YYYY-MM-DD
func (d RawDate) MarshalJSON() ([]byte, error) {
	s := d.String()
	return json.Marshal(s)
}

// UnmarshalText implements the encoding.TextUnmarshaler for YYYY-MM-DD
func (d *RawDate) UnmarshalText(data []byte) error {
	parsed, err := Parse(time.DateOnly, string(data))
	if err != nil {
		return err
	}

	*d = parsed
	return nil
}

// UnmarshalJSON implements `json.Unmarshaler`; parses the date as time.DateOnly (YYYY-MM-DD).
func (d *RawDate) UnmarshalJSON(data []byte) error {
	s := ""
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	var parsed RawDate
	parsed, err = Parse(time.DateOnly, s)
	if err != nil {
		return err
	}

	*d = parsed
	return nil
}

// Value implements `driver.Valuer`; marshals to a `time.Time`
// Notice that the timezone can be changed using global `TimeZoneDB`
func (d RawDate) Value() (driver.Value, error) {
	return d.Time(TimeZoneDB), nil
}

// GoString implements `fmt.GoStringer`.
func (d RawDate) GoString() string {
	return fmt.Sprintf("rawdate.MustNew(%d, time.%s, %d)", d.Year0+1, time.Month(d.Month0+1), d.Day0+1)
}

// Scan implements `sql.Scanner` aun unmarshalls `time.Time`
func (d *RawDate) Scan(src any) error {
	var t time.Time

	switch srcTyped := src.(type) {
	case time.Time:
		t = srcTyped
	default:
		return fmt.Errorf("wrong type to scan as RawDate: type=%T", src)
	}
	if t.Hour() != 0 || t.Minute() != 0 || t.Second() != 0 || t.Nanosecond() != 0 {
		return fmt.Errorf("timestamp contains time information: %s", t.String())
	}
	raw := FromTime(t)
	if !raw.IsValid() {
		return fmt.Errorf("not a valid time for RawDate: %s", t.String())
	}
	*d = raw
	return nil
}
