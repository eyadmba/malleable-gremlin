package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

var queryPingTimeout = 5 * time.Second // Define a timeout for ping within ExecuteQuery

type ExecuteQueryArgument struct {
	ConnectionStringID string `json:"connectionStringId"`
	ConnectionString   string `json:"connectionString"`
	Query              string `json:"query"`
}

type ExecuteQueryResult struct {
	Rows  []map[string]interface{} `json:"rows"`
	Error string                   `json:"error,omitempty"`
	Time  time.Duration            `json:"time"`
}

func (cm *ConnectionManager) ExecuteQuery(ctx context.Context, req *ExecuteQueryArgument) (*ExecuteQueryResult, error) {
	var db *sql.DB
	var err error
	var connStr string

	// Validate inputs
	connProvided := req.ConnectionString != ""
	idProvided := req.ConnectionStringID != ""

	if !connProvided && !idProvided {
		return nil, ErrNeitherInputProvided
	}
	if connProvided && idProvided {
		return nil, ErrBothInputsProvided
	}

	if connProvided {
		connStr = req.ConnectionString
	} else {
		cm.mu.RLock()
		var exists bool
		connStr, exists = cm.connections[req.ConnectionStringID]
		cm.mu.RUnlock()
		if !exists {
			return nil, fmt.Errorf("%w: '%s'", ErrConnIDNotFound, req.ConnectionStringID)
		}
	}

	// Open a new connection
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConnectionSetupFailed, err)
	}
	defer db.Close()

	// Ping the database using derived context with timeout
	ctxPing, cancelPing := context.WithTimeout(ctx, queryPingTimeout)
	defer cancelPing()
	if err := db.PingContext(ctxPing); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConnectionFailed, err)
	}

	return executeQuery(ctx, db, req.Query), nil
}

func executeQuery(ctx context.Context, db *sql.DB, query string) *ExecuteQueryResult {
	start := time.Now()
	result := ExecuteQueryResult{}

	// Execute query using the passed context
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		result.Error = err.Error()
		// Don't return early, capture time and return result below
	} else {
		// Only process rows if the query execution didn't error
		defer rows.Close()

		// Get column names
		columns, err := rows.Columns()
		if err != nil {
			result.Error = err.Error() // Capture column error
		} else {
			// Prepare result set only if columns were retrieved
			result.Rows = make([]map[string]interface{}, 0)
			for rows.Next() {
				// Create a slice of interface{} to hold the values
				values := make([]interface{}, len(columns))
				valuePtrs := make([]interface{}, len(columns))
				for i := range columns {
					valuePtrs[i] = &values[i]
				}

				// Scan the result into the pointers
				if err := rows.Scan(valuePtrs...); err != nil {
					result.Error = err.Error() // Capture scan error
					break                      // Stop processing rows on scan error
				}

				// Create a map for this row
				rowMap := make(map[string]interface{}) // Renamed row -> rowMap
				for i, col := range columns {
					var v interface{}
					val := values[i]
					b, ok := val.([]byte)
					if ok {
						v = string(b)
					} else {
						v = val
					}
					rowMap[col] = v
				}
				result.Rows = append(result.Rows, rowMap)
			}

			// Check for errors encountered during iteration
			if err := rows.Err(); err != nil {
				if result.Error == "" { // Don't overwrite a previous error (e.g., scan error)
					result.Error = err.Error()
				}
			}
		}
	}

	result.Time = time.Since(start)
	return &result
}
