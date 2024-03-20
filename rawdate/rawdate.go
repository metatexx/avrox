// Package rawdate provides a simple date handling utility without time.
// It uses zero-based values for year, month, and day to maintain consistency
// with time.Time, where an uninitialized instance defaults to the year 0001-01-01.
package rawdate

import (
	"errors"
	"fmt"
	"time"
)

// NOTICE: When changing code ensure all functions correctly interpret the zero-based
// values and document any specific behaviors or limitations due to this design choice.

// Zero represents a zero value for RawDate corresponding to the date 0001-01-01.
var Zero = RawDate{0, 0, 0}

// ISODate represents the format for a date-only value in ISO standard.
const ISODate = time.DateOnly

// RawDate represents a date without time information.
// Year0, Month0, and Day0 are zero-based values for year, month, and day respectively.
// These fields are exported for marshaling purposes and should be used with the understanding
// that they are zero-based.
type RawDate struct {
	Year0  int
	Month0 int8
	Day0   int8
}

// New creates a new RawDate given the year, month, and day.
// Returns an error if the provided values do not form a valid date.
func New(y int, m time.Month, d int) (RawDate, error) {
	r := RawDate{Year0: y - 1, Month0: int8(m) - 1, Day0: int8(d) - 1}
	if !r.IsValid() {
		return Zero, errors.New("not a date")
	}
	return r, nil
}

// MustNew creates a new RawDate given the year, month, and day.
// Panics if the provided values do not form a valid date.
func MustNew(y int, m time.Month, d int) RawDate {
	r := RawDate{Year0: y - 1, Month0: int8(m) - 1, Day0: int8(d) - 1}
	if !r.IsValid() {
		panic("not a date")
	}
	return r
}

// Today creates a new RawDate with the current date (in local time)
func Today() RawDate {
	t := time.Now()
	return RawDate{Year0: t.Year() - 1, Month0: int8(t.Month() - 1), Day0: int8(t.Day() - 1)}
}

// Day returns the day of the month for the RawDate.
func (r RawDate) Day() int {
	return int(r.Day0) + 1
}

// Month returns the month of the year for the RawDate.
func (r RawDate) Month() time.Month {
	return time.Month(r.Month0 + 1)
}

// Year returns the year for the RawDate.
func (r RawDate) Year() int {
	return r.Year0 + 1
}

// Weekday calculates and returns the day of the week for the RawDate.
func (r RawDate) Weekday() time.Weekday {
	t := time.Date(r.Year(), r.Month(), r.Day(), 0, 0, 0, 0, time.UTC)
	return t.Weekday()
}

// IsValid checks whether the RawDate represents a valid date.
func (r RawDate) IsValid() bool {
	t := time.Date(r.Year0+1, time.Month(r.Month0)+1, int(r.Day0)+1, 0, 0, 0, 0, time.UTC)
	if t.Day() != r.Day() || t.Month() != r.Month() || t.Year() != r.Year() {
		return false
	}
	return true
}

// FromTime creates a RawDate from a time.Time value, ignoring the time portion.
func FromTime(t time.Time) RawDate {
	return RawDate{Year0: t.Year() - 1, Month0: int8(t.Month() - 1), Day0: int8(t.Day() - 1)}
}

// Parse parses a string representing a date and returns the corresponding RawDate.
// An error is returned if the string does not represent a valid date.
func Parse(layout, s string) (RawDate, error) {
	t, err := time.Parse(layout, s)
	if err != nil {
		return Zero, err
	}
	return RawDate{t.Year() - 1, int8(t.Month() - 1), int8(t.Day() - 1)}, nil
}

// Compare compares two RawDates.
// Returns 1 if a > b, -1 if a < b, and 0 if a == b.
func Compare(a, b RawDate) int {
	s := a.Year0 - b.Year0
	if s == 0 {
		s = int(a.Month0 - b.Month0)
	}
	if s == 0 {
		s = int(a.Day0 - b.Day0)
	}
	if s > 0 {
		return 1
	}
	if s < 0 {
		return -1
	}
	return 0
}

// Compare compares the RawDate with another RawDate.
// Returns 1 if r > a, -1 if r < a, and 0 if r == a.
func (r RawDate) Compare(a RawDate) int {
	return Compare(r, a)
}

// After reports whether the RawDate a is after RawDate b.
func (a RawDate) After(b RawDate) bool {
	return Compare(a, b) > 0
}

// Before reports whether the RawDate a is before RawDate b.
func (a RawDate) Before(b RawDate) bool {
	return Compare(a, b) < 0
}

// IsZero checks if the RawDate is a zero date.
func (r RawDate) IsZero() bool {
	return r.Year0 == 0 && r.Month0 == 0 && r.Day0 == 0
}

// String returns a string representation of the RawDate in YYYY-MM-DD format.
func (r RawDate) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", r.Year0+1, r.Month0+1, r.Day0+1)
}

// Time converts the RawDate to a time.Time value, with the provided location and zero time.
func (r RawDate) Time(location *time.Location) time.Time {
	return time.Date(r.Year(), r.Month(), r.Day(), 0, 0, 0, 0, location)
}

// Format formats the RawDate according to the provided layout string.
func (r RawDate) Format(format string) string {
	tm := r.Time(time.UTC)
	return tm.Format(format)
}

// AddDate adds the specified number of years, months, and days to the RawDate.
// Returns a new RawDate with the updated values.
func (r RawDate) AddDate(years, months, days int) RawDate {
	t := r.Time(time.UTC).AddDate(years, months, days)
	return MustNew(t.Year(), t.Month(), t.Day())
}

// MonthStart returns a new RawDate that represents the first day of the month for the given RawDate.
func (r RawDate) MonthStart() RawDate {
	return MustNew(r.Year(), r.Month(), 1)
}

// MonthEnd returns a new RawDate that represents the last day of the month for the given RawDate.
func (r RawDate) MonthEnd() RawDate {
	daysInMonth := time.Date(r.Year(), r.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	return MustNew(r.Year(), r.Month(), daysInMonth)
}

// NextWeekday returns the next date that falls on the given weekday.
// If orToday is true and the given weekday is today, it returns the current date.
func (r RawDate) NextWeekday(weekday time.Weekday, orToday bool) RawDate {
	wd := r.Weekday()
	if orToday && wd == weekday {
		return r
	}

	difference := int(weekday - wd)
	if difference <= 0 {
		difference += 7
	}

	return r.AddDate(0, 0, difference)
}

// PreviousWeekday returns the previous date that falls on the given weekday.
// If orToday is true and the given weekday is today, it returns the current date.
func (r RawDate) PreviousWeekday(weekday time.Weekday, orToday bool) RawDate {
	wd := r.Weekday()
	if orToday && wd == weekday {
		return r
	}

	var difference int
	if r.Weekday() > weekday {
		difference = int(wd - weekday)
	} else {
		difference = int(wd - weekday + 7)
	}

	return r.AddDate(0, 0, -difference)
}
