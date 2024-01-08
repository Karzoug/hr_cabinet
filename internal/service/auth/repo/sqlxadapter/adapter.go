package sqlxadapter

import (
	"database/sql"
	"fmt"
	"runtime"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/jmoiron/sqlx"
)

// CasbinRule ...
type CasbinRule struct {
	ID    int64          `db:"id"`
	PType string         `db:"ptype"`
	V0    sql.NullString `db:"v0"`
	V1    sql.NullString `db:"v1"`
	V2    sql.NullString `db:"v2"`
	V3    sql.NullString `db:"v3"`
	V4    sql.NullString `db:"v4"`
	V5    sql.NullString `db:"v5"`
}

// Adapter represents the sqlx adapter for policy storage.
type Adapter struct {
	db         *sqlx.DB
	tableName  string
	isFiltered bool
}

// Filter ...
type Filter struct {
	PType []string
	V0    []string
	V1    []string
	V2    []string
	V3    []string
	V4    []string
	V5    []string
}

// AdapterOptions contains all possible configuration options.
type AdapterOptions struct {
	DriverName     string
	DataSourceName string
	TableName      string
	DB             *sqlx.DB
}

func finalizer(a *Adapter) {
	a.db.Close()
}

func loadPolicyLine(line CasbinRule, model model.Model) {
	lineText := line.PType
	if line.V0.Valid {
		lineText += ", " + line.V0.String
	}
	if line.V1.Valid {
		lineText += ", " + line.V1.String
	}
	if line.V2.Valid {
		lineText += ", " + line.V2.String
	}
	if line.V3.Valid {
		lineText += ", " + line.V3.String
	}
	if line.V4.Valid {
		lineText += ", " + line.V4.String
	}
	if line.V5.Valid {
		lineText += ", " + line.V5.String
	}
	persist.LoadPolicyLine(lineText, model)
}

func savePolicyLine(ptype string, rule []string) CasbinRule {
	line := CasbinRule{}
	line.PType = ptype
	if len(rule) > 0 {
		line.V0.String = rule[0]
	}
	if len(rule) > 1 {
		line.V1.String = rule[1]
	}
	if len(rule) > 2 {
		line.V2.String = rule[2]
	}
	if len(rule) > 3 {
		line.V3.String = rule[3]
	}
	if len(rule) > 4 {
		line.V4.String = rule[4]
	}
	if len(rule) > 5 {
		line.V5.String = rule[5]
	}
	return line
}

func (a *Adapter) dropTable() {
	_, err := a.db.Exec(fmt.Sprintf("DELETE FROM %s", a.tableName))
	if err != nil {
		panic(err)
	}
}

func (a *Adapter) ensureTable() {
	_, err := a.db.Exec(fmt.Sprintf("SELECT 1 FROM %s LIMIT 1", a.tableName))
	if err != nil {
		panic(err)
	}
}

func (a *Adapter) insertPolicyLine(line *CasbinRule) (err error) {
	query := fmt.Sprintf(
		"INSERT INTO %s (ptype, v0, v1, v2, v3, v4, v5) VALUES (:ptype, :v0, :v1, :v2, :v3, :v4, :v5)",
		a.tableName,
	)
	_, err = a.db.NamedExec(query, line)
	if err != nil {
		return
	}
	return
}

func (a *Adapter) deletePolicyLine(line *CasbinRule) (err error) {
	query := fmt.Sprintf(
		"DELETE FROM %s WHERE ptype = :ptype AND v0 = :v0 AND v1 = :v1 AND v2 = :v2 AND v3 = :v3 AND v4 = :v4 AND v5 = :v5",
		a.tableName,
	)
	_, err = a.db.NamedExec(query, line)
	if err != nil {
		return
	}
	return
}

// NewAdapter is the constructor for Adapter
// Deprecated: Use NewAdapterFromOptions instead
func NewAdapter(driverName string, dataSourceName string) *Adapter {
	db, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	a := &Adapter{
		db:        db,
		tableName: "casbin_rule",
	}
	a.ensureTable()
	// Call the destructor when the object is released.
	runtime.SetFinalizer(a, finalizer)
	return a
}

// NewAdapterByDB is the constructor for Adapter with existed connection
// Deprecated: Use NewAdapterFromOptions instead
func NewAdapterByDB(db *sqlx.DB) *Adapter {
	a := &Adapter{
		db:        db,
		tableName: "casbin_rule",
	}
	a.ensureTable()
	return a
}

// NewAdapterFromOptions is the constructor for Adapter with existed connection
func NewAdapterFromOptions(opts *AdapterOptions) *Adapter {
	a := &Adapter{tableName: "casbin_rule"}
	if opts.TableName != "" {
		a.tableName = opts.TableName
	}
	if opts.DB != nil {
		a.db = opts.DB
	} else {
		db, err := sqlx.Connect(opts.DriverName, opts.DataSourceName)
		if err != nil {
			panic(err)
		}
		a.db = db
		runtime.SetFinalizer(a, finalizer)
	}
	a.ensureTable()
	return a
}

// LoadPolicy loads policy from database.
func (a *Adapter) LoadPolicy(model model.Model) error {
	var lines []CasbinRule
	err := a.db.Select(&lines, fmt.Sprintf("SELECT * FROM %s", a.tableName))
	if err != nil {
		return err
	}
	for _, line := range lines {
		loadPolicyLine(line, model)
	}
	return nil
}

// SavePolicy saves policy to database.
func (a *Adapter) SavePolicy(model model.Model) (err error) {
	a.dropTable()
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			err = a.insertPolicyLine(&line)
			if err != nil {
				return
			}
		}
	}
	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			err = a.insertPolicyLine(&line)
			if err != nil {
				return
			}
		}
	}
	return
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) (err error) {
	line := savePolicyLine(ptype, rule)
	err = a.insertPolicyLine(&line)
	if err != nil {
		return
	}
	return err
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) (err error) {
	line := savePolicyLine(ptype, rule)
	err = a.deletePolicyLine(&line)
	if err != nil {
		return
	}
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) (err error) {
	line := CasbinRule{}
	line.PType = ptype
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0.String = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1.String = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2.String = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3.String = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4.String = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5.String = fieldValues[5-fieldIndex]
	}
	err = a.rawDelete(&line)
	if err != nil {
		return
	}
	return
}

func (a *Adapter) rawDelete(line *CasbinRule) (err error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE ptype = :ptype", a.tableName)
	if line.V0.Valid {
		query += " AND v0 = :v0"
	}
	if line.V1.Valid {
		query += " AND v1 = :v1"
	}
	if line.V2.Valid {
		query += " AND v2 = :v2"
	}
	if line.V3.Valid {
		query += " AND v3 = :v3"
	}
	if line.V4.Valid {
		query += " AND v4 = :v4"
	}
	if line.V5.Valid {
		query += " AND v5 = :v5"
	}
	_, err = a.db.NamedExec(query, line)
	if err != nil {
		return
	}
	return
}

// IsFiltered returns true if the loaded policy has been filtered.
func (a *Adapter) IsFiltered() bool {
	return a.isFiltered
}
