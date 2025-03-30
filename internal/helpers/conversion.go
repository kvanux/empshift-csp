package conversion

import (
	"empshift-csp/internal/models"
)

func ConvertStaff(staffReq models.StaffRequest, min int, max int) models.Employee {
	return models.Employee{
		ID:          staffReq.ID,
		Name:        staffReq.Name,
		MinShifts:   min,
		MaxShifts:   max,
		Unavailable: sliceToSet(staffReq.Unavailable),
		Preferences: sliceToSet(staffReq.Preferred),
	}
}

func ConvertShift(scheduleReq models.ScheduleRequest, min int, max int) models.Shift {
	return models.Shift{
		ID:       scheduleReq.ID,
		Name:     scheduleReq.Name,
		Day:      scheduleReq.Day,
		MinStaff: min,
		MaxStaff: max,
	}
}

func sliceToSet(slice []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, item := range slice {
		set[item] = struct{}{}
	}
	return set
}
