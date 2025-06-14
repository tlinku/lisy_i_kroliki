package main

type Organism interface {
	GetIcon() string
	GetID() int
	GetPosition() (int, int)
	GetX() int
	GetY() int
	GetEnergy() int
	Move(x, y int)
	NewTurn()
	CanMove() bool
	GetType() string
	CanBreed() bool
	HasBred() bool
	Breed()
	GetDiet() []string
}
