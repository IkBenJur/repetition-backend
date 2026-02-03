package service

import (
	"database/sql"
	"fmt"
	"strings"
)

type BaseController[Type any] struct {
	DB        *sql.DB
	TableName string

	SelectColumns string
	InsertColumns string
	UpdateColumns string
	InsertIndices string

	// List of columns
	PrimaryKeyDefinition ColumnDefinitionInterface
	columnDefinitions    []ColumnDefinitionInterface
	insertDefinitions    []ColumnDefinitionInterface
	updateDefinitions    []ColumnDefinitionInterface
}

type ColumnDefinitionInterface interface {
	GetColumnName() string
	IsMutable() bool
	IsPrimaryKey() bool
	ScanValue(dest any) any             // Return pointer to field reference
	GetValue(source any) any            // Return the value for specific fields
	SetValue(dest any, value any) error // Set the value for specific field
}

type ColumnDefinition[Type any, Entity any] struct {
	columnName   string
	isMutable    bool
	isPrimaryKey bool
	fieldAccesor func(*Entity) *Type // Return pointer to field
}

func (c *ColumnDefinition[Type, Entity]) GetColumnName() string {
	return c.columnName
}

func (c *ColumnDefinition[Type, Entity]) IsMutable() bool {
	return c.isMutable
}

// For writing values to entity. ScanRows. Returns pointer
func (c *ColumnDefinition[Type, Entity]) ScanValue(dest any) any {
	return c.fieldAccesor(dest.(*Entity))
}

// For reading from entity. Handles nil pointers
func (c *ColumnDefinition[Type, Entity]) GetValue(source any) any {
	ptr := c.fieldAccesor(source.(*Entity))
	if ptr == nil {
		return nil
	}
	return *ptr
}

// TODO Prob remove. Not used
func (c *ColumnDefinition[Type, Entity]) SetValue(dest any, val any) error {
	entity := dest.(*Entity)
	typed, ok := val.(Type)
	if !ok {
		return fmt.Errorf("expected type %T, got %T", *new(Type), val)
	}
	*c.fieldAccesor(entity) = typed
	return nil
}

func (c *ColumnDefinition[Type, Entity]) IsPrimaryKey() bool {
	return c.isPrimaryKey
}

func NewPrimaryKeyColumnDefinition[Type any, Entity any](
	name string,
	isMutable bool,
	fieldAccesor func(*Entity) *Type,
) *ColumnDefinition[Type, Entity] {
	return &ColumnDefinition[Type, Entity]{
		columnName:   name,
		isPrimaryKey: true,
		isMutable:    isMutable,
		fieldAccesor: fieldAccesor,
	}
}

func NewColumnDefinition[Type any, Entity any](
	name string,
	isMutable bool,
	fieldAccesor func(*Entity) *Type,
) *ColumnDefinition[Type, Entity] {
	return &ColumnDefinition[Type, Entity]{
		columnName:   name,
		isPrimaryKey: false,
		isMutable:    isMutable,
		fieldAccesor: fieldAccesor,
	}
}

