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
		t    rawdate.RawDate
		args args
		want int
	}{
		{"equal", rawdate.MustNew(2024, 2, 20), args{rawdate.MustNew(2024, 2, 20)}, 0},
		{"larger", rawdate.MustNew(2024, 2, 21), args{rawdate.MustNew(2024, 2, 20)}, 1},
		{"smaller", rawdate.MustNew(2024, 2, 19), args{rawdate.MustNew(2024, 2, 20)}, -1},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if got := tt.t.Compare(tt.args.a); got != tt.want {
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
		t    rawdate.RawDate
		args args
		want string
	}{
		{"easy", rawdate.MustNew(2024, 2, 20), args{rawdate.ISODate}, "2024-02-20"},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {

			if got := tt.t.Format(tt.args.format); got != tt.want {
				t1.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRawDate_String(t1 *testing.T) {

	tests := []struct {
		name string
		t    rawdate.RawDate
		want string
	}{
		{"easy", rawdate.MustNew(2024, 2, 20), "2024-02-20"},
		{"medium", rawdate.MustNew(1, 1, 1), "0001-01-01"},
		{"hard", rawdate.MustNew(0, 1, 1), "0000-01-01"},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {

			if got := tt.t.String(); got != tt.want {
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
		t    rawdate.RawDate
		args args
		want time.Time
	}{
		{
			"easy",
			rawdate.MustNew(now.Year(), int(now.Month()), now.Day()),
			args{time.Local},
			time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {

			if got := tt.t.Time(tt.args.location); !reflect.DeepEqual(got, tt.want) {
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
		{"easy", args{now}, rawdate.MustNew(now.Year(), int(now.Month()), now.Day())},
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
		m int
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
