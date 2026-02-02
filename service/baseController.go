package service

import (
	"database/sql"
	"fmt"
	"strings"
)

type BaseController[Type any] struct {
	DB        *sql.DB
	TableName string

	SelectColumnsDefincition string
	// TODO Add INSERT Column definition
	// TODO Change names. Something like Select and Update Query
	UpdateColumnsDefinition string
	InsertIndices           string

	// List of columns
	// TODO Keep original list as well
	selectDefinitions []ColumnDefinitionInterface
	updateDefinitions []ColumnDefinitionInterface
}

type ColumnDefinitionInterface interface {
	GetColumnName() string
	IsMutable() bool
	ScanValue(dest any) any             // Return pointer to field reference
	GetValue(source any) any            // Return the value for specific fields
	SetValue(dest any, value any) error // Set the value for specific field
}

type ColumnDefinition[Type any, Entity any] struct {
	columnName string
	isMutable  bool
	// TODO Add is Primary key. This field should not end up in the INSERT or UPDATE queries
	fieldAccesor func(*Entity) *Type // Return pointer to field
}

func (c *ColumnDefinition[Type, Entity]) GetColumnName() string {
	return c.columnName
}

func (c *ColumnDefinition[Type, Entity]) IsMutable() bool {
	return c.isMutable
}

func (c *ColumnDefinition[Type, Entity]) ScanValue(dest any) any {
	return c.fieldAccesor(dest.(*Entity))
}

func (c *ColumnDefinition[Type, Entity]) GetValue(source any) any {
	return *c.fieldAccesor(source.(*Entity))
}

func (c *ColumnDefinition[Type, Entity]) SetValue(dest any, val any) error {
	entity := dest.(*Entity)
	typed, ok := val.(Type)
	if !ok {
		return fmt.Errorf("expected type %T, got %T", *new(Type), val)
	}
	*c.fieldAccesor(entity) = typed
	return nil
}

func NewColumnDefinition[Type any, Entity any](
	name string,
	isMutable bool,
	fieldAccesor func(*Entity) *Type,
) *ColumnDefinition[Type, Entity] {
	return &ColumnDefinition[Type, Entity]{
		columnName:   name,
		isMutable:    isMutable,
		fieldAccesor: fieldAccesor,
	}
}

func NewBaseController[Type any](db *sql.DB, tableName string, columnDefinitions []ColumnDefinitionInterface) *BaseController[Type] {

	selectDefinitions := make([]ColumnDefinitionInterface, 0)
	updateDefinitions := make([]ColumnDefinitionInterface, 0)

	var selectColumnsBuilder strings.Builder
	var columnIndicesBuilder strings.Builder
	var UpdateColumnsBuilder strings.Builder

	curIdx := 1
	updateIdx := 1
	for _, columnDefinition := range columnDefinitions {
		isIdField := columnDefinition.GetColumnName() == "id"
		if isIdField {
			continue
		}

		if curIdx > 1 {
			selectColumnsBuilder.WriteString(", ")
			columnIndicesBuilder.WriteString(", ")
		}

		selectColumnsBuilder.WriteString(columnDefinition.GetColumnName())
		fmt.Fprintf(&columnIndicesBuilder, "$%d", curIdx)
		curIdx++

		selectDefinitions = append(selectDefinitions, columnDefinition)

		columnIsNotMutable := !columnDefinition.IsMutable()
		if columnIsNotMutable {
			continue
		}

		if updateIdx > 1 {
			UpdateColumnsBuilder.WriteString(", ")
		}

		fmt.Fprintf(&UpdateColumnsBuilder, "%s = $%d", columnDefinition.GetColumnName(), updateIdx)
		updateIdx++

		updateDefinitions = append(updateDefinitions, columnDefinition)
	}

	return &BaseController[Type]{
		DB:                       db,
		TableName:                tableName,
		SelectColumnsDefincition: selectColumnsBuilder.String(),
		UpdateColumnsDefinition:  UpdateColumnsBuilder.String(),
		InsertIndices:            columnIndicesBuilder.String(),
		selectDefinitions:        selectDefinitions,
		updateDefinitions:        updateDefinitions,
	}
}

// Create executes an insert query within a transaction
func (bc *BaseController[Type]) Create(entity *Type) (sql.Result, error) {
	tx, err := bc.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create the query form the columns
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		bc.TableName,
		bc.SelectColumnsDefincition,
		bc.InsertIndices,
	)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	getters := make([]any, 0, len(bc.selectDefinitions))
	for _, column := range bc.selectDefinitions {
		getters = append(getters, column.GetValue(entity))
	}

	result, err := stmt.Exec(getters...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}

// CreateBatch executes multiple inserts in a single transaction
func (bc *BaseController[Type]) CreateBatch(query string, argsList [][]interface{}) error {
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
// Example: INSERT INTO table (a, b) VALUES ($1, $2), ($3, $4), ($5, $6)
func (bc *BaseController[Type]) CreateBulk(baseQuery string, numColumns int, argsList [][]interface{}) error {
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

		// rowPlaceholder = $1, $2, $3, ...
		var rowPlaceholders []string
		for i := 0; i < numColumns; i++ {
			rowPlaceholders = append(rowPlaceholders, fmt.Sprintf("$%d", paramIndex))
			paramIndex++
		}

		// Placeholder = ($1, $2, $3)
		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ", ")))

		// flatArgs = ($1, $2, $3), ($4, $5, $6) ...
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
func (bc *BaseController[Type]) Update(updateId int64, entity Type) (sql.Result, error) {
	tx, err := bc.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create the query form the columns
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = %d",
		bc.TableName,
		bc.UpdateColumnsDefinition,
		updateId,
	)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	getters := make([]any, 0, len(bc.updateDefinitions))
	for _, column := range bc.updateDefinitions {
		getters = append(getters, column.GetValue(entity))
	}

	result, err := stmt.Exec(getters...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, err
}

// Delete removes a record by ID
func (bc *BaseController[Type]) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", bc.TableName)

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
func (bc *BaseController[Type]) WithTransaction(fn func(*sql.Tx) error) error {
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
func (bc *BaseController[Type]) QueryRow(query string, args ...interface{}) *sql.Row {
	return bc.DB.QueryRow(query, args...)
}

// Query executes a query that returns multiple rows
func (bc *BaseController[Type]) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := bc.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return rows, nil
}
