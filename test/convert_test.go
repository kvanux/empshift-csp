package test

import (
	conversion "empshift-csp/internal/helpers"
	"empshift-csp/internal/models"
	"reflect"
	"testing"
)

func TestConvertStaff(t *testing.T) {
	staffReq := models.StaffRequest{
		ID:          1,
		Name:        "John Doe",
		Unavailable: []string{"Monday", "Tuesday"},
		Preferred:   []string{"Friday"},
	}
	min := 2
	max := 5

	expected := models.Employee{
		ID:          1,
		Name:        "John Doe",
		MinShifts:   min,
		MaxShifts:   max,
		Unavailable: map[string]struct{}{"Monday": {}, "Tuesday": {}},
		Preferences: map[string]struct{}{"Friday": {}},
	}

	result := conversion.ConvertStaff(staffReq, min, max)

	if result.ID != expected.ID || result.Name != expected.Name || result.MinShifts != expected.MinShifts || result.MaxShifts != expected.MaxShifts {
		t.Errorf("ConvertStaff failed: expected %+v, got %+v", expected, result)
	}

	if !reflect.DeepEqual(result.Unavailable, expected.Unavailable) {
		t.Errorf("ConvertStaff failed for Unavailable: expected %+v, got %+v", expected.Unavailable, result.Unavailable)
	}

	if !reflect.DeepEqual(result.Preferences, expected.Preferences) {
		t.Errorf("ConvertStaff failed for Preferences: expected %+v, got %+v", expected.Preferences, result.Preferences)
	}
}

func TestConvertShift(t *testing.T) {
	scheduleReq := models.ScheduleRequest{
		ID:   "456",
		Name: "Morning Shift",
		Day:  "Monday",
	}
	min := 3
	max := 8

	expected := models.Shift{
		ID:       "456",
		Name:     "Morning Shift",
		Day:      "Monday",
		MinStaff: min,
		MaxStaff: max,
	}

	result := conversion.ConvertShift(scheduleReq, min, max)

	if result != expected {
		t.Errorf("ConvertShift failed: expected %+v, got %+v", expected, result)
	}
}
