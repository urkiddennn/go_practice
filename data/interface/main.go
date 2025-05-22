package main

import (
	"fmt"
	"strings"
)

// This means that the Ocean is have a value of Array of Creatures
type Ocean struct {
	Creatures []string
}

// this functions deals with ocean, getting the ocean value as o
// returning a join array as strings
func (o Ocean) String() string {
	return strings.Join(o.Creatures, ", ")
}

func log(header string, s fmt.Stringer) {
	fmt.Println(header, ":", s)
}

func main() {
	o := Ocean{
		Creatures: []string{
			"sea urchin",
			"lobster",
			"shark",
		},
	}
	log("ocean contains", o)
}
