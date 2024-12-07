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
		pages   [TABLE_MAX_PAGES]*Page
		numRows uint32
	}
)

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

func (st *Statement) executeInsert(table *Table) ExecuteResult {
	if table.numRows >= TABLE_MAX_ROWS {
		return EXECUTE_TABLE_FULL
	}
	rowToInsert := st.rowToInsert
	rowToInsert.SerializeRow(&table.rowSLot(table.numRows).rows[table.numRows])
	table.numRows++

	return EXECUTE_SUCCESS
}

func (st *Statement) executeStatement(table *Table) ExecuteResult {
	switch st.Type {
	case STATEMENT_INSERT:
		return st.executeInsert(table)
	}
	return EXECUTE_SUCCESS
}

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

func (table *Table) rowSLot(rowNum uint32) *Page {
	pageNum := rowNum / ROWS_PER_PAGE
	if pageNum >= uint32(len(table.pages)) {
		panic("pageNum out of bounds")
	}
	page := table.pages[pageNum]
	if page == nil {
		page = new(Page)
	}
	rowOffset := rowNum % ROWS_PER_PAGE
	byteOffset := rowOffset * ROW_SIZE
	page.currentRow = rowNum
	page.offset = byteOffset
	return page
}
