package main

import "fmt"

// GoodFunction is a well-documented function.
func GoodFunction() {
	fmt.Println("This is a very long and specific magic string.")
}

func BadFunction() {
	fmt.Println("This function needs a comment!")
	fmt.Println("This function needs a comment!")
	fmt.Println("This function needs a comment!")
	fmt.Println("This function needs a comment!")
}
