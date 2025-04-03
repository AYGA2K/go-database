package main

import (
	"fmt"

	inputbuffer "github.com/AYGA2K/go-database/internal/input_buffer"
)

func main() {
	fmt.Println("Welcome to the DB you have always dreamed about ")
	// inputBuffer := inputBuffer.NewInputBuffer()
	// for {
	// 	inputBuffer.ReadInput()
	// }

	table, err := inputbuffer.NewTable("database.db")
	if err != nil {
		fmt.Println("Error initializing table:", err)
		return
	}
	defer table.Close()

	// Insert a row
	// statement := inputbuffer.NewStatement{
	// 	rowToInsert: inputbuffer.Row{
	// 		id:       1,
	// 		userName: []byte("john_doe"),
	// 		email:    []byte("john@example.com"),
	// 	},
	// 	Type: inputbuffer.STATEMENT_INSERT,
	// }
	statement := inputbuffer.NewStatement(inputbuffer.NewRow([]byte("john doe"), []byte("john@example.com"), 1), inputbuffer.STATEMENT_INSERT)
	result := statement.ExecuteInsert(table)
	if result != inputbuffer.EXECUTE_SUCCESS {
		fmt.Println("Insert failed")
		return
	}

	// Select rows
	statement.Type = inputbuffer.STATEMENT_SELECT
	statement.ExecuteSelect(table)
}
