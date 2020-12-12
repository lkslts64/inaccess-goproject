package ptask

import (
	"testing"
	"time"
)

const layout = "20060102T150405Z"
const tz = "Europe/Athens"
const limit = 1 << 10

func TestPtlist(t *testing.T) {

	parseTimePanicOnErr := func(value string) time.Time {
		t, err := parseTime(layout, value, tz)
		if err != nil {
			panic(err)
		}
		return t
	}

	tests := []struct {
		t1   time.Time
		t2   time.Time
		p    string
		want int
	}{
		{t1: parseTimePanicOnErr("20200301T000000Z"), t2: parseTimePanicOnErr("20200301T235959Z"), p: "1h", want: 24},
		{t1: parseTimePanicOnErr("20200301T000000Z"), t2: parseTimePanicOnErr("20200301T235959Z"), p: "1d", want: 1},
		{t1: parseTimePanicOnErr("20200301T000000Z"), t2: parseTimePanicOnErr("20200301T235959Z"), p: "1mo", want: 0},
		{t1: parseTimePanicOnErr("19900301T000000Z"), t2: parseTimePanicOnErr("20200301T235959Z"), p: "1y", want: 30},
		{t1: parseTimePanicOnErr("20200301T000000Z"), t2: parseTimePanicOnErr("20200302T142959Z"), p: "1h", want: 39},
		{t1: parseTimePanicOnErr("20200101T000000Z"), t2: parseTimePanicOnErr("20280101T000000Z"), p: "1y", want: 8},
		{t1: parseTimePanicOnErr("17000301T000000Z"), t2: parseTimePanicOnErr("20200301T235959Z"), p: "1h", want: 1025},
	}

	for _, tt := range tests {
		list, err := List(tt.t1, tt.t2, tt.p, limit)
		if err != nil {
			if err == errLimit && tt.want > limit {
				//we wanted to get a limitErr
				continue
			}
			t.Error(err)
		}
		if len(list) != tt.want {
			t.Errorf("want : %d, got : %d", tt.want, len(list))
		}
	}

}

// parses `value` based on `layout` and returns a time.Time whose's location is
// set as `name`.
func parseTime(layout, value, name string) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, err
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}
