package attache

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// Storeable is the interface implemented by any type that can be
// stored/retrieved via database/sql
type Storeable interface {
	// Table returns the name of the table that represents
	// the Storeable object in the database
	Table() string

	// Insert returns the columns to be inserted, and the values
	// to be inserted in those columns for the Storeable object
	Insert() (columns []string, values []interface{})

	// Update returns the columns to be updated, and the
	// updated values for those columns for the Storeable object
	Update() (columns []string, values []interface{})

	// Select returns the columns to be selected,
	// as well as pointers to values that will be used
	// to store the values retrieved from the database
	Select() (columns []string, into []interface{})

	// KeyColumns returns the columns composing the primary key
	// for the Storeable object
	KeyColumns() []string

	// KeyValues returns the values representing the primary key
	// values for the Storeable object
	KeyValues() []interface{}
}

type (
	// BeforeInserter can be implemented by a Storeable type
	// to provide a callback before insertion. If a non-nil error
	// is returned, the insert operation is aborted
	BeforeInserter interface{ BeforeInsert() error }

	// AfterInserter can be implemented by a Storeable type
	// to provide a callback after insertion. The callback
	// has access to the returned sql.Result
	AfterInserter interface{ AfterInsert(sql.Result) }
)

type (
	// BeforeUpdater can be implemented by a Storeable type
	// to provide a callback before an update. If a non-nil
	// error is returned, the update operation is aborted
	BeforeUpdater interface{ BeforeUpdate() (err error) }

	// AfterUpdater can be implemented by a Storeable type
	// to provide a callback after an update. The callback
	// has access to the returned sql.Result
	AfterUpdater interface{ AfterUpdate(sql.Result) }
)

type (
	// BeforeDeleter can be implemented by a Storeable type
	// to provide a callback before a deletion. If a non-nil
	// error is returned, the deletion is aborted
	BeforeDeleter interface{ BeforeDelete() (err error) }

	// AfterDeleter can be implemented by a Storeable type
	// to provide a callback after a deletion. The callback
	// has access to the returned sql.Result
	AfterDeleter interface{ AfterDelete(sql.Result) }
)

var (
	tStoreable = reflect.TypeOf((*Storeable)(nil)).Elem()
)

// A DB provides methods for querying the database, executing sql
// operations, and storing / deleting / querying Storeable objects
type DB interface {
	// Query is identical to DB.Query from database/sql
	Query(query string, args ...interface{}) (*sql.Rows, error)
	// Exec is identical to DB.Exec from database/sql
	Exec(query string, args ...interface{}) (sql.Result, error)

	// Insert inserts the Storeable object into the database.
	// Any error encountered is returned
	Insert(s Storeable) error

	// Update updates the Storeable object in the database.
	// Any error encountered is returned
	Update(s Storeable) error

	// Delete removes the Storeable object from the database.
	// Any error encountered is returned
	Delete(s Storeable) error

	// All queries all objects in the database represented
	// by the concrete type returned by newFn.
	// Any error encountered is returned
	All(newFn func() Storeable) ([]Storeable, error)

	// Where queries all objects in the database represented by
	// the concrete type returned by newFn that match the
	// given where clause.
	// Any error encountered is returned
	Where(newFn func() Storeable, where string) ([]Storeable, error)

	// Find locates the object in the database represented by
	// into's concrete type and by the key value(s) provided in args
	Find(into Storeable, args ...interface{}) error

	// Find locates the object in the database represented by
	// into's concrete type and by the given field-value combination
	FindBy(into Storeable, field string, val interface{}) error

	private()
}

type db struct {
	conn *sql.DB
}

func (d db) private() {}

func (d db) Exec(query string, args ...interface{}) (sql.Result, error) {
	sqlLog(query, args)
	return d.conn.Exec(query, args...)
}

func (d db) Query(query string, args ...interface{}) (*sql.Rows, error) {
	sqlLog(query, args)
	return d.conn.Query(query, args...)
}

func (d db) Insert(s Storeable) error {
	if impl, ok := s.(BeforeInserter); ok {
		if err := impl.BeforeInsert(); err != nil {
			return err
		}
	}

	cols, vals := s.Insert()

	query := new(bytes.Buffer)
	fmt.Fprintf(query, "INSERT INTO %s (", s.Table())
	for i, name := range cols {
		if i != 0 {
			query.WriteString(", ")
		}

		fmt.Fprintf(query, "%s", name)
	}

	query.WriteString(") VALUES (")
	for i := range vals {
		if i != 0 {
			query.WriteString(", ")
		}
		query.WriteString("?")
	}
	query.WriteString(")")

	result, err := d.Exec(query.String(), vals...)
	if err != nil {
		return err
	}

	if impl, ok := s.(AfterInserter); ok {
		impl.AfterInsert(result)
	}

	return nil
}

