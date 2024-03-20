package rawdate

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

// The variables below are used to configure the database serialization behavior.
// WARNING:
// These variables are global and modifying them can lead to inconsistencies in
// concurrent routines, especially if different routines expect different settings.
// Primarily, these should be set according to the specific database variant in use.
// Caution is advised as these settings also affect imported packages that rely on them.

var SQLDateFormat = "2006-01-02" // This is the format for the DB when Parse uses a string (default)
var SQLZeroValue = "0000-01-01"  // This is what is used as Zero value in the DB (only with Parse using string)
var SQLValueTime = false         // Setting this to true uses a time.Time when unmarshalling for SQL
var SQLTimeZone = time.UTC       // This is the timezone that gets used for the time.Time value

// These are alternate formats which work well with old MSSQL Servers

const MSSQLDateFmt = "20060102"
const MSSQLZeroDate = "19000101"

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

// Value implements `driver.Valuer`; marshals to a string
func (d RawDate) Value() (driver.Value, error) {
	if SQLValueTime {
		return d.ValueTime()
	}
	// using 1900-01-01 as zero value when the given date is the zero value
	if d.IsZero() {
		return SQLZeroValue, nil
	}
	// the internal format is UTC make sure we use local time here!
	return d.Format(SQLDateFormat), nil
}

// This implementation will use the time.Time serialisation
// With test fields as destination this creates something like "2024-03-20 00:00:00+01:00"
func (d RawDate) ValueTime() (driver.Value, error) {
	return d.Time(SQLTimeZone), nil
}

// GoString implements `fmt.GoStringer`.
func (d RawDate) GoString() string {
	return fmt.Sprintf("rawdate.MustNew(%d, time.%s, %d)", d.Year0+1, time.Month(d.Month0+1), d.Day0+1)
}

// Scan implements `sql.Scanner` by unmarshalling from `time.Time`
func (d *RawDate) Scan(src any) error {
	var rd RawDate
	var err error

	switch st := src.(type) {
	case string:
		rd, err = Parse(SQLDateFormat, st)
		if err != nil {
			return err
		}

	case time.Time:
		if st.Hour() != 0 || st.Minute() != 0 || st.Second() != 0 || st.Nanosecond() != 0 {
			return fmt.Errorf("timestamp contains time information: %s", st.String())
		}
		rd, err = New(st.Year(), st.Month(), st.Day())
		if err != nil {
			return err
		}
		if !rd.IsValid() {
			return fmt.Errorf("not a valid time for RawDate: %s", st.String())
		}
	default:
		return fmt.Errorf("wrong type to scan as RawDate: type=%T", src)
	}
	*d = rd
	return nil
}
