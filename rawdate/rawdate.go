package rawdate

import (
	"errors"
	"fmt"
	"time"
)

var Zero = RawDate{1, 1, 1}

const ISODate = time.DateOnly

type RawDate struct {
	Year  int
	Month int8
	Day   int8
}

func New(y, m, d int) (RawDate, error) {
	r := RawDate{Year: y, Month: int8(m), Day: int8(d)}
	if !r.IsValid() {
		return Zero, errors.New("not a date")
	}
	return r, nil
}

func MustNew(y, m, d int) RawDate {
	r := RawDate{Year: y, Month: int8(m), Day: int8(d)}
	if !r.IsValid() {
		panic("not a date")
	}
	return r
}

func (r RawDate) IsValid() bool {
	t := time.Date(r.Year, time.Month(r.Month), int(r.Day), 0, 0, 0, 0, time.UTC)
	if t.Day() != int(r.Day) || t.Month() != time.Month(r.Month) || t.Year() != r.Year {
		return false
	}
	return true
}

func FromTime(t time.Time) RawDate {
	return RawDate{Year: t.Year(), Month: int8(t.Month()), Day: int8(t.Day())}
}

func Parse(layout, s string) (RawDate, error) {
	/*
		parts := strings.Split(s, "-")
		if len(parts) != 3 {
			return Zero, errors.New("not a date")
		}
		var y, m, d int
		var err error
		y, err = strconv.Atoi(parts[0])
		if err != nil {
			return Zero, errors.New("not a date")
		}
		m, err = strconv.Atoi(parts[1])
		if err != nil || m < 1 || m > 12 {
			return Zero, errors.New("not a date")
		}
		d, err = strconv.Atoi(parts[2])
		if err != nil {
			return Zero, errors.New("not a date")
		}
		t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
		if t.Day() != d || t.Month()!=time.Month(m) || t.Year() !=y {
			return Zero, errors.New("not a date")
		}
	*/
	t, err := time.Parse(layout, s)
	if err != nil {
		return Zero, err
	}
	return RawDate{t.Year(), int8(t.Month()), int8(t.Day())}, nil
}

func Compare(a, b RawDate) int {
	s := a.Year - b.Year
	if s == 0 {
		s = int(a.Month - b.Month)
	}
	if s == 0 {
		s = int(a.Day - b.Day)
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

func (r RawDate) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", r.Year, r.Month, r.Day)
}

func (r RawDate) Time(location *time.Location) time.Time {
	return time.Date(int(r.Year), time.Month(r.Month), int(r.Day), 0, 0, 0, 0, location)
}

func (r RawDate) Format(format string) string {
	tm := r.Time(time.UTC)
	return tm.Format(format)
}
