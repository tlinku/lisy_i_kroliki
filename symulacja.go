package main

import (
	"math/rand"
)

type World struct {
	Grid   [][]Organism
	Width  int
	Height int
	Turn   int
	nextID int
}

func NewWorld(width, height int) *World {
	grid := make([][]Organism, height)
	for i := range grid {
		grid[i] = make([]Organism, width)
	}
	return &World{
		Grid:   grid,
		Width:  width,
		Height: height,
		Turn:   0,
		nextID: 1,
	}
}

func (w *World) IsValidPosition(x, y int) bool {
	return x >= 0 && x < w.Width && y >= 0 && y < w.Height
}

func (w *World) IsEmpty(x, y int) bool {
	return w.IsValidPosition(x, y) && w.Grid[y][x] == nil
}

func (w *World) GetOrganism(x, y int) Organism {
	if !w.IsValidPosition(x, y) {
		return nil
	}
	return w.Grid[y][x]
}

func (w *World) PlaceOrganism(organism Organism) bool {
	x, y := organism.GetPosition()
	if !w.IsEmpty(x, y) {
		return false
	}
	w.Grid[y][x] = organism
	return true
}

func (w *World) RemoveOrganism(x, y int) {
	if w.IsValidPosition(x, y) {
		w.Grid[y][x] = nil
	}
}

func (w *World) MoveOrganism(fromX, fromY, toX, toY int) bool {
	if !w.IsValidPosition(fromX, fromY) || !w.IsEmpty(toX, toY) {
		return false
	}

	organism := w.Grid[fromY][fromX]
	if organism == nil || !organism.CanMove() {
		return false
	}

	w.Grid[fromY][fromX] = nil
	w.Grid[toY][toX] = organism
	organism.Move(toX, toY)
	return true
}

func (w *World) GetEmptyNeighborPositions(x, y int) [][2]int {
	var positions [][2]int
	directions := [][]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

	for _, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		if w.IsEmpty(nx, ny) {
			positions = append(positions, [2]int{nx, ny})
		}
	}
	return positions
}

func (w *World) FindFood(x, y int, diet []string) []Organism {
	var food []Organism
	directions := [][]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

	for _, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		if organism := w.GetOrganism(nx, ny); organism != nil {
			for _, foodType := range diet {
				if organism.GetType() == foodType {
					food = append(food, organism)
					break
				}
			}
		}
	}
	return food
}
func (w *World) GetStatistics() map[string]int {
	stats := map[string]int{"Fox": 0, "Rabbit": 0, "Grass": 0}

	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if organism := w.Grid[y][x]; organism != nil {
				stats[organism.GetType()]++
			}
		}
	}
	return stats
}

func (w *World) GetOrganismsByType(organismType string) []Organism {
	var organisms []Organism
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if organism := w.Grid[y][x]; organism != nil && organism.GetType() == organismType {
				organisms = append(organisms, organism)
			}
		}
	}
	return organisms
}

func (w *World) Simulate() {
	organisms := w.getAllLivingOrganisms()
	rand.Shuffle(len(organisms), func(i, j int) {
		organisms[i], organisms[j] = organisms[j], organisms[i]
	})

	for _, organism := range organisms {
		if organism.GetEnergy() <= 0 {
			continue
		}

		x, y := organism.GetPosition()
		if food := w.FindFood(x, y, organism.GetDiet()); len(food) > 0 {
			switch org := organism.(type) {
			case *Rabbit:
				if org.GetEatingCooldown() == 0 {
					org.Eat()
					fx, fy := food[0].GetPosition()
					w.RemoveOrganism(fx, fy)
				}
			case *Fox:
				if org.GetEatingCooldown() == 0 {
					org.Eat()
					fx, fy := food[0].GetPosition()
					w.RemoveOrganism(fx, fy)
				}
			}
		}
		if organism.GetEnergy() > 0 && organism.CanMove() {
			moved := false
			if organism.CanBreed() && !organism.HasBred() && organism.GetType() != "Grass" {
				moved = w.moveTowardsPartner(organism)
			}
			if !moved {
				if positions := w.GetEmptyNeighborPositions(x, y); len(positions) > 0 {
					newPos := positions[rand.Intn(len(positions))]
					w.MoveOrganism(x, y, newPos[0], newPos[1])
				}
			}
		}
		if organism.GetEnergy() > 0 && organism.CanBreed() && !organism.HasBred() {
			w.tryBreeding(organism)
		}
	}

	w.updateAndCleanup()
	if w.Turn%5 == 0 {
		w.spawnRandomGrass(5)
	}
	w.Turn++
}

