package main

import (
	"fmt"
	"strings"
)

func main() {
	flavors := []string{"chocolate", "vanilla", "strawberry", "banana"}

	fmt.Println("Write a fruit")
	var fruit string

	_, err := fmt.Scanln(&fruit)
	if err != nil {
		panic(err)
	}

	fruit = strings.TrimSpace(fruit)

	fmt.Println(fruit)

	for _, flav := range flavors {
		switch fruit {
		case fruit:
			fmt.Println(fruit, "is my favorite!")
		case fruit:
			fmt.Println(fruit, "is great!")
		default:
			fmt.Println("I've never tried", flav, "before")

		}
	}
}
