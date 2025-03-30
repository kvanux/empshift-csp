package models

type Shift struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Day      string `json:"day"`
	MinStaff int    `json:"min_staff"`
	MaxStaff int    `json:"max_staff"`
}

type Employee struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	MaxShifts   int                 `json:"max_shifts"`
	MinShifts   int                 `json:"min_shifts"`
	Unavailable map[string]struct{} `json:"unavailable"`
	Preferences map[string]struct{} `json:"preferences"`
}

type Schedule struct {
	Assignments map[string][]int `json:"assignments"`
	Locked      map[string][]int `json:"locked"`
	Fitness     float64          `json:"fitness"`
}

// DTOs
type StaffRequest struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Unavailable []string `json:"unavailable"`
	Preferred   []string `json:"preferred"`
	Assigned    []string `json:"assigned"`
	IsOkay      bool     `json:"isOkay"`
}

type ScheduleRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Day         string `json:"day"`
	Assigned    []int  `json:"assigned"`
	Preferred   []int  `json:"preferred"`
	Unavailable []int  `json:"unavailable"`
	IsOkay      bool   `json:"isOkay"`
	IsLocked    bool   `json:"isLocked"`
}

type SchedulePackageRequest struct {
	Staffs    []StaffRequest    `json:"staffs"`
	Schedules []ScheduleRequest `json:"schedules"`
	MinStaff  int               `json:"minStaff"`
	MaxStaff  int               `json:"maxStaff"`
	MinShift  int               `json:"minShift"`
	MaxShift  int               `json:"maxShift"`
}

type ScheduleResponse struct {
	Schedules []ScheduleRequest `json:"schedules"`
	Rating    float64
}
