package serror

import (
	"errors"
	"fmt"

	"golang.org/x/xerrors"
)

var (
	// ErrNoMoreRow indicates there is no rows left to read
	ErrNoMoreRow = errors.New("no more rows available")
	// ErrInsertCommand indicates error with insert query operation
	ErrInsertCommand = errors.New("insert command operation")
	// ErrDeleteCommand indicates error with insert query operation
	ErrDeleteCommand = errors.New("delete command operation")
	// ErrUpdateCommand indicates error with insert query operation
	ErrUpdateCommand = errors.New("update command operation")
)
var (
	// ErrUserNotFound indicates no user associated with either ID or emailID
	ErrUserNotFound = errors.New("no user found")
)

var (
	// ErrTodoNotFound indicates no todo associated with either todoID or userID
	ErrTodoNotFound = errors.New("no todo found")
)

// QueryError reports the error and QueryType in compact form
// that are returned when any db triggers an error
// QueryError should be returned as a part of API.
// Any underlying error which is part of any external package
// should not be wrapped
type QueryError struct {
	// frame
	frame xerrors.Frame

	// Query
	Query string
	Err   error
	// UnderlyingErrorString stores the original error if any
	UnderlyingErrorString string
}

// FormatError define how the error will be represented
// when we print it to the console or
// otherwise want to retrieve the simple / default value
func (qe QueryError) FormatError(p xerrors.Printer) error {
	p.Printf("[%v] - underlying error [%s]", qe.Err, qe.UnderlyingErrorString)
	qe.frame.Format(p)
	return nil
}

// Format provide backwards compatibility
func (qe QueryError) Format(f fmt.State, c rune) {
	xerrors.FormatError(qe, f, c)
}

// Error Implements Error interface
func (qe QueryError) Error() string {
	return fmt.Sprintf("[%v] - underlying error [%s]", qe.Err, qe.UnderlyingErrorString)
}

// NewQueryError returns an initialized error of Error type
func NewQueryError(queryName string, err error, originalErrS string) error {
	pe := QueryError{
		Err:                   err,
		Query:                 queryName,
		UnderlyingErrorString: originalErrS,
		frame:                 xerrors.Caller(1), // skip the first frame
	}
	return xerrors.Errorf("error executing query %q: %w", queryName, pe)
}
