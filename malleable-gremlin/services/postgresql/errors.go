package postgresql

import "errors"

var (
	// ErrBothInputsProvided indicates that both a connection string and connection ID were given when only one is allowed.
	ErrBothInputsProvided = errors.New("only one of connection string or connection ID should be provided")

	// ErrNeitherInputProvided indicates that neither a connection string nor connection ID was provided.
	ErrNeitherInputProvided = errors.New("either connection string or connection ID must be provided")

	// ErrConnIDNotFound indicates that the provided connection ID does not exist in the manager.
	ErrConnIDNotFound = errors.New("connection ID not found")

	// ErrConnectionSetupFailed indicates a failure during the initial setup of the database connection pool (e.g., sql.Open failed due to bad DSN or driver).
	ErrConnectionSetupFailed = errors.New("failed to configure database connection")

	// ErrConnectionFailed indicates a failure during the establishing or validating of a database connection (e.g., sql.Open, db.Ping).
	ErrConnectionFailed = errors.New("failed to connect to database")
)
