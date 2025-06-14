package main

import "fmt"

type Fox struct {
	energy           int
	ID               int
	icon             string
	ate              bool
	canEat           bool
	canBreed         bool
	bred             bool
	x, y             int
	canMove          bool
	eats             []string
	eatingCooldown   int
	breedingCooldown int
}

func NewFox(id int, x int, y int) *Fox {
	return &Fox{
		ID:               id,
		canEat:           true,
		ate:              false,
		energy:           15,
		canBreed:         false,
		bred:             false,
		canMove:          true,
		icon:             "🦊",
		x:                x,
		y:                y,
		eats:             []string{"Rabbit"},
		eatingCooldown:   2,
		breedingCooldown: 6,
	}
}

func (r *Fox) PrintInfo() {
	fmt.Printf("ID: %d\nEnergy: %d\nPosition: (%d,%d)\nAte: %t\nCan Breed: %t\n", r.ID, r.energy, r.x, r.y, r.ate, r.canBreed)
}

func (r *Fox) GetIcon() string {
	return r.icon
}
func (r *Fox) GetType() string {
	return "Fox"
}
func (r *Fox) GetDiet() []string {
	return r.eats
}
func (r *Fox) GetID() int {
	return r.ID
}
func (r *Fox) GetEnergy() int {
	return r.energy
}
func (r *Fox) GetPosition() (int, int) {
	return r.x, r.y
}
func (r *Fox) GetX() int {
	return r.x
}
func (r *Fox) GetY() int {
	return r.y
}
func (r *Fox) HasAte() bool {
	return r.ate
}
func (r *Fox) CanBreed() bool {
	return r.canBreed
}
func (r *Fox) HasBred() bool {
	return r.bred
}
func (r *Fox) CanMove() bool {
	return r.canMove
}
func (r *Fox) GetEatingCooldown() int {
	return r.eatingCooldown
}

func (r *Fox) GetBreedingCooldown() int {
	return r.breedingCooldown
}

func (r *Fox) Eat() {
	if r.eatingCooldown == 0 {
		r.ate = true
		r.energy += 10
		r.eatingCooldown = 8
	}
}

func (r *Fox) Breed() {
	if r.canBreed && r.breedingCooldown == 0 {
		r.bred = true
		r.canBreed = false
		r.energy -= 2
		r.breedingCooldown = 7
	}
}

func (r *Fox) Move(x int, y int) {
	r.x = x
	r.y = y
}
func (r *Fox) Die() {
	r.energy = 0
	r.canMove = false
}

func (r *Fox) NewTurn() {
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