func NewBaseController[Type any](db *sql.DB, tableName string, columnDefinitions []ColumnDefinitionInterface) *BaseController[Type] {

	insertDefinitions := make([]ColumnDefinitionInterface, 0)
	updateDefinitions := make([]ColumnDefinitionInterface, 0)

	var selectColumnsBuilder strings.Builder
	var insertColumnsBuilder strings.Builder
	var insertIndicesBuilder strings.Builder
	var updateColumnsBuilder strings.Builder

	var primaryKeyDefinition ColumnDefinitionInterface

	selectIdx := 1
	insertIdx := 1
	updateIdx := 1
	for _, columnDefinition := range columnDefinitions {

		if selectIdx > 1 {
			selectColumnsBuilder.WriteString(", ")
		}

		selectColumnsBuilder.WriteString(columnDefinition.GetColumnName())
		selectIdx++

		if columnDefinition.IsPrimaryKey() {
			primaryKeyDefinition = columnDefinition
			continue
		}

		if insertIdx > 1 {
			insertColumnsBuilder.WriteString(", ")
			insertIndicesBuilder.WriteString(", ")
		}

		insertColumnsBuilder.WriteString(columnDefinition.GetColumnName())
		fmt.Fprintf(&insertIndicesBuilder, "$%d", insertIdx)
		insertIdx++

		insertDefinitions = append(insertDefinitions, columnDefinition)

		columnIsNotMutable := !columnDefinition.IsMutable()
		if columnIsNotMutable {
			continue
		}

		if updateIdx > 1 {
			updateColumnsBuilder.WriteString(", ")
		}

		fmt.Fprintf(&updateColumnsBuilder, "%s = $%d", columnDefinition.GetColumnName(), updateIdx)
		updateIdx++

		updateDefinitions = append(updateDefinitions, columnDefinition)
	}

	return &BaseController[Type]{
		DB:                   db,
		TableName:            tableName,
		SelectColumns:        selectColumnsBuilder.String(),
		InsertColumns:        insertColumnsBuilder.String(),
		UpdateColumns:        updateColumnsBuilder.String(),
		InsertIndices:        insertIndicesBuilder.String(),
		PrimaryKeyDefinition: primaryKeyDefinition,
		columnDefinitions:    columnDefinitions,
		insertDefinitions:    insertDefinitions,
		updateDefinitions:    updateDefinitions,
	}
}

// Create executes an insert query within a transaction
func (bc *BaseController[Type]) Create(entity *Type) (int, error) {
	tx, err := bc.DB.Begin()
	if err != nil {
		return -1, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create the query form the columns
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		bc.TableName,
		bc.InsertColumns,
		bc.InsertIndices,
	)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return -1, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	getters := make([]any, 0, len(bc.insertDefinitions))
	for _, column := range bc.insertDefinitions {
		getters = append(getters, column.GetValue(entity))
	}

	var newId int
	err = stmt.QueryRow(getters...).Scan(&newId)
	if err != nil {
		return -1, fmt.Errorf("failed to execute statement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return -1, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return newId, nil
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
func (bc *BaseController[Type]) CreateBulk(entities []*Type) ([]int64, error) {
	ids := make([]int64, len(entities))
	tx, err := bc.DB.Begin()
	if err != nil {
		return ids, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Build the bulk insert query
	numColumns := len(bc.insertDefinitions)
	numRows := len(entities)
	getters := make([]any, 0, numColumns*numRows)
	bulkInsertValues := make([]string, 0, numRows)

	paramIndex := 1
	for _, entity := range entities {
		rowPlaceholders := make([]string, 0, numColumns)

		for _, column := range bc.insertDefinitions {
			rowPlaceholders = append(rowPlaceholders, fmt.Sprintf("$%d", paramIndex))
			getters = append(getters, column.GetValue(entity))
			paramIndex++
		}

		bulkInsertValues = append(bulkInsertValues, fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ", ")))
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s RETURNING id",
		bc.TableName,
		bc.InsertColumns,
		strings.Join(bulkInsertValues, ", "),
	)

	rows, err := tx.Query(query, getters...)
	if err != nil {
		return ids, fmt.Errorf("failed to execute bulk insert: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return ids, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// TODO Rows is empty with batch insert for standard sql driver.
	// Should be fixed when switching the the pgx driver as well
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids, nil
}

// Update executes an update query within a transaction
func (bc *BaseController[Type]) Update(entity *Type) (int64, error) {
	tx, err := bc.DB.Begin()
	if err != nil {
		return -1, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create the query form the columns TODO SQL INJECTION
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = $%d RETURNING id",
		bc.TableName,
		bc.UpdateColumns,
		len(bc.updateDefinitions)+1,
	)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return -1, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	updateFields := make([]any, 0, len(bc.updateDefinitions)+1) // Plus 1 because of ID field
	for _, column := range bc.updateDefinitions {
		updateFields = append(updateFields, column.GetValue(entity))
	}

	updateFields = append(updateFields, bc.PrimaryKeyDefinition.GetValue(entity))

	var updatedId int64
	err = stmt.QueryRow(updateFields...).Scan(&updatedId)
	if err != nil {
		return -1, fmt.Errorf("failed to execute statement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return -1, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return updatedId, err
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
