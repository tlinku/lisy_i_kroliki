package main

import "fmt"

type Grass struct {
	energy           int
	ID               int
	icon             string
	ate              bool
	canBreed         bool
	bred             bool
	x, y             int
	canMove          bool
	eats             []string
	canEat           bool
	eatingCooldown   int
	breedingCooldown int
}

func NewGrass(id int, x int, y int) *Grass {
	return &Grass{
		ID:               id,
		canEat:           false,
		energy:           6,
		canBreed:         false,
		bred:             false,
		canMove:          false,
		icon:             "🌱",
		x:                x,
		y:                y,
		eats:             []string{},
		breedingCooldown: 2,
	}
}

func (r *Grass) PrintInfo() {
	fmt.Printf("ID: %d\nEnergy: %d\nPosition: (%d,%d)\nAte: %t\nCan Breed: %t\n", r.ID, r.energy, r.x, r.y, r.ate, r.canBreed)
}

func (r *Grass) GetIcon() string {
	return r.icon
}
func (r *Grass) GetType() string {
	return "Grass"
}
func (r *Grass) GetDiet() []string {
	return r.eats
}
func (r *Grass) GetID() int {
	return r.ID
}
func (r *Grass) GetEnergy() int {
	return r.energy
}
func (r *Grass) GetPosition() (int, int) {
	return r.x, r.y
}
func (r *Grass) GetX() int {
	return r.x
}
func (r *Grass) GetY() int {
	return r.y
}
func (r *Grass) HasAte() bool {
	return r.ate
}
func (r *Grass) CanBreed() bool {
	return r.canBreed
}
func (r *Grass) HasBred() bool {
	return r.bred
}
func (r *Grass) CanMove() bool {
	return r.canMove
}
func (r *Grass) GetEatingCooldown() int {
	return r.eatingCooldown
}

func (r *Grass) GetBreedingCooldown() int {
	return r.breedingCooldown
}

func (r *Grass) Breed() {
	if r.canBreed && r.breedingCooldown == 0 {
		r.bred = true
		r.canBreed = false
		r.energy -= 2
		r.breedingCooldown = 4
	}
}

func (r *Grass) Move(x int, y int) {
}
func (r *Grass) Die() {
	r.energy = 0
}

func (r *Grass) NewTurn() {
	if r.breedingCooldown > 0 {
		r.breedingCooldown--
		if r.breedingCooldown == 0 {
			r.canBreed = true
			r.bred = false
		}
	}
	r.energy -= 1
	if r.energy <= 0 {
		r.Die()
	}
}
