package main

import (
	"fmt"
	"strings"
)

func main() {
	// in Go maps[key]{value} this means you will call the key
	// to get it's value, now you also specify what datatypes you used for it.
	usernames := map[string]string{
		"Richard":   "richard-dragon",
		"Diosel":    "diosel-cute",
		"Joshua":    "joshua-bayot",
		"Realchard": "the-boys",
		"King":      "king-gwapo",
	}

	for {
		fmt.Println("Enter name: ")
		var name string
		_, err := fmt.Scanln(&name)
		if err != nil {
			panic(err)
		}
		name = strings.TrimSpace(name)

		if u, ok := usernames[name]; ok {
			fmt.Printf("%q is the username of %q\n", u, name)
			continue
		}
		// the continue indicates that if it satisfy the condition, it will go back condition
		// the start else it just continue and one more
		// input the code excutetion will be done.

		fmt.Printf("I dont have %v's username, whay is it?\n", name)

		var username string
		_, err = fmt.Scanln(&username)
		if err != nil {
			panic(err)
		}
		username = strings.TrimSpace(username)

		usernames[name] = username

		fmt.Println("data updated.")
	}
}
