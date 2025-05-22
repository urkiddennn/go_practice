package main

import "fmt"

// a value string and list of strings
type Boat struct {
	Name      string
	occupants []string
}

// b as value of Boat and a function AddOccupant that return the
// set value of b form the appended occupants
func (b *Boat) AddOccupant(name string) *Boat {
	b.occupants = append(b.occupants, name)

	return b
}

// Function Manifest() that get the value of Boat and print all the occupants
func (b Boat) Manifest() {
	fmt.Println("The", b.Name, "has the followinf occupants:")
	for _, n := range b.occupants {
		fmt.Println("\t", n)
	}
}

func main() {
	b := &Boat{
		Name: "S.S. DigitalOcean",
	}
	b.AddOccupant("Sammy the Sahrk")
	b.AddOccupant("Larry the Lobster")

	b.Manifest()
}
