package inputbuffer

import (
	"bytes"
	"fmt"
	"os"
	"unsafe"
)

type Row struct {
	userName []byte
	email    []byte
	id       uint32
}

type Statement struct {
	rowToInsert Row
	Type        StatementType
}
type (
	Page struct {
		rows       [][]byte
		currentRow uint32
		offset     uint32
	}
	Table struct {
		file    *os.File
		numRows uint32
	}
)

func NewRow(userName []byte, email []byte, id uint32) Row {
	row := Row{
		id: id,
	}
	// copy(row.userName[:], userName)
	// copy(row.email[:], email)
	row.userName = append([]byte(nil), userName...)
	row.email = append([]byte(nil), email...)
	return row
}

func NewStatement(row Row, statementType StatementType) Statement {
	return Statement{
		rowToInsert: row,
		Type:        statementType,
	}

}

func (table *Table) Close() error {
	return table.file.Close()
}
func NewTable(filename string) (*Table, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	// Determine the number of rows in the file
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	numRows := fileInfo.Size() / ROW_SIZE

	return &Table{
		file:    file,
		numRows: uint32(numRows),
	}, nil
}

func (st *Statement) prepareStatement(Ib *InputBuffer) PrepareResult {
	if bytes.Equal(Ib.Char[:6], []byte("insert")) {
		st.Type = STATEMENT_INSERT
		var tmpId int
		argsAssigned, err := fmt.Sscanf(
			string(Ib.Char),
			"insert %d %s %s",
			&tmpId,
			&st.rowToInsert.userName,
			&st.rowToInsert.email,
		)
		if err != nil || argsAssigned < 3 {
			fmt.Println("Error: Prepare syntax error: ", err.Error())
			os.Exit(1)
		}
		st.rowToInsert.id = uint32(tmpId)
		return PREPARE_SUCCESS
	}

	if bytes.Equal(Ib.Char[:6], []byte("select")) {
		st.Type = STATEMENT_SELECT
		return PREPARE_SUCCESS
	}

	return PREPARE_UNRECOGNIZED_STATEMENT
}

func (st *Statement) ExecuteInsert(table *Table) ExecuteResult {
	if table.numRows >= TABLE_MAX_ROWS {
		return EXECUTE_TABLE_FULL
	}
	rowToInsert := st.rowToInsert
	offset := int64(table.numRows) * ROW_SIZE
	buffer := make([]byte, ROW_SIZE)

	// page := table.rowSLot(table.numRows)
	// if page.rows[page.currentRow] == nil {
	// 	page.rows[page.currentRow] = make([]byte, ROW_SIZE)
	// }
	rowToInsert.SerializeRow(&buffer)

	// Write the serialized row to the file
	_, err := table.file.WriteAt(buffer, offset)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return EXECUTE_WRITE_ERROR
	}

	table.numRows++

	return EXECUTE_SUCCESS
}
func (st *Statement) ExecuteSelect(table *Table) ExecuteResult {
	for i := uint32(0); i < table.numRows; i++ {
		offset := int64(i) * ROW_SIZE
		buffer := make([]byte, ROW_SIZE)
		// page := table.rowSLot(i)
		// if page.rows[page.currentRow] == nil {
		// 	continue // Skip uninitialized rows
		// }
		// row := Row{}
		// row.DeserializeRow(&page.rows[page.currentRow])
		// fmt.Printf("(%d , %s ,%s) \n", row.id, row.userName, row.email)
		_, err := table.file.ReadAt(buffer[:], offset)
		if err != nil {
			fmt.Println("Error reading from file:", err)
			continue
		}

		// Deserialize the row
		row := Row{}
		row.DeserializeRow(&buffer)

		// Print the row
		fmt.Printf("(%d, %s, %s)\n", row.id, row.userName, row.email)
	}
	return EXECUTE_SUCCESS
}

func (st *Statement) ExecuteStatement(table *Table) ExecuteResult {
	switch st.Type {
	case STATEMENT_INSERT:
		return st.ExecuteInsert(table)
	case STATEMENT_SELECT:
		return st.ExecuteSelect(table)
	}
	return EXECUTE_SUCCESS
}

