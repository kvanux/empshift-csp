package models

type Shift struct {
	ID       string // Format: "Name-Weekday"
	Name     string // rand
	Day      string // Monday - Sunday
	MinStaff int    // Min number of employees required for each shift
	MaxStaff int    // Max number of employees allowed for each shift
}

type Employee struct {
	ID          int
	Name        string
	MaxShifts   int                 // Max shifts/week
	MinShifts   int                 // Min shifts/week
	Unavailable map[string]struct{} // Set of shift IDs
	Preferences map[string]struct{} // Set of preferred shift IDs
}

type Schedule struct {
	Assignments map[string][]int // Shift ID → Employee IDs
	Locked      map[string][]int // Locked assignment by users
	Fitness     float64
}
