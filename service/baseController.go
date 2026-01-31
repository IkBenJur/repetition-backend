package service

import (
	"database/sql"
	"fmt"
	"strings"
)

type BaseController struct {
	DB        *sql.DB
	TableName string

	// Map containing columns and if they are muttable
	columnAndMutability map[string]bool
}

func NewBaseController(db *sql.DB, tableName string, columnAndMutability map[string]bool) *BaseController {
	return &BaseController{
		DB:                  db,
		TableName:           tableName,
		columnAndMutability: columnAndMutability,
	}
}

// Create executes an insert query within a transaction
func (bc *BaseController) Create(args ...interface{}) (sql.Result, error) {
	tx, err := bc.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create the query form the columns
	var columnNamesBuilder strings.Builder
	var columnIndicesBuilder strings.Builder
	curIdx := 1
	for columnName, _ := range bc.columnAndMutability {
		if curIdx > 1 {
			columnNamesBuilder.WriteString(", ")
			columnIndicesBuilder.WriteString(", ")
		}

		columnNamesBuilder.WriteString(columnName)
		fmt.Fprintf(&columnIndicesBuilder, "?%d", curIdx)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		bc.TableName,
		columnNamesBuilder.String(),
		columnIndicesBuilder.String(),
	)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}

// CreateBatch executes multiple inserts in a single transaction
func (bc *BaseController) CreateBatch(query string, argsList [][]interface{}) error {
	if len(argsList) == 0 {
		return nil
	}

	tx, err := bc.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, args := range argsList {
		if _, err := stmt.Exec(args...); err != nil {
			return fmt.Errorf("failed to execute batch insert: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CreateBulk creates a single INSERT with multiple value sets (more efficient for PostgreSQL)
// Example: INSERT INTO table (a, b) VALUES (?1, ?2), (?3, ?4), (?5, ?6)
func (bc *BaseController) CreateBulk(baseQuery string, numColumns int, argsList [][]interface{}) error {
	if len(argsList) == 0 {
		return nil
	}

	tx, err := bc.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Build the bulk insert query
	var placeholders []string
	var flatArgs []interface{}
	paramIndex := 1

	for _, args := range argsList {

		// rowPlaceholder = ?1, ?2, ?3, ...
		var rowPlaceholders []string
		for i := 0; i < numColumns; i++ {
			rowPlaceholders = append(rowPlaceholders, fmt.Sprintf("?%d", paramIndex))
			paramIndex++
		}

		// Placeholder = (?1, ?2, ?3)
		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ", ")))

		// flatArgs = (?1, ?2, ?3), (?4, ?5, ?6) ...
		flatArgs = append(flatArgs, args...)
	}

	query := fmt.Sprintf("%s VALUES %s", baseQuery, strings.Join(placeholders, ", "))

	_, err = tx.Exec(query, flatArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute bulk insert: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Update executes an update query within a transaction
func (bc *BaseController) Update(updateId int, args ...interface{}) (sql.Result, error) {
	tx, err := bc.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create the query form the columns
	var updateQuery strings.Builder
	curIdx := 1
	for columnName, _ := range bc.columnAndMutability {
		columnIsNotMutable := !bc.columnAndMutability[columnName]
		if columnIsNotMutable {
			continue
		}

		if curIdx > 1 {
			updateQuery.WriteString(", ")
		}

		fmt.Fprintf(&updateQuery, "%s = ?%d", columnName, curIdx)
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = %d",
		bc.TableName,
		updateQuery.String(),
		updateId,
	)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, err
}

// Delete removes a record by ID
func (bc *BaseController) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?1", bc.TableName)

	tx, err := bc.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no record found with id %d", id)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithTransaction allows you to execute multiple operations in a single transaction
func (bc *BaseController) WithTransaction(fn func(*sql.Tx) error) error {
	tx, err := bc.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// QueryRow executes a query that returns a single row
func (bc *BaseController) QueryRow(query string, args ...interface{}) *sql.Row {
	return bc.DB.QueryRow(query, args...)
}

// Query executes a query that returns multiple rows
func (bc *BaseController) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := bc.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return rows, nil
}
