package inputbuffer

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

type InputBuffer struct {
	Char        []byte
	InputLength int
}

func NewInputBuffer() *InputBuffer {
	return &InputBuffer{
		Char:        []byte{},
		InputLength: 0,
	}
}

func (Ib *InputBuffer) ReadInput() {
	scanner := bufio.NewScanner(os.Stdin)
	printPrompt()

	if scanner.Scan() {
		input := scanner.Bytes()
		Ib.Char = input
		Ib.InputLength = len(input)
		if Ib.Char[0] == '.' {
			switch Ib.metaCommand() {
			case META_COMMAND_SUCCESS:
			case META_COMMAND_UNRECOGNIZED_COMMAND:
				fmt.Printf("Unreconized command %s \n", string(Ib.Char))
			}
		}
		var statement Statement
		switch statement.prepareStatement(Ib) {
		case PREPARE_SUCCESS:
			fmt.Println(statement.rowToInsert.id)
			fmt.Println(string(statement.rowToInsert.email))
			fmt.Println(string(statement.rowToInsert.userName))

		case PREPARE_UNRECOGNIZED_STATEMENT:
			fmt.Printf("Unrecognized keyword at start of %s.\n", Ib.Char)
		}
		statement.executeStatement()
		fmt.Println("Executed")
	} else if err := scanner.Err(); err != nil {

		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)
	}
}

func (Ib *InputBuffer) metaCommand() MetaCommandResult {
	if bytes.Equal(Ib.Char, []byte(".exit")) {
		fmt.Println("Exiting...")
		os.Exit(0)
	}
	return META_COMMAND_UNRECOGNIZED_COMMAND
}

func printPrompt() {
	fmt.Print("db > ")
}
