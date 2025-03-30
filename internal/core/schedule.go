package schedule

import (
	conversion "empshift-csp/internal/helpers"
	"empshift-csp/internal/models"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func ComputeSchedule(req models.SchedulePackageRequest) (models.ScheduleResponse, error) {
	start := time.Now()
	defer func() {
		fmt.Printf("Schedule generated in %v\n", time.Since(start))
	}()

	minShift := req.MinShift
	maxShift := req.MaxShift
	minStaff := req.MinStaff
	maxStaff := req.MaxStaff

	employees := make([]models.Employee, 0, len(req.Staffs))
	shifts := make([]models.Shift, 0, len(req.Schedules))
	current := models.Schedule{
		Assignments: make(map[string][]int),
		Locked:      make(map[string][]int),
	}

	for _, staff := range req.Staffs {
		employees = append(employees, conversion.ConvertStaff(staff, minShift, maxShift))
	}
	for _, shift := range req.Schedules {
		shifts = append(shifts, conversion.ConvertShift(shift, minStaff, maxStaff))
		if shift.IsLocked {
			current.Locked[shift.ID] = shift.Assigned
		} else {
			current.Assignments[shift.ID] = shift.Assigned
		}
	}

	populationSize := 200
	population := GenerateRandomSchedules(shifts, employees, populationSize, current)
	maxGenerations := 1500

	for gen := 0; gen < maxGenerations; gen++ {
		for i := range population {
			population[i].Fitness = CalculateFitness(population[i], shifts, employees)
		}
		parents := selectTopSchedules(population, 0.1)

		var offspring []models.Schedule
		for len(offspring) < populationSize {
			parentA := parents[rand.Intn(len(parents))]
			parentB := parents[rand.Intn(len(parents))]
			child := crossover(parentA, parentB, shifts)
			offspring = append(offspring, child)
		}
		population = offspring
	}

	for i := range population {
		population[i].Fitness = CalculateFitness(population[i], shifts, employees)
	}
	bestSchedule := getBestSchedule(population)
	count := make(map[int]int)
	for j := range employees {
		count[j] = 0
	}
	for _, shift := range bestSchedule.Assignments {
		for _, value := range shift {
			if _, exists := count[value]; exists {
				count[value]++
			}
		}
	}

	var result models.ScheduleResponse
	resultSchedules := req.Schedules

	for i := range resultSchedules {
		for key, values := range bestSchedule.Assignments {
			if resultSchedules[i].ID == key {
				resultSchedules[i].Assigned = values
			}
		}
	}

	result.Schedules = resultSchedules
	result.Rating = bestSchedule.Fitness

	return result, nil
}

func GenerateRandomSchedules(shifts []models.Shift, employees []models.Employee, populationSize int, current models.Schedule) []models.Schedule {
	population := make([]models.Schedule, populationSize)

	for i := 0; i < populationSize; i++ {
		schedule := models.Schedule{
			Assignments: make(map[string][]int),
			Locked:      make(map[string][]int),
		}

		for shiftID, assignedEmps := range current.Locked {
			schedule.Locked[shiftID] = assignedEmps
			schedule.Assignments[shiftID] = assignedEmps
		}

		for _, shift := range shifts {
			if _, isLocked := schedule.Locked[shift.ID]; isLocked {
				continue
			}

			numAssigned := 0
			for numAssigned < shift.MinStaff {
				numAssigned = rand.Intn(shift.MaxStaff + 1)
			}

			if numAssigned > len(employees) {
				numAssigned = len(employees)
			}

			shuffledEmps := shuffleEmployees(employees)
			employeeIDs := make([]int, 0)
			for _, emp := range shuffledEmps[:numAssigned] {
				employeeIDs = append(employeeIDs, emp.ID)
			}
			schedule.Assignments[shift.ID] = employeeIDs
		}
		population[i] = schedule
	}

	return population
}

func shuffleEmployees(employees []models.Employee) []models.Employee {
	shuffled := make([]models.Employee, len(employees))
	copy(shuffled, employees)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

func CalculateFitness(s models.Schedule, shifts []models.Shift, employees []models.Employee) float64 {
	score := 100.0

	empMap := make(map[int]models.Employee)
	for _, emp := range employees {
		empMap[emp.ID] = emp
	}

	empShiftCount := make(map[int]int)
	empDayShifts := make(map[int]map[string]int)

	for _, shift := range shifts {
		assigned := s.Assignments[shift.ID]
		unique := make(map[int]struct{})

		for _, empID := range assigned {
			if _, exists := unique[empID]; exists {
				score -= 20
			}
			unique[empID] = struct{}{}

			empShiftCount[empID]++

			if empDayShifts[empID] == nil {
				empDayShifts[empID] = make(map[string]int)
			}
			empDayShifts[empID][shift.Day]++

			emp, exist := empMap[empID]
			if exist {
				if _, notAvailable := emp.Unavailable[shift.Day]; notAvailable {
					score -= 40
				}
				if _, prefers := emp.Preferences[shift.Day]; prefers {
					score += 5
				}
			}
		}

		if len(unique) < shift.MinStaff {
			score -= float64(shift.MinStaff-len(unique)) * 10
		}
		if len(unique) > shift.MaxStaff {
			score -= float64(len(unique)-shift.MaxStaff) * 5
		}
	}

	for _, days := range empDayShifts {
		for _, count := range days {
			if count > 1 {
				penalty := (count - 1) * 10
				score -= float64(penalty)
			}
		}
	}

	for _, emp := range employees {
		count := empShiftCount[emp.ID]
		if count > emp.MaxShifts {
			score -= float64(count-emp.MaxShifts) * 8
		}
		if count < emp.MinShifts {
			score -= float64(emp.MinShifts-count) * 8
		}
	}
	return score
}

func selectTopSchedules(population []models.Schedule, topPercent float64) []models.Schedule {
	sort.Slice(population, func(i, j int) bool {
		return population[i].Fitness > population[j].Fitness
	})
	topN := int(float64(len(population)) * topPercent)
	return population[:topN]
}

func crossover(parentA, parentB models.Schedule, shifts []models.Shift) models.Schedule {
	child := models.Schedule{
		Assignments: make(map[string][]int),
		Locked:      parentA.Locked,
	}
	for _, shift := range shifts {
		if _, isLocked := child.Locked[shift.ID]; isLocked {
			child.Assignments[shift.ID] = parentA.Locked[shift.ID]
		} else {
			if rand.Float32() < 0.5 {
				child.Assignments[shift.ID] = parentA.Assignments[shift.ID]
			} else {
				child.Assignments[shift.ID] = parentB.Assignments[shift.ID]
			}
		}
	}
	return child
}

func getBestSchedule(population []models.Schedule) models.Schedule {
	if len(population) == 0 {
		panic("population is empty")
	}
	best := population[0]
	for _, s := range population {
		if s.Fitness > best.Fitness {
			best = s
		}
	}
	return best
}
