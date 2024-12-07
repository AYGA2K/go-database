package main

import (
	"fmt"

	inputBuffer "github.com/AYGA2K/go-database/internal/input_buffer"
)

func main() {
	fmt.Println("Welcome to the DB you have always dreamed about ")
	inputBuffer := inputBuffer.NewInputBuffer()
	for {
		inputBuffer.ReadInput()
	}
}
