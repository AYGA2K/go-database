package inputbuffer

type MetaCommandResult int

const (
	META_COMMAND_SUCCESS MetaCommandResult = iota
	META_COMMAND_UNRECOGNIZED_COMMAND
)

type PrepareResult int

const (
	PREPARE_SUCCESS PrepareResult = iota
	PREPARE_UNRECOGNIZED_STATEMENT
)

type StatementType int

const (
	STATEMENT_INSERT StatementType = iota
	STATEMENT_SELECT
)

type ExecuteResult int

const (
	EXECUTE_SUCCESS ExecuteResult = iota
	EXECUTE_TABLE_FULL
	EXECUTE_WRITE_ERROR
)

const (
	ID_OFFSET       = 0
	ID_SIZE         = 4
	USERNAME_OFFSET = ID_OFFSET + ID_SIZE
	USERNAME_SIZE   = 32
	EMAIL_OFFSET    = USERNAME_OFFSET + USERNAME_SIZE
	EMAIL_SIZE      = 255
	TOTAL_SIZE      = EMAIL_OFFSET + EMAIL_SIZE
	ROW_SIZE        = ID_SIZE + USERNAME_SIZE + EMAIL_SIZE
)

const (
	PAGE_SIZE       = 4096
	TABLE_MAX_PAGES = 100
	ROWS_PER_PAGE   = PAGE_SIZE / ROW_SIZE
)
const (
	TABLE_MAX_ROWS = 1000
)
