package main

import "fmt"

type Rabbit struct {
	energy           int
	ID               int
	icon             string
	canEat           bool
	ate              bool
	canBreed         bool
	bred             bool
	x, y             int
	canMove          bool
	eats             []string
	eatingCooldown   int
	breedingCooldown int
}

func NewRabbit(id int, x int, y int) *Rabbit {
	return &Rabbit{
		ID:               id,
		ate:              false,
		energy:           10,
		canEat:           true,
		canBreed:         false,
		bred:             false,
		canMove:          true,
		icon:             "🐰",
		x:                x,
		y:                y,
		eats:             []string{"Grass"},
		eatingCooldown:   0,
		breedingCooldown: 2,
	}
}

func (r *Rabbit) PrintInfo() {
	fmt.Printf("ID: %d\nEnergy: %d\nPosition: (%d,%d)\nAte: %t\nCan Breed: %t\n", r.ID, r.energy, r.x, r.y, r.ate, r.canBreed)
}

func (r *Rabbit) GetIcon() string {
	return r.icon
}
func (r *Rabbit) GetDiet() []string {
	return r.eats
}
func (r *Rabbit) GetID() int {
	return r.ID
}
func (r *Rabbit) GetEnergy() int {
	return r.energy
}
func (r *Rabbit) GetPosition() (int, int) {
	return r.x, r.y
}
func (r *Rabbit) GetX() int {
	return r.x
}
func (r *Rabbit) GetY() int {
	return r.y
}
func (r *Rabbit) HasAte() bool {
	return r.ate
}
func (r *Rabbit) CanBreed() bool {
	return r.canBreed
}
func (r *Rabbit) GetType() string {
	return "Rabbit"
}
func (r *Rabbit) HasBred() bool {
	return r.bred
}
func (r *Rabbit) CanMove() bool {
	return r.canMove
}
func (r *Rabbit) GetEatingCooldown() int {
	return r.eatingCooldown
}

func (r *Rabbit) GetBreedingCooldown() int {
	return r.breedingCooldown
}

func (r *Rabbit) Eat() {
	if r.eatingCooldown == 0 {
		r.ate = true
		r.energy += 6
		r.eatingCooldown = 3
	}
}

func (r *Rabbit) Breed() {
	if r.canBreed && r.breedingCooldown == 0 {
		r.bred = true
		r.canBreed = false
		r.energy -= 1
		r.breedingCooldown = 5
	}
}

func (r *Rabbit) Move(x int, y int) {
	r.x = x
	r.y = y
}
func (r *Rabbit) Die() {
	r.energy = 0
	r.canMove = false
}

func (r *Rabbit) NewTurn() {
	if r.eatingCooldown > 0 {
		r.eatingCooldown--
	}
	if r.breedingCooldown > 0 {
		r.breedingCooldown--
		if r.breedingCooldown == 0 {
			r.canBreed = true
			r.bred = false
		}
	}
	if r.eatingCooldown == 0 {
		r.ate = false
	}
	r.energy -= 1
}
