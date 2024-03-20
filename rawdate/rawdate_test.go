package rawdate_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/metatexx/avrox/rawdate"
)

func TestCompare(t *testing.T) {
	type args struct {
		a rawdate.RawDate
		b rawdate.RawDate
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"equal", args{rawdate.MustNew(2024, 2, 20), rawdate.MustNew(2024, 2, 20)}, 0},
		{"larger", args{rawdate.MustNew(2024, 2, 21), rawdate.MustNew(2024, 2, 20)}, 1},
		{"smaller", args{rawdate.MustNew(2024, 2, 19), rawdate.MustNew(2024, 2, 20)}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rawdate.Compare(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_Compare(t1 *testing.T) {
	type args struct {
		a rawdate.RawDate
	}
	tests := []struct {
		name string
		base rawdate.RawDate
		args args
		want int
	}{
		{"equal", rawdate.MustNew(2024, 2, 20), args{rawdate.MustNew(2024, 2, 20)}, 0},
		{"larger", rawdate.MustNew(2024, 2, 21), args{rawdate.MustNew(2024, 2, 20)}, 1},
		{"smaller", rawdate.MustNew(2024, 2, 19), args{rawdate.MustNew(2024, 2, 20)}, -1},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if got := tt.base.Compare(tt.args.a); got != tt.want {
				t1.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_Equal(t1 *testing.T) {
	type args struct {
		a rawdate.RawDate
	}
	tests := []struct {
		name string
		base rawdate.RawDate
		args args
		want bool
	}{
		{"equal", rawdate.MustNew(2024, 2, 20), args{rawdate.MustNew(2024, 2, 20)}, true},
		{"before", rawdate.MustNew(2024, 2, 20), args{rawdate.MustNew(2024, 2, 19)}, false},
		{"after", rawdate.MustNew(2024, 2, 20), args{rawdate.MustNew(2024, 2, 21)}, false},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if got := tt.base.Equal(tt.args.a); got != tt.want {
				t1.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_Before(t1 *testing.T) {
	type args struct {
		a rawdate.RawDate
	}
	tests := []struct {
		name string
		base rawdate.RawDate
		args args
		want bool
	}{
		{"equal", rawdate.MustNew(2024, 2, 20), args{rawdate.MustNew(2024, 2, 20)}, false},
		{"before", rawdate.MustNew(2024, 2, 20), args{rawdate.MustNew(2024, 2, 19)}, false},
		{"after", rawdate.MustNew(2024, 2, 20), args{rawdate.MustNew(2024, 2, 21)}, true},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if got := tt.base.Before(tt.args.a); got != tt.want {
				t1.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_After(t1 *testing.T) {
	type args struct {
		a rawdate.RawDate
	}
	tests := []struct {
		name string
		base rawdate.RawDate
		args args
		want bool
	}{
		{"equal", rawdate.MustNew(2024, 2, 20), args{rawdate.MustNew(2024, 2, 20)}, false},
		{"before", rawdate.MustNew(2024, 2, 20), args{rawdate.MustNew(2024, 2, 19)}, true},
		{"after", rawdate.MustNew(2024, 2, 20), args{rawdate.MustNew(2024, 2, 21)}, false},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if got := tt.base.After(tt.args.a); got != tt.want {
				t1.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_Format(t1 *testing.T) {

	type args struct {
		format string
	}
	tests := []struct {
		name string
		base rawdate.RawDate
		args args
		want string
	}{
		{"easy", rawdate.MustNew(2024, 2, 20), args{rawdate.ISODate}, "2024-02-20"},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {

			if got := tt.base.Format(tt.args.format); got != tt.want {
				t1.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_String(t1 *testing.T) {

	tests := []struct {
		name string
		base rawdate.RawDate
		want string
	}{
		{"easy", rawdate.MustNew(2024, 2, 20), "2024-02-20"},
		{"medium", rawdate.MustNew(1, 1, 1), "0001-01-01"},
		{"hard", rawdate.MustNew(0, 1, 1), "0000-01-01"},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {

			if got := tt.base.String(); got != tt.want {
				t1.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_Time(t1 *testing.T) {
	now := time.Now()
	type args struct {
		location *time.Location
	}
	tests := []struct {
		name string
		base rawdate.RawDate
		args args
		want time.Time
	}{
		{
			"easy",
			rawdate.MustNew(now.Year(), now.Month(), now.Day()),
			args{time.Local},
			time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {

			if got := tt.base.Time(tt.args.location); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromTime(t *testing.T) {
	now := time.Now()
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want rawdate.RawDate
	}{
		{"easy", args{now}, rawdate.MustNew(now.Year(), now.Month(), now.Day())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rawdate.FromTime(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		layout string
		s      string
	}
	tests := []struct {
		name    string
		args    args
		want    rawdate.RawDate
		wantErr bool
	}{
		{"easy", args{rawdate.ISODate, "2024-02-20"}, rawdate.MustNew(2024, 2, 20), false},
		{"twisted", args{rawdate.ISODate, "2024-02-00"}, rawdate.Zero, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rawdate.Parse(tt.args.layout, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		y int
		m time.Month
		d int
	}
	tests := []struct {
		name    string
		args    args
		want    rawdate.RawDate
		wantErr bool
	}{
		{"easy", args{2024, 2, 20}, rawdate.MustNew(2024, 2, 20), false},
		{"bloede", args{2024, 2, 0}, rawdate.Zero, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rawdate.New(tt.args.y, tt.args.m, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_IsZero(t *testing.T) {
	type fields struct {
		Year  int
		Month int8
		Day   int8
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"zero", fields{Year: 0, Month: 0, Day: 0}, true},
		{"not zero", fields{Year: 1, Month: 1, Day: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rawdate.RawDate{
				Year0:  tt.fields.Year,
				Month0: tt.fields.Month,
				Day0:   tt.fields.Day,
			}
			if got := r.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_AddDate(t *testing.T) {
	type args struct {
		years  int
		months int
		days   int
	}
	tests := []struct {
		name string
		base rawdate.RawDate
		args args
		want rawdate.RawDate
	}{
		{"zero", rawdate.MustNew(2024, 2, 20),
			args{
				years:  0,
				months: 0,
				days:   0,
			}, rawdate.MustNew(2024, 2, 20)},
		{"one day", rawdate.MustNew(2024, 2, 20),
			args{
				years:  0,
				months: 0,
				days:   1,
			}, rawdate.MustNew(2024, 2, 21)},
		{"one month", rawdate.MustNew(2024, 2, 20),
			args{
				years:  0,
				months: 1,
				days:   0,
			}, rawdate.MustNew(2024, 3, 20)},
		{"evil date", rawdate.MustNew(2024, 2, 29),
			args{
				years:  1,
				months: 0,
				days:   0,
			}, rawdate.MustNew(2025, 3, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.base.AddDate(tt.args.years, tt.args.months, tt.args.days); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_Weekday(t *testing.T) {
	tests := []struct {
		name string
		base rawdate.RawDate
		want time.Weekday
	}{
		{"monday", rawdate.MustNew(2024, 2, 19), time.Monday},
		{"sunday", rawdate.MustNew(2024, 2, 18), time.Sunday},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.base.Weekday(); got != tt.want {
				t.Errorf("Weekday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_MonthStart(t *testing.T) {
	tests := []struct {
		name string
		base rawdate.RawDate
		want rawdate.RawDate
	}{
		{"march 2024", rawdate.MustNew(2024, time.March, 11), rawdate.MustNew(2024, time.March, 1)},
		{"february 2024", rawdate.MustNew(2024, time.February, 13), rawdate.MustNew(2024, time.February, 1)},
		{"february 2023", rawdate.MustNew(2023, time.February, 15), rawdate.MustNew(2023, time.February, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.base.MonthStart(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MonthStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_MonthEnd(t *testing.T) {
	tests := []struct {
		name string
		base rawdate.RawDate
		want rawdate.RawDate
	}{
		{"march 2024", rawdate.MustNew(2024, time.March, 11), rawdate.MustNew(2024, time.March, 31)},
		{"april 2024", rawdate.MustNew(2024, time.April, 11), rawdate.MustNew(2024, time.April, 30)},
		{"february 2024", rawdate.MustNew(2024, time.February, 13), rawdate.MustNew(2024, time.February, 29)},
		{"february 2023", rawdate.MustNew(2023, time.February, 15), rawdate.MustNew(2023, time.February, 28)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.base.MonthEnd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MonthEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_NextWeekday(t *testing.T) {
	tests := []struct {
		name    string
		base    rawdate.RawDate
		weekday time.Weekday
		orToday bool
		want    rawdate.RawDate
	}{
		{"1. march 2024 false", rawdate.MustNew(2024, time.March, 1), time.Monday, false, rawdate.MustNew(2024, time.March, 4)},
		{"1. march 2024 false", rawdate.MustNew(2024, time.March, 1), time.Friday, false, rawdate.MustNew(2024, time.March, 8)},
		{"1. march 2024 true", rawdate.MustNew(2024, time.March, 1), time.Friday, true, rawdate.MustNew(2024, time.March, 1)},
		{"15. april 2024 true", rawdate.MustNew(2024, time.April, 15), time.Sunday, true, rawdate.MustNew(2024, time.April, 21)},
		{"15. april 2024 false", rawdate.MustNew(2024, time.April, 15), time.Sunday, false, rawdate.MustNew(2024, time.April, 21)},
		{"16. april 2024 true", rawdate.MustNew(2024, time.April, 15), time.Thursday, true, rawdate.MustNew(2024, time.April, 18)},
		{"16. april 2024 false", rawdate.MustNew(2024, time.April, 15), time.Thursday, false, rawdate.MustNew(2024, time.April, 18)},
		{"21. april 2024 true", rawdate.MustNew(2024, time.April, 21), time.Friday, true, rawdate.MustNew(2024, time.April, 26)},
		{"21. april 2024 false", rawdate.MustNew(2024, time.April, 21), time.Friday, false, rawdate.MustNew(2024, time.April, 26)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.base.NextWeekday(tt.weekday, tt.orToday); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NextWeekday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_PreviousWeekday(t *testing.T) {
	tests := []struct {
		name    string
		base    rawdate.RawDate
		weekday time.Weekday
		orToday bool
		want    rawdate.RawDate
	}{
		{"1. march 2024 false", rawdate.MustNew(2024, time.March, 1), time.Monday, false, rawdate.MustNew(2024, time.February, 26)},
		{"1. march 2023 false", rawdate.MustNew(2023, time.March, 1), time.Friday, false, rawdate.MustNew(2023, time.February, 24)},
		{"1. march 2024 true", rawdate.MustNew(2024, time.March, 1), time.Friday, true, rawdate.MustNew(2024, time.March, 1)},
		{"1. march 2024 true", rawdate.MustNew(2024, time.March, 1), time.Friday, false, rawdate.MustNew(2024, time.February, 23)},
		{"15. april 2024 true", rawdate.MustNew(2024, time.April, 15), time.Sunday, true, rawdate.MustNew(2024, time.April, 14)},
		{"15. april 2024 false", rawdate.MustNew(2024, time.April, 15), time.Sunday, false, rawdate.MustNew(2024, time.April, 14)},
		{"16. april 2024 true", rawdate.MustNew(2024, time.April, 15), time.Thursday, true, rawdate.MustNew(2024, time.April, 11)},
		{"16. april 2024 false", rawdate.MustNew(2024, time.April, 15), time.Thursday, false, rawdate.MustNew(2024, time.April, 11)},
		{"21. april 2024 true", rawdate.MustNew(2024, time.April, 21), time.Sunday, true, rawdate.MustNew(2024, time.April, 21)},
		{"21. april 2024 false", rawdate.MustNew(2024, time.April, 21), time.Sunday, false, rawdate.MustNew(2024, time.April, 14)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.base.PreviousWeekday(tt.weekday, tt.orToday); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PreviousWeekday() = %v, want %v", got, tt.want)
			}
		})
	}
}