func (w *World) tryBreeding(organism Organism) {
	x, y := organism.GetPosition()
	emptyPositions := w.GetEmptyNeighborPositions(x, y)

	if len(emptyPositions) == 0 {
		return
	}
	if organism.GetType() == "Grass" {
		if organism.GetEnergy() >= 4 {
			organism.Breed()
			newPos := emptyPositions[rand.Intn(len(emptyPositions))]
			newGrass := NewGrass(w.nextID, newPos[0], newPos[1])
			w.PlaceOrganism(newGrass)
			w.nextID++
		}
		return
	}
	partner := w.findNearbyPartner(organism, x, y)
	if partner == nil {
		return
	}

	minEnergy := 3
	if organism.GetType() == "Fox" {
		minEnergy = 4
	}

	if organism.GetEnergy() >= minEnergy && partner.GetEnergy() >= minEnergy {
		organism.Breed()
		partner.Breed()

		newPos := emptyPositions[rand.Intn(len(emptyPositions))]
		var newOrganism Organism

		if organism.GetType() == "Rabbit" {
			newOrganism = NewRabbit(w.nextID, newPos[0], newPos[1])
		} else if organism.GetType() == "Fox" {
			newOrganism = NewFox(w.nextID, newPos[0], newPos[1])
		}

		if newOrganism != nil {
			w.PlaceOrganism(newOrganism)
			w.nextID++
		}
	}
}

func (w *World) findNearbyPartner(organism Organism, x, y int) Organism {
	directions := [][]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

	for _, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		if partner := w.GetOrganism(nx, ny); partner != nil {
			if partner.GetType() == organism.GetType() &&
				partner.CanBreed() &&
				!partner.HasBred() &&
				partner.GetEnergy() > 0 {
				return partner
			}
		}
	}
	return nil
}

func (w *World) getAllLivingOrganisms() []Organism {
	var organisms []Organism
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if organism := w.Grid[y][x]; organism != nil && organism.GetEnergy() > 0 {
				organisms = append(organisms, organism)
			}
		}
	}
	return organisms
}
func (w *World) spawnRandomGrass(count int) {
	for i := 0; i < count; i++ {
		for attempts := 0; attempts < 50; attempts++ {
			x, y := rand.Intn(w.Width), rand.Intn(w.Height)
			if w.IsEmpty(x, y) {
				grass := NewGrass(w.nextID, x, y)
				w.PlaceOrganism(grass)
				w.nextID++
				break
			}
		}
	}
}

func (w *World) updateAndCleanup() {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if organism := w.Grid[y][x]; organism != nil {
				if organism.GetEnergy() > 0 {
					organism.NewTurn()
				}
				if organism.GetEnergy() <= 0 {
					w.Grid[y][x] = nil
				}
			}
		}
	}
}

func (w *World) PopulateRandomly(foxCount, rabbitCount, grassCount int) {
	for i := 0; i < foxCount; i++ {
		for attempts := 0; attempts < 100; attempts++ {
			x, y := rand.Intn(w.Width), rand.Intn(w.Height)
			if w.IsEmpty(x, y) {
				fox := NewFox(w.nextID, x, y)
				w.PlaceOrganism(fox)
				w.nextID++
				break
			}
		}
	}
	for i := 0; i < rabbitCount; i++ {
		for attempts := 0; attempts < 100; attempts++ {
			x, y := rand.Intn(w.Width), rand.Intn(w.Height)
			if w.IsEmpty(x, y) {
				rabbit := NewRabbit(w.nextID, x, y)
				w.PlaceOrganism(rabbit)
				w.nextID++
				break
			}
		}
	}
	for i := 0; i < grassCount; i++ {
		for attempts := 0; attempts < 100; attempts++ {
			x, y := rand.Intn(w.Width), rand.Intn(w.Height)
			if w.IsEmpty(x, y) {
				grass := NewGrass(w.nextID, x, y)
				w.PlaceOrganism(grass)
				w.nextID++
				break
			}
		}
	}
}

func (w *World) IsExtinct() bool {
	stats := w.GetStatistics()
	return stats["Fox"] == 0 && stats["Rabbit"] == 0
}

func (w *World) moveTowardsPartner(organism Organism) bool {
	x, y := organism.GetPosition()
	orgType := organism.GetType()
	var closestPartner Organism
	var closestDistance int = 100

	for dy := -3; dy <= 3; dy++ {
		for dx := -3; dx <= 3; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}

			nx, ny := x+dx, y+dy
			if !w.IsValidPosition(nx, ny) {
				continue
			}

			if partner := w.GetOrganism(nx, ny); partner != nil {
				if partner.GetType() == orgType &&
					partner.CanBreed() &&
					!partner.HasBred() &&
					partner.GetEnergy() > 0 {
					distance := dx*dx + dy*dy
					if distance < closestDistance {
						closestDistance = distance
						closestPartner = partner
					}
				}
			}
		}
	}

	if closestPartner == nil {
		return false
	}
	px, py := closestPartner.GetPosition()
	var moveX, moveY int = x, y
	if px > x {
		moveX = x + 1
	} else if px < x {
		moveX = x - 1
	}

	if py > y {
		moveY = y + 1
	} else if py < y {
		moveY = y - 1
	}
	if w.IsEmpty(moveX, moveY) {
		w.MoveOrganism(x, y, moveX, moveY)
		return true
	}

	return false
}
