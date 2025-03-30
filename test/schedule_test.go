package test

import (
	schedule "empshift-csp/internal/core"
	"empshift-csp/internal/models"
	"testing"
)

func TestComputeSchedule(t *testing.T) {
	req := models.SchedulePackageRequest{
		MinShift: 2,
		MaxShift: 5,
		MinStaff: 1,
		MaxStaff: 3,
		Staffs: []models.StaffRequest{
			{ID: 1, Name: "Alice", Unavailable: []string{"Monday"}, Preferred: []string{"Friday"}},
			{ID: 2, Name: "Bob", Unavailable: []string{"Tuesday"}, Preferred: []string{"Monday"}},
			{ID: 3, Name: "Chuck", Unavailable: []string{"Tuesday"}, Preferred: []string{"Monday"}},
		},
		Schedules: []models.ScheduleRequest{
			{ID: "shift1", Name: "Morning", Day: "Monday", IsLocked: false, Assigned: []int{}},
			{ID: "shift2", Name: "Evening", Day: "Tuesday", IsLocked: false, Assigned: []int{}},
		},
	}

	result, err := schedule.ComputeSchedule(req)
	if err != nil {
		t.Fatalf("ComputeSchedule returned an error: %v", err)
	}

	if len(result.Schedules) != len(req.Schedules) {
		t.Errorf("Expected %d schedules, got %d", len(req.Schedules), len(result.Schedules))
	}

	// if result.Rating == 0 {
	// 	t.Errorf("Expected a non-zero rating, got %f", result.Rating)
	// }

	for _, sched := range result.Schedules {
		if len(sched.Assigned) < req.MinStaff || len(sched.Assigned) > req.MaxStaff {
			t.Errorf("Schedule %s has invalid number of assignments: %d", sched.ID, len(sched.Assigned))
		}
	}
}

func TestGenerateRandomSchedules(t *testing.T) {
	shifts := []models.Shift{
		{ID: "shift1", Name: "Morning", Day: "Monday", MinStaff: 1, MaxStaff: 3},
		{ID: "shift2", Name: "Evening", Day: "Tuesday", MinStaff: 1, MaxStaff: 3},
	}
	employees := []models.Employee{
		{ID: 1, Name: "Alice", MinShifts: 1, MaxShifts: 7},
		{ID: 2, Name: "Bob", MinShifts: 1, MaxShifts: 7},
		{ID: 3, Name: "Chuck", MinShifts: 1, MaxShifts: 7},
	}
	current := models.Schedule{
		Assignments: make(map[string][]int),
		Locked:      make(map[string][]int),
	}

	population := schedule.GenerateRandomSchedules(shifts, employees, 10, current)
	if len(population) != 10 {
		t.Errorf("Expected population size of 10, got %d", len(population))
	}

	for _, sched := range population {
		for shiftID, assigned := range sched.Assignments {
			if len(assigned) < shifts[0].MinStaff || len(assigned) > shifts[0].MaxStaff {
				t.Errorf("Shift %s has invalid number of assignments: %d", shiftID, len(assigned))
			}
		}
	}
}