func (d db) Update(s Storeable) error {
	if impl, ok := s.(BeforeUpdater); ok {
		if err := impl.BeforeUpdate(); err != nil {
			return err
		}
	}

	cols, vals := s.Update()
	query := new(bytes.Buffer)
	fmt.Fprintf(query, "UPDATE %s SET ", s.Table())
	for i, name := range cols {
		if i != 0 {
			query.WriteString(", ")
		}

		fmt.Fprintf(query, "%s=?", name)
	}
	query.WriteString(" WHERE ")
	for i, name := range s.KeyColumns() {
		if i != 0 {
			query.WriteString(" AND ")
		}

		fmt.Fprintf(query, "%s=?", name)
	}

	result, err := d.Exec(query.String(), append(vals, s.KeyValues()...)...)
	if err != nil {
		return err
	}

	if impl, ok := s.(AfterUpdater); ok {
		impl.AfterUpdate(result)
	}

	return nil
}

func (d db) Delete(s Storeable) error {
	if impl, ok := s.(BeforeDeleter); ok {
		if err := impl.BeforeDelete(); err != nil {
			return err
		}
	}

	query := new(bytes.Buffer)
	fmt.Fprintf(query, "DELETE FROM %s WHERE ", s.Table())

	for i, name := range s.KeyColumns() {
		if i != 0 {
			query.WriteString(" AND ")
		}

		fmt.Fprintf(query, "%s=?", name)
	}

	result, err := d.Exec(query.String(), s.KeyValues()...)
	if err != nil {
		return err
	}

	if impl, ok := s.(AfterDeleter); ok {
		impl.AfterDelete(result)
	}

	return nil
}

func (d db) All(newFn func() Storeable) ([]Storeable, error) {
	storeable := newFn()
	cols, _ := storeable.Select()
	result := make([]Storeable, 0, 64)
	query := selectQuery(cols, storeable.Table(), nil, false, 0)
	rows, err := d.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		into := newFn()
		_, targs := into.Select()
		if err := rows.Scan(targs...); err != nil {
			return nil, err
		}
		result = append(result, into)
	}

	if len(result) == 0 {
		return nil, sql.ErrNoRows
	}

	return result, nil
}

func (d db) Where(newFn func() Storeable, where string) ([]Storeable, error) {
	if len(where) == 0 {
		return d.All(newFn)
	}

	storeable := newFn()
	cols, _ := storeable.Select()
	result := make([]Storeable, 0, 64)
	query := selectQueryUnsafe(cols, storeable.Table(), where, 0)
	rows, err := d.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		into := newFn()
		_, targs := into.Select()
		if err := rows.Scan(targs...); err != nil {
			return nil, err
		}
		result = append(result, into)
	}

	if len(result) == 0 {
		return nil, sql.ErrNoRows
	}

	return result, nil
}

func (d db) Find(into Storeable, args ...interface{}) error {
	cols, targets := into.Select()
	query := selectQuery(cols, into.Table(), into.KeyColumns(), false, 1)
	rows, err := d.Query(query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}

	return rows.Scan(targets...)
}

func (d db) FindBy(into Storeable, field string, val interface{}) error {
	cols, targets := into.Select()
	query := selectQuery(cols, into.Table(), []string{field}, false, 1)
	rows, err := d.Query(query, val)
	if err != nil {
		return err
	}

	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}

	return rows.Scan(targets...)
}

func selectQuery(cols []string, table string, searchFields []string, or bool, limit int) string {
	query := new(strings.Builder)
	query.WriteString("SELECT ")
	for i, name := range cols {
		if i != 0 {
			query.WriteString(", ")
		}

		query.WriteString(name)
	}

	query.WriteString(" FROM ")
	query.WriteString(table)
	if len(searchFields) > 0 {
		query.WriteString(" WHERE ")
		for i, field := range searchFields {
			if i != 0 {
				if or {
					query.WriteString(" OR ")
				} else {
					query.WriteString(" AND ")
				}
			}

			fmt.Fprintf(query, "%s=?", field)
		}
	}

	if limit > 0 {
		fmt.Fprintf(query, " LIMIT %d", limit)
	}

	query.WriteByte(';')
	return query.String()
}

func selectQueryUnsafe(cols []string, table, where string, limit int) string {
	query := new(strings.Builder)
	query.WriteString("SELECT ")
	for i, name := range cols {
		if i != 0 {
			query.WriteString(", ")
		}

		query.WriteString(name)
	}

	query.WriteString(" FROM ")
	query.WriteString(table)
	query.WriteString(" WHERE ")
	query.WriteString(where)

	if limit > 0 {
		fmt.Fprintf(query, " LIMIT %d", limit)
	}

	query.WriteByte(';')
	return query.String()
}

func sqlLog(query string, args []interface{}) {
	log.Println("SQL:", query, args)
}
