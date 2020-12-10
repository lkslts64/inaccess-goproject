package ptask

import (
	"errors"
	"time"
)

const (
	hour  = "1h"
	day   = "1d"
	month = "1mo"
	year  = "1y"
)

// List returns the period task list between (t1,t2) with period `period`.
// Period should be one of {1h,1d,1mo,1y}.
func List(t1, t2 time.Time, period string) ([]time.Time, error) {
	if t1.After(t2) {
		return nil, errors.New("t1 comes after t2")
	}
	pt, err := newPeriodTask(period, t1)
	if err != nil {
		return nil, err
	}
	ptlist := []time.Time{}
	// if t1 matches the start of the period , append it in the list.
	if pt.matchesInvocationPeriod(t1) {
		ptlist = append(ptlist, t1)
	}
	next := pt.nextInvocationTime()
	for ; next.Before(t2); next = pt.nextInvocationTime() {
		ptlist = append(ptlist, next)
	}
	return ptlist, nil
}

type periodTask struct {
	period string
	t      time.Time
}

func newPeriodTask(period string, t time.Time) (*periodTask, error) {
	switch period {
	case hour, day, month, year:
		return &periodTask{period: period, t: t}, nil
	default:
		return nil, errors.New("uknown period")
	}
}

// we find the next invocation by adding the period to the previous one.
func (pt *periodTask) nextInvocationTime() time.Time {
	t := pt.t
	switch pt.period {
	case hour:
		//`t` will always be equal to `prev` except maybe from the first invocation
		// of this func.
		prev := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
		pt.t = prev.Add(time.Hour)
	case day:
		prev := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		pt.t = prev.AddDate(0, 0, 1)
	case month:
		prev := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		pt.t = prev.AddDate(0, 1, 0)
	case year:
		prev := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
		pt.t = prev.AddDate(1, 0, 0)
	default:
		panic(pt.period)
	}
	return pt.t
}

func (pt *periodTask) matchesInvocationPeriod(t time.Time) bool {
	switch pt.period {
	case hour:
		return isStartOfHour(t)
	case day:
		return isStartOfDay(t)
	case month:
		return isStartOfMonth(t)
	case year:
		return isStartOfYear(t)
	default:
		panic(pt.period)
	}
}

func isStartOfHour(t time.Time) bool {
	return t.Minute() == 0 && t.Second() == 0 && t.Nanosecond() == 0
}

func isStartOfDay(t time.Time) bool {
	return isStartOfHour(t) && t.Hour() == 0
}

func isStartOfMonth(t time.Time) bool {
	return isStartOfDay(t) && t.Day() == 1
}

func isStartOfYear(t time.Time) bool {
	return isStartOfMonth(t) && t.Month() == 1
}
