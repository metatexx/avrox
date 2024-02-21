package rawdate

import (
	"errors"
	"fmt"
	"time"
)

var Zero = RawDate{0, 0, 0}

const ISODate = time.DateOnly

type RawDate struct {
	Year0  int
	Month0 int8
	Day0   int8
}

func New(y int, m time.Month, d int) (RawDate, error) {
	r := RawDate{Year0: y - 1, Month0: int8(m) - 1, Day0: int8(d) - 1}
	if !r.IsValid() {
		return Zero, errors.New("not a date")
	}
	return r, nil
}

func MustNew(y int, m time.Month, d int) RawDate {
	r := RawDate{Year0: y - 1, Month0: int8(m) - 1, Day0: int8(d) - 1}
	if !r.IsValid() {
		panic("not a date")
	}
	return r
}

func (r RawDate) Day() int {
	return int(r.Day0) + 1
}

func (r RawDate) Month() time.Month {
	return time.Month(r.Month0 + 1)
}

func (r RawDate) Year() int {
	return r.Year0 + 1
}

func (r RawDate) Weekday() time.Weekday {
	t := time.Date(r.Year(), r.Month(), r.Day(), 0, 0, 0, 0, time.UTC)
	return t.Weekday()
}

func (r RawDate) IsValid() bool {
	t := time.Date(r.Year0+1, time.Month(r.Month0)+1, int(r.Day0)+1, 0, 0, 0, 0, time.UTC)
	if t.Day() != r.Day() || t.Month() != r.Month() || t.Year() != r.Year() {
		return false
	}
	return true
}

func FromTime(t time.Time) RawDate {
	return RawDate{Year0: t.Year() - 1, Month0: int8(t.Month() - 1), Day0: int8(t.Day() - 1)}
}

func Parse(layout, s string) (RawDate, error) {
	t, err := time.Parse(layout, s)
	if err != nil {
		return Zero, err
	}
	return RawDate{t.Year() - 1, int8(t.Month() - 1), int8(t.Day() - 1)}, nil
}

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

func (r RawDate) Compare(a RawDate) int {
	return Compare(r, a)
}

func (r RawDate) IsZero() bool {
	return r.Year0 == 0 && r.Month0 == 0 && r.Day0 == 0
}

func (r RawDate) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", r.Year0+1, r.Month0+1, r.Day0+1)
}

func (r RawDate) Time(location *time.Location) time.Time {
	return time.Date(r.Year(), r.Month(), r.Day(), 0, 0, 0, 0, location)
}

func (r RawDate) Format(format string) string {
	tm := r.Time(time.UTC)
	return tm.Format(format)
}

func (r RawDate) AddDate(years, months, days int) RawDate {
	t := r.Time(time.UTC).AddDate(years, months, days)
	return MustNew(t.Year(), t.Month(), t.Day())
}
