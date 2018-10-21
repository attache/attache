package attache

import (
	"database/sql"
	"fmt"

	"github.com/attache/attache/filter"
	"github.com/gocraft/dbr"
)

var ErrRecordNotFound = dbr.ErrNotFound

// Record is the interface implemented by any type that can be
// stored/retrieved via database/sql
type Record interface {
	// Table returns the name of the table that represents
	// the Record object in the database
	Table() string

	// Key returns the columns and values that compose the primary key
	// for the Record object
	Key() ([]string, []interface{})

	// InsertColumns returns the columns and values to be used for a default INSERT operation
	Insert() (columns []string, values []interface{})

	// UpdateColumns returns the columns and values to be used for a default UPDATE operation
	Update() (columns []string, values []interface{})
}

type (
	// BeforeInserter can be implemented by a Record type
	// to provide a callback before insertion. If a non-nil error
	// is returned, the insert operation is aborted
	BeforeInserter interface{ BeforeInsert() error }

	// AfterInserter can be implemented by a Record type
	// to provide a callback after insertion. The callback
	// has access to the returned sql.Result
	AfterInserter interface{ AfterInsert(sql.Result) }
)

type (
	// BeforeUpdater can be implemented by a Record type
	// to provide a callback before an update. If a non-nil
	// error is returned, the update operation is aborted
	BeforeUpdater interface{ BeforeUpdate() (err error) }

	// AfterUpdater can be implemented by a Record type
	// to provide a callback after an update. The callback
	// has access to the returned sql.Result
	AfterUpdater interface{ AfterUpdate(sql.Result) }
)

type (
	// BeforeDeleter can be implemented by a Record type
	// to provide a callback before a deletion. If a non-nil
	// error is returned, the deletion is aborted
	BeforeDeleter interface{ BeforeDelete() (err error) }

	// AfterDeleter can be implemented by a Record type
	// to provide a callback after a deletion. The callback
	// has access to the returned sql.Result
	AfterDeleter interface{ AfterDelete(sql.Result) }
)

// DBRunner is the set of functions provided by both DB and TX
type DBRunner interface {
	// Insert inserts the Record into the database
	Insert(Record) error

	// Update updates the Record in the database
	Update(Record) error

	// Delete deletes the Record from the database
	Delete(Record) error

	// All reutrns all Records of the specified type
	All(typ Record) ([]Record, error)

	// Get fetches a single record whose key values match those provided into the target Record
	Get(into Record, keyVals ...interface{}) error

	// GetBy fetches a single record matching the where condition into the target Record
	GetBy(into Record, where string, args ...interface{}) error

	// Where finds all records of the specified type matching the where query
	Where(typ Record, where string, args ...interface{}) ([]Record, error)

	// WhereFilter finds all records of the specified type matching the filterString
	WhereFilter(typ Record, filterString string) ([]Record, error)
}

var (
	_ DBRunner = DB{}
	_ DBRunner = TX{}
)

type DB struct{ s *dbr.Session }

func (db DB) Raw() *dbr.Session                             { return db.s }
func (db DB) Insert(r Record) error                         { return doInsert(db.s, r) }
func (db DB) Update(r Record) error                         { return doUpdate(db.s, r) }
func (db DB) Delete(r Record) error                         { return doDelete(db.s, r) }
func (db DB) All(typ Record) ([]Record, error)              { return doAll(db.s, typ) }
func (db DB) Get(into Record, keyVals ...interface{}) error { return doGet(db.s, into, keyVals) }
func (db DB) GetBy(into Record, where string, args ...interface{}) error {
	return doGetBy(db.s, into, where, args)
}
func (db DB) Where(typ Record, where string, args ...interface{}) ([]Record, error) {
	return doWhere(db.s, typ, where, args)
}
func (db DB) WhereFilter(typ Record, where string) ([]Record, error) {
	return doWhereFilter(db.s, typ, where)
}

func (db DB) Tx(block func(tx TX) error) error {
	tx, err := db.s.Begin()
	if err != nil {
		return err
	}

	defer tx.RollbackUnlessCommitted()
	if err := block(TX{tx}); err != nil {
		return err
	}
	return tx.Commit()
}

type TX struct{ s *dbr.Tx }

