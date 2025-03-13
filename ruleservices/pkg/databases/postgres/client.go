package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"

	config "github.com/kiranpt03/factorysphere/devicesphere/services/ruleservices/pkg/config"
	log "github.com/kiranpt03/factorysphere/devicesphere/services/ruleservices/pkg/utils/loggers"
	_ "github.com/lib/pq"
)

// PostgreSQLRepository represents a PostgreSQL repository.
type PostgreSQLRepository struct {
	db *sql.DB
}

// NewPostgreSQLRepository creates a new PostgreSQL repository.
func NewPostgreSQLRepository(config *config.Config) (*PostgreSQLRepository, error) {
	// Create a connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.Dbname)
	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Debug("Connected to the database!")

	// Table creation queries
	rulesTableQuery := `
                CREATE TABLE IF NOT EXISTS rules (
                        id VARCHAR(255) PRIMARY KEY,
                        name VARCHAR(255) NOT NULL,
                        severity VARCHAR(255) NOT NULL,
                        status VARCHAR(255) NOT NULL,
                        type VARCHAR(255) NOT NULL,
                        description TEXT NOT NULL,
                        created_at VARCHAR(255) NOT NULL,
                        updated_at VARCHAR(255) NOT NULL
                );
        `
	conditionsTableQuery := `
                CREATE TABLE IF NOT EXISTS conditions (
                        id VARCHAR(255) PRIMARY KEY,
                        rule_id VARCHAR(255) NOT NULL,
                        position VARCHAR(255) NOT NULL,
                        type VARCHAR(255) NOT NULL,
                        device_id VARCHAR(255),
                        device_name VARCHAR(255),
                        property_id VARCHAR(255),
                        property_name VARCHAR(255),
                        operator_id VARCHAR(255),
                        operator_symbol VARCHAR(255),
                        value VARCHAR(255),
                        FOREIGN KEY (rule_id) REFERENCES rules(id)
                );
        `

	// Create tables if they don't exist
	_, err = db.Exec(rulesTableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create rules table: %w", err)
	}

	_, err = db.Exec(conditionsTableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create conditions table: %w", err)
	}

	return &PostgreSQLRepository{db: db}, nil
}

// GetAll retrieves all records from the specified table.
func (r *PostgreSQLRepository) GetAll(tableName string) ([]map[string]interface{}, error) {
	log.Info("Get all devices from database")
	var records []map[string]interface{}

	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		record := make(map[string]interface{})
		for i, col := range cols {
			record[col] = columns[i]
		}

		records = append(records, record)
	}

	return records, nil
}

// Get retrieves a record by ID from the specified table.
func (r *PostgreSQLRepository) Get(tableName string, id string) (map[string]interface{}, error) {
	var record map[string]interface{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tableName)
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		record = make(map[string]interface{})
		for i, col := range cols {
			record[col] = columns[i]
		}
	}

	return record, nil
}

// Create creates a new record in the specified table.
func (r *PostgreSQLRepository) Create(tableName string, data map[string]interface{}) (string, error) {
	// Generate a random ID
	id := uuid.New().String()

	// Add the ID to the data map
	data["id"] = id

	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	params := make([]string, 0, len(data))

	for k, v := range data {
		columns = append(columns, k)
		values = append(values, v)
		params = append(params, fmt.Sprintf("$%d", len(params)+1))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(params, ", "))

	var createdID string
	err := r.db.QueryRow(query, values...).Scan(&createdID)
	fmt.Println(err)

	return createdID, err
}

// Update updates an existing record in the specified table.
func (r *PostgreSQLRepository) Update(tableName string, id string, data map[string]interface{}) (int64, error) {
	var affectedRows int64

	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	for k, v := range data {
		columns = append(columns, fmt.Sprintf("%s = $%d", k, len(values)+1))
		values = append(values, v)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d RETURNING id",
		tableName,
		strings.Join(columns, ", "),
		len(values)+1)

	values = append(values, id)

	err := r.db.QueryRow(query, values...).Scan(&affectedRows)

	return affectedRows, err
}

// Delete deletes a record by the specified column and value from the specified table.
func (r *PostgreSQLRepository) Delete(tableName string, column string, value interface{}) error {

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", tableName, column)
	_, err := r.db.Exec(query, value)

	return err
}

// CheckExist checks if a record with the specified column and value exists in the table.
func (r *PostgreSQLRepository) CheckExist(tableName string, column string, value interface{}) (bool, error) {
	var exists bool

	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE %s = $1)", tableName, column)
	err := r.db.QueryRow(query, value).Scan(&exists)

	return exists, err
}

// ExecuteQuery executes a raw SQL query and returns the results.
func (r *PostgreSQLRepository) ExecuteQuery(query string, args ...interface{}) ([]map[string]interface{}, error) {
	log.Debug("Executing query: %s", query)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		record := make(map[string]interface{})
		for i, col := range cols {
			record[col] = columns[i]
		}

		results = append(results, record)
	}

	return results, nil
}
