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

	populationSize := 200
	population := generateRandomSchedules(shifts, employees, populationSize)
	maxGenerations := 1000

	for gen := 0; gen < maxGenerations; gen++ {
		for i := range population {
			population[i].Fitness = calculateFitness(population[i], shifts, employees)
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
		population[i].Fitness = calculateFitness(population[i], shifts, employees)
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
	printSchedule(bestSchedule)
	fmt.Printf("%v", &count)

	fmt.Println("")
	fmt.Println("Program exited.")
	fmt.Println("")
}

func generateRandomSchedules(shifts []models.Shift, employees []models.Employee, populationSize int) []models.Schedule {
	population := make([]models.Schedule, populationSize)
	for i := 0; i < populationSize; i++ {
		schedule := models.Schedule{
			Assignments: make(map[string][]int),
			Locked:      make(map[string][]int),
		}
		// for shiftID, emps := range lockedAssignments {
		//     schedule.Assignments[shiftID] = emps
		// }
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

func printSchedule(s models.Schedule) {
	for shiftID, emps := range s.Assignments {
		fmt.Printf("Shift %s: Employees %v\n", shiftID, emps)
	}
}