func (db TX) Raw() *dbr.Tx                                  { return db.s }
func (db TX) Insert(r Record) error                         { return doInsert(db.s, r) }
func (db TX) Update(r Record) error                         { return doUpdate(db.s, r) }
func (db TX) Delete(r Record) error                         { return doDelete(db.s, r) }
func (db TX) All(typ Record) ([]Record, error)              { return doAll(db.s, typ) }
func (db TX) Get(into Record, keyVals ...interface{}) error { return doGet(db.s, into, keyVals) }
func (db TX) GetBy(into Record, where string, args ...interface{}) error {
	return doGetBy(db.s, into, where, args)
}
func (db TX) Where(typ Record, where string, args ...interface{}) ([]Record, error) {
	return doWhere(db.s, typ, where, args)
}
func (db TX) WhereFilter(typ Record, where string) ([]Record, error) {
	return doWhereFilter(db.s, typ, where)
}

// general implementations

type (
	dbrInserter interface {
		InsertInto(string) *dbr.InsertStmt
	}

	dbrUpdater interface {
		Update(string) *dbr.UpdateStmt
	}

	dbrDeleter interface {
		DeleteFrom(string) *dbr.DeleteStmt
	}

	dbrSelecter interface {
		Select(...string) *dbr.SelectStmt
	}
)

func doInsert(s dbrInserter, r Record) error {
	if impl, ok := r.(BeforeInserter); ok {
		if err := impl.BeforeInsert(); err != nil {
			return err
		}
	}

	ic, iv := r.Insert()
	result, err := s.InsertInto(r.Table()).Columns(ic...).Values(iv...).Exec()
	if err != nil {
		return err
	}

	if impl, ok := r.(AfterInserter); ok {
		impl.AfterInsert(result)
	}

	return nil
}

func doUpdate(s dbrUpdater, r Record) error {
	if impl, ok := r.(BeforeUpdater); ok {
		if err := impl.BeforeUpdate(); err != nil {
			return err
		}
	}

	uc, uv := r.Update()
	query := s.Update(r.Table())
	for i := 0; i < len(uc); i++ {
		query.Set(uc[i], uv[i])
	}

	kc, kv := r.Key()
	for i := 0; i < len(kc); i++ {
		query.Where(dbr.Eq(kc[i], kv[i]))
	}

	result, err := query.Exec()
	if err != nil {
		return err
	}

	if impl, ok := r.(AfterUpdater); ok {
		impl.AfterUpdate(result)
	}

	return nil
}

func doDelete(s dbrDeleter, r Record) error {
	if impl, ok := r.(BeforeDeleter); ok {
		if err := impl.BeforeDelete(); err != nil {
			return err
		}
	}

	kc, kv := r.Key()
	query := s.DeleteFrom(r.Table())
	for i := 0; i < len(kc); i++ {
		query.Where(dbr.Eq(kc[i], kv[i]))
	}

	result, err := query.Exec()
	if err != nil {
		return err
	}

	if impl, ok := r.(AfterDeleter); ok {
		impl.AfterDelete(result)
	}

	return nil
}

func doAll(s dbrSelecter, typ Record) ([]Record, error) {
	result := []Record{}
	_, err := s.Select("*").From(typ.Table()).Load(dbr.InterfaceLoader(&result, typ))
	return result, err
}

func doGet(s dbrSelecter, into Record, keyVals []interface{}) error {
	kc, _ := into.Key()
	if got, want := len(keyVals), len(kc); got != want {
		return fmt.Errorf("expected %d values in key, got %d", want, got)
	}

	query := s.Select("*").From(into.Table())
	for i := 0; i < len(kc); i++ {
		query.Where(dbr.Eq(kc[i], keyVals[i]))
	}

	return query.LoadOne(into)
}

func doGetBy(s dbrSelecter, into Record, where string, args []interface{}) error {
	return s.Select("*").From(into.Table()).Where(where, args...).LoadOne(into)
}

func doWhere(s dbrSelecter, typ Record, where string, args []interface{}) ([]Record, error) {
	result := []Record{}
	_, err := s.Select("*").From(typ.Table()).Where(where, args...).Load(dbr.InterfaceLoader(&result, typ))
	return result, err
}

func doWhereFilter(s dbrSelecter, typ Record, where string) ([]Record, error) {
	result := []Record{}
	cond, err := filter.Parse(where)
	if err != nil {
		return nil, err
	}
	_, err = s.Select("*").From(typ.Table()).Where(cond).Load(dbr.InterfaceLoader(&result, typ))
	return result, err
}
