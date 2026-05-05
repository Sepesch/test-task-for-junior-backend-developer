package task

import (
	"fmt"
	"time"
)

type Recurrence interface {
	Type() string
	Validate() error
	NextOccurrence(after time.Time) (time.Time, bool)
}

type DailyRecurrence struct {
	Interval int `json:"interval"`
}

func (r DailyRecurrence) Type() string { return "daily" }
func (r DailyRecurrence) Validate() error {
	if r.Interval <= 0 {
		return fmt.Errorf("interval must be positive")
	}
	return nil
}
func (r DailyRecurrence) NextOccurrence(after time.Time) (time.Time, bool) {
	next := after.AddDate(0, 0, r.Interval)
	return beginningOfDay(next), true
}

type MonthlyRecurrence struct {
	DayOfMonth int `json:"day_of_month"`
}

func (r MonthlyRecurrence) Type() string { return "monthly" }
func (r MonthlyRecurrence) Validate() error {
	if r.DayOfMonth < 1 || r.DayOfMonth > 30 {
		return fmt.Errorf("day_of_month must be between 1 and 30")
	}
	return nil
}
func (r MonthlyRecurrence) NextOccurrence(after time.Time) (time.Time, bool) {
	candidate := time.Date(after.Year(), after.Month(), r.DayOfMonth, 0, 0, 0, 0, time.UTC)
	if candidate.Before(after) || candidate.Equal(after) {
		candidate = candidate.AddDate(0, 1, 0)
	}
	for candidate.Day() != r.DayOfMonth {
		candidate = candidate.AddDate(0, 1, 0)
	}
	return beginningOfDay(candidate), true
}

type WeeklyRecurrence struct {
	WeekDays []int `json:"week_days"`
}

func (r WeeklyRecurrence) Type() string { return "weekly" }
func (r WeeklyRecurrence) Validate() error {
	if len(r.WeekDays) == 0 {
		return fmt.Errorf("week_days cannot be empty")
	}
	seen := make(map[int]bool)
	for _, d := range r.WeekDays {
		if d < 1 || d > 7 {
			return fmt.Errorf("week day must be between 1 (mon) and 7 (sun)")
		}
		if seen[d] {
			return fmt.Errorf("duplicate week day: %d", d)
		}
		seen[d] = true
	}
	return nil
}
func (r WeeklyRecurrence) NextOccurrence(after time.Time) (time.Time, bool) {
	for d := 1; d <= 7; d++ {
		candidate := after.AddDate(0, 0, d)
		wd := isoWeekday(candidate)
		for _, target := range r.WeekDays {
			if wd == target {
				return beginningOfDay(candidate), true
			}
		}
	}
	return time.Time{}, false
}

type SpecificDatesRecurrence struct {
	Dates []string `json:"dates"`
}

func (r SpecificDatesRecurrence) Type() string { return "specific_dates" }
func (r SpecificDatesRecurrence) Validate() error {
	if len(r.Dates) == 0 {
		return fmt.Errorf("dates cannot be empty")
	}
	for _, d := range r.Dates {
		if _, err := time.Parse("2006-01-02", d); err != nil {
			return fmt.Errorf("invalid date %q: %w", d, err)
		}
	}
	return nil
}
func (r SpecificDatesRecurrence) NextOccurrence(after time.Time) (time.Time, bool) {
	var earliest *time.Time
	for _, d := range r.Dates {
		t, _ := time.Parse("2006-01-02", d)
		t = beginningOfDay(t)
		if t.After(after) {
			if earliest == nil || t.Before(*earliest) {
				earliest = &t
			}
		}
	}
	if earliest == nil {
		return time.Time{}, false
	}
	return *earliest, true
}

type EvenOddRecurrence struct {
	Parity string `json:"parity"`
}

func (r EvenOddRecurrence) Type() string { return "even_odd" }
func (r EvenOddRecurrence) Validate() error {
	if r.Parity != "even" && r.Parity != "odd" {
		return fmt.Errorf("parity must be 'even' or 'odd'")
	}
	return nil
}
func (r EvenOddRecurrence) NextOccurrence(after time.Time) (time.Time, bool) {
	cur := beginningOfDay(after).AddDate(0, 0, 1)
	for i := 0; i < 31; i++ {
		day := cur.Day()
		if r.Parity == "even" && day%2 == 0 {
			return cur, true
		}
		if r.Parity == "odd" && day%2 == 1 {
			return cur, true
		}
		cur = cur.AddDate(0, 0, 1)
	}
	return time.Time{}, false
}

func beginningOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func isoWeekday(t time.Time) int {
	wd := t.Weekday()
	if wd == 0 {
		return 7
	}
	return int(wd)
}