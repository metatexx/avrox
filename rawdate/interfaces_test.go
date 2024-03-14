package rawdate_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	testifyrequire "github.com/stretchr/testify/require"

	"github.com/metatexx/avrox/rawdate"
)

func TestDate_MarshalText(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Name     string
		Date     rawdate.RawDate
		Expected string
	}

	cases := []testCase{
		{Name: "Remote past", Date: rawdate.RawDate{Year0: 1997 - 1, Month0: int8(time.July - 1), Day0: 15 - 1}, Expected: "1997-07-15"},
		{Name: "Recent past", Date: rawdate.RawDate{Year0: 2021 - 1, Month0: int8(time.February - 1), Day0: 20 - 1}, Expected: "2021-02-20"},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			asBytes, err := tc.Date.MarshalText()
			assert.Nil(err)
			assert.Equal(tc.Expected, string(asBytes))
		})
	}
}

func TestDate_MarshalJSON(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Name     string
		Date     *rawdate.RawDate
		Expected string
	}

	cases := []testCase{
		{Name: "Remote past", Date: &rawdate.RawDate{Year0: 1997 - 1, Month0: int8(time.July - 1), Day0: 15 - 1}, Expected: `"1997-07-15"`},
		{Name: "Recent past", Date: &rawdate.RawDate{Year0: 2021 - 1, Month0: int8(time.February - 1), Day0: 20 - 1}, Expected: `"2021-02-20"`},
		{Name: "Unset", Date: nil, Expected: "null"},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			asBytes, err := json.Marshal(tc.Date)
			assert.Nil(err)
			assert.Equal(tc.Expected, string(asBytes))
		})
	}
}

func TestDate_UnmarshalText(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Input []byte
		Date  rawdate.RawDate
		Error string
	}

	cases := []testCase{
		{Input: []byte(`x`), Error: `parsing time "x" as "2006-01-02": cannot parse "x" as "2006"`},
		{Input: []byte(`10`), Error: `parsing time "10" as "2006-01-02": cannot parse "10" as "2006"`},
		{Input: []byte("01/26/2018"), Error: `parsing time "01/26/2018" as "2006-01-02": cannot parse "01/26/2018" as "2006"`},
		{Input: []byte("1997-07-15"), Date: rawdate.RawDate{Year0: 1997 - 1, Month0: int8(time.July - 1), Day0: 15 - 1}},
		{Input: []byte("2021-02-20"), Date: rawdate.RawDate{Year0: 2021 - 1, Month0: int8(time.February - 1), Day0: 20 - 1}},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(string(tc.Input), func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			d := rawdate.RawDate{}
			err := d.UnmarshalText(tc.Input)
			if err != nil {
				assert.Equal(tc.Error, fmt.Sprintf("%v", err))
				assert.Equal(rawdate.RawDate{}, d)
			} else {
				assert.Equal("", tc.Error)
				assert.Equal(tc.Date, d)
			}
		})
	}
}

func TestDate_UnmarshalJSON(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Input []byte
		Date  rawdate.RawDate
		Error string
	}

	cases := []testCase{
		{Input: []byte(`x`), Error: "invalid character 'x' looking for beginning of value"},
		{Input: []byte(`10`), Error: "json: cannot unmarshal number into Go value of type string"},
		{Input: []byte(`"abc"`), Error: `parsing time "abc" as "2006-01-02": cannot parse "abc" as "2006"`},
		{Input: []byte(`"01/26/2018"`), Error: `parsing time "01/26/2018" as "2006-01-02": cannot parse "01/26/2018" as "2006"`},
		{Input: []byte(`"1997-07-15"`), Date: rawdate.RawDate{Year0: 1997 - 1, Month0: int8(time.July - 1), Day0: 15 - 1}},
		{Input: []byte(`"2021-02-20"`), Date: rawdate.RawDate{Year0: 2021 - 1, Month0: int8(time.February - 1), Day0: 20 - 1}},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(string(tc.Input), func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			d := rawdate.RawDate{}
			err := json.Unmarshal(tc.Input, &d)
			if err != nil {
				assert.Equal(tc.Error, fmt.Sprintf("%v", err))
				assert.Equal(rawdate.RawDate{}, d)
			} else {
				assert.Equal("", tc.Error)
				assert.Equal(tc.Date, d)
			}
		})
	}
}

func TestDate_Scan(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	// Wrong type
	d := rawdate.RawDate{}
	err := d.Scan(1)
	assert.NotNil(err)
	assert.Equal("wrong type to scan as RawDate: type=int", fmt.Sprintf("%v", err))
	assert.Equal(rawdate.RawDate{}, d)

	// Time but not date
	d = rawdate.RawDate{}
	tz, err := time.LoadLocation("Europe/Berlin")
	assert.Nil(err)
	src := time.Date(2001, time.August, 4, 11, 10, 55, 0, tz)
	err = d.Scan(src)
	assert.NotNil(err)
	assert.Equal("timestamp contains time information: 2001-08-04 11:10:55 +0200 CEST", fmt.Sprintf("%v", err))
	assert.Equal(rawdate.RawDate{}, d)

	// Happy path
	d = rawdate.RawDate{}
	src = time.Date(1991, time.April, 26, 0, 0, 0, 0, time.Local)
	err = d.Scan(src)
	assert.Nil(err)
	expected := rawdate.RawDate{Year0: 1991 - 1, Month0: int8(time.April - 1), Day0: 26 - 1}
	assert.Equal(expected, d)
}

func TestDate_Value(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	d := rawdate.RawDate{Year0: 1991 - 1, Month0: int8(time.April - 1), Day0: 26 - 1}
	v, err := d.Value()
	assert.Nil(err)
	expected := time.Date(1991, time.April, 26, 0, 0, 0, 0, time.Local)
	assert.Equal(expected, v)
}

func TestDate_String(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Date     rawdate.RawDate
		Expected string
	}

	cases := []testCase{
		{Date: rawdate.RawDate{Year0: 2021 - 1, Month0: int8(time.May - 1), Day0: 11 - 1}, Expected: "2021-05-11"},
		{Date: rawdate.RawDate{Year0: 2024 - 1, Month0: int8(time.January - 1), Day0: 31 - 1}, Expected: "2024-01-31"},
		{Date: rawdate.RawDate{Year0: 1999 - 1, Month0: int8(time.December - 1), Day0: 24 - 1}, Expected: "1999-12-24"},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(tc.Expected, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			assert.Equal(tc.Expected, tc.Date.String())
		})
	}
}

func TestDate_GoString(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Date     rawdate.RawDate
		Expected string
	}

	cases := []testCase{
		{Date: rawdate.RawDate{Year0: 2021 - 1, Month0: int8(time.May - 1), Day0: 11 - 1}, Expected: "rawdate.MustNew(2021, time.May, 11)"},
		{Date: rawdate.RawDate{Year0: 2024 - 1, Month0: int8(time.January - 1), Day0: 31 - 1}, Expected: "rawdate.MustNew(2024, time.January, 31)"},
		{Date: rawdate.RawDate{Year0: 1999 - 1, Month0: int8(time.December - 1), Day0: 24 - 1}, Expected: "rawdate.MustNew(1999, time.December, 24)"},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(tc.Expected, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			assert.Equal(tc.Expected, tc.Date.GoString())
		})
	}
}