// func (source *Row) SerializeRow(destination *[]byte) {
// 	// Serialize id
// 	copy((*destination)[ID_OFFSET:ID_OFFSET+ID_SIZE], (*[ID_SIZE]byte)(unsafe.Pointer(&source.id))[:])
// 	fmt.Printf("Serialized ID: %v\n", (*destination)[ID_OFFSET:ID_OFFSET+ID_SIZE])

// 	// Serialize userName
// 	copy((*destination)[USERNAME_OFFSET:USERNAME_OFFSET+USERNAME_SIZE], source.userName[:])
// 	for i := len(source.userName); i < USERNAME_SIZE; i++ {
// 		(*destination)[USERNAME_OFFSET+i] = 0 // Pad with zeros
// 	}
// 	fmt.Printf("Serialized UserName: %v\n", (*destination)[USERNAME_OFFSET:USERNAME_OFFSET+USERNAME_SIZE])

// 	// Serialize email
// 	copy((*destination)[EMAIL_OFFSET:EMAIL_OFFSET+EMAIL_SIZE], source.email[:])
// 	for i := len(source.email); i < EMAIL_SIZE; i++ {
// 		(*destination)[EMAIL_OFFSET+i] = 0 // Pad with zeros
// 	}
// 	fmt.Printf("Serialized Email: %v\n", (*destination)[EMAIL_OFFSET:EMAIL_OFFSET+EMAIL_SIZE])
// }
// func (destination *Row) DeserializeRow(source *[]byte) {
// 	// Deserialize id
// 	copy((*[ID_SIZE]byte)(unsafe.Pointer(&destination.id))[:], (*source)[ID_OFFSET:ID_OFFSET+ID_SIZE])

// 	// Deserialize userName
// 	copy(destination.userName[:], (*source)[USERNAME_OFFSET:USERNAME_OFFSET+USERNAME_SIZE])

// 	// Deserialize email
// 	copy(destination.email[:], (*source)[EMAIL_OFFSET:EMAIL_OFFSET+EMAIL_SIZE])
// }

// SerializeRow serializes the Row struct into a byte slice.
func (source *Row) SerializeRow(destination *[]byte) {
	copy((*destination)[ID_OFFSET:ID_OFFSET+ID_SIZE], (*[ID_SIZE]byte)(unsafe.Pointer(&source.id))[:])

	copy((*destination)[USERNAME_OFFSET:USERNAME_OFFSET+USERNAME_SIZE], source.userName[:])

	copy((*destination)[EMAIL_OFFSET:EMAIL_OFFSET+EMAIL_SIZE], source.email[:])
}

// DeserializeRow deserializes a byte slice into a Row struct.
func (destination *Row) DeserializeRow(source *[]byte) {
	copy((*[ID_SIZE]byte)(unsafe.Pointer(&destination.id))[:], (*source)[ID_OFFSET:ID_OFFSET+ID_SIZE])

	copy(destination.userName[:], (*source)[USERNAME_OFFSET:USERNAME_OFFSET+USERNAME_SIZE])

	copy(destination.email[:], (*source)[EMAIL_OFFSET:EMAIL_OFFSET+EMAIL_SIZE])
}

// func (table *Table) rowSLot(rowNum uint32) *Page {
// 	pageNum := rowNum / ROWS_PER_PAGE
// 	if pageNum >= TABLE_MAX_PAGES {
// 		panic("pageNum out of bounds")
// 	}

// 	// Initialize the page if it doesn't exist
// 	if table.pages[pageNum] == nil {
// 		table.pages[pageNum] = &Page{
// 			rows:       make([][]byte, ROWS_PER_PAGE),
// 			currentRow: 0,
// 			offset:     0,
// 		}
// 	}

// 	rowOffset := rowNum % ROWS_PER_PAGE
// 	page := table.pages[pageNum]
// 	page.currentRow = rowNum
// 	page.offset = rowOffset
// 	return page
// }
