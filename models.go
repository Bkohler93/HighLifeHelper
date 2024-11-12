package main

type Recipe struct {
	ID          int
	Ingredients []Ingredient
	IsSelected  bool
}

type Ingredient struct {
	Name   string
	Amount int
}
