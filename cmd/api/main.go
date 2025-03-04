package main

import (
	"empshift-csp/internal/dummy"
	"empshift-csp/internal/models"
	"fmt"
	"math/rand"
	"sort"
)

func main() {
	fmt.Println("")
	fmt.Println("Program start:")
	fmt.Println("")
	// Input Start
	var employeeNo int
	var minShift int
	var maxShift int
	var daysPerWeek int
	var shiftsPerDay int
	var minStaff int
	var maxStaff int

	fmt.Print("How many days per week: ")
	fmt.Scanf("%v\n", &daysPerWeek)
	fmt.Print("How many shifts per day: ")
	fmt.Scanf("%v\n", &shiftsPerDay)
	fmt.Print("Minimum staffs per shift: ")
	fmt.Scanf("%v\n", &minStaff)
	fmt.Print("Maximum staffs per shift: ")
	fmt.Scanf("%v\n", &maxStaff)
	fmt.Print("How many employees: ")
	fmt.Scanf("%v\n", &employeeNo)
	fmt.Print("Minimum shift per employee: ")
	fmt.Scanf("%v\n", &minShift)
	fmt.Print("Maximum shift per employee: ")
	fmt.Scanf("%v\n", &maxShift)

	var employees []models.Employee
	var shifts []models.Shift

	for i := range employeeNo {
		newEmp := models.Employee{
			ID:        i,
			Name:      dummy.Names[i],
			MaxShifts: maxShift,
			MinShifts: minShift,
		}
		employees = append(employees, newEmp)
	}

	for i := range daysPerWeek {
		for j := range shiftsPerDay {
			newShift := models.Shift{
				ID:       fmt.Sprintf("%s-%s", dummy.DummyShifts[j], dummy.Weekdays[i]),
				Name:     dummy.DummyShifts[j],
				Day:      dummy.Weekdays[i],
				MinStaff: minStaff,
				MaxStaff: maxStaff,
			}
			shifts = append(shifts, newShift)
		}
	}
	// Input End

	// Processing Start
	populationSize := 200
	population := generateRandomSchedules(shifts, employees, populationSize)
	maxGenerations := 1000
	// mutationRate := float32(0.1)

	for gen := 0; gen < maxGenerations; gen++ {
		for i := range population {
			population[i].Fitness = calculateFitness(population[i], shifts, employees)
		}
		parents := selectTopSchedules(population, 0.2)

		var offspring []models.Schedule
		for len(offspring) < populationSize {
			parentA := parents[rand.Intn(len(parents))]
			parentB := parents[rand.Intn(len(parents))]
			child := crossover(parentA, parentB, shifts)
			// child = mutate(child, shifts, employees, mutationRate)
			offspring = append(offspring, child)
		}
		population = offspring
	}
	// Processing End

	// Returning Start
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
	printSchedule(bestSchedule)
	fmt.Printf("%v", &count)
	// Returning End

	fmt.Println("")
	fmt.Println("Program exited.")
	fmt.Println("")
}

func generateRandomSchedules(shifts []models.Shift, employees []models.Employee, populationSize int) []models.Schedule {
	population := make([]models.Schedule, populationSize)
	for i := 0; i < populationSize; i++ {
		schedule := models.Schedule{
			Assignments: make(map[string][]int),
			// Locked:      make(map[string][]int), // Assume locked assignments pre-loaded
		}
		// Copy locked assignments: no way to get input currently, will be done in the future
		// for shiftID, emps := range lockedAssignments {
		//     schedule.Assignments[shiftID] = emps
		// }
		// Assign non-locked shifts
		for _, shift := range shifts {
			// if _, isLocked := schedule.Locked[shift.ID]; !isLocked {}
			numAssigned := 0
			for numAssigned < shift.MinStaff {
				numAssigned = rand.Intn(shift.MaxStaff + 1)
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

func getEmployeeByID(id int, employees []models.Employee) models.Employee {
	for _, emp := range employees {
		if emp.ID == id {
			return emp
		}
	}
	panic(fmt.Sprintf("employee with ID %d not found", id))
}

func calculateFitness(s models.Schedule, shifts []models.Shift, employees []models.Employee) float64 {
	score := 0.0
	for _, shift := range shifts {
		assigned := len(s.Assignments[shift.ID])
		if assigned < shift.MinStaff {
			score -= float64(shift.MinStaff-assigned) * 10
		}
		if assigned > shift.MaxStaff {
			score -= float64(assigned-shift.MaxStaff) * 5
		}
	}

	empShiftCount := make(map[int]int)
	for shiftID, emps := range s.Assignments {
		for _, empID := range emps {
			empShiftCount[empID]++

			emp := getEmployeeByID(empID, employees)
			if _, unavailable := emp.Unavailable[shiftID]; unavailable {
				score -= 15
			}

			if _, preferred := emp.Preferences[shiftID]; preferred {
				score += 3
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

func mutate(schedule models.Schedule, shifts []models.Shift, employees []models.Employee, mutationRate float32) models.Schedule {
	mutated := models.Schedule{
		Assignments: make(map[string][]int),
		Locked:      schedule.Locked,
	}
	for shiftID, emps := range schedule.Assignments {
		if _, isLocked := schedule.Locked[shiftID]; isLocked {
			mutated.Assignments[shiftID] = emps
			continue
		}
		if rand.Float32() < mutationRate {

			newEmps := emps
			if len(emps) > 0 && rand.Float32() < 0.5 {
				newEmps = emps[:len(emps)-1]
			} else {
				newEmp := employees[rand.Intn(len(employees))].ID
				newEmps = append(newEmps, newEmp)
			}
			mutated.Assignments[shiftID] = newEmps
		} else {
			mutated.Assignments[shiftID] = emps
		}
	}
	return mutated
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

func printSchedule(s models.Schedule) {
	for shiftID, emps := range s.Assignments {
		fmt.Printf("Shift %s: Employees %v\n", shiftID, emps)
	}
}
