package attache

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/gocraft/dbr"
	"github.com/mccolljr/attache/filter"
)

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

var tRecord = reflect.TypeOf((*Record)(nil)).Elem()

type DB struct {
	s *dbr.Session
}

func (db DB) Insert(r Record) error {
	if impl, ok := r.(BeforeInserter); ok {
		if err := impl.BeforeInsert(); err != nil {
			return err
		}
	}

	ic, iv := r.Insert()
	result, err := db.s.InsertInto(r.Table()).Columns(ic...).Values(iv...).Exec()
	if err != nil {
		return err
	}

	if impl, ok := r.(AfterInserter); ok {
		impl.AfterInsert(result)
	}

	return nil
}

func (db DB) Update(r Record) error {
	if impl, ok := r.(BeforeUpdater); ok {
		if err := impl.BeforeUpdate(); err != nil {
			return err
		}
	}

	uc, uv := r.Update()
	query := db.s.Update(r.Table())
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

func (db DB) Delete(r Record) error {
	if impl, ok := r.(BeforeDeleter); ok {
		if err := impl.BeforeDelete(); err != nil {
			return err
		}
	}

	kc, kv := r.Key()
	query := db.s.DeleteFrom(r.Table())
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

func (db DB) FindByKey(into Record, keyVals ...interface{}) error {
	kc, _ := into.Key()
	if got, want := len(keyVals), len(kc); got != want {
		return fmt.Errorf("expected %d values in key, got %d", want, got)
	}

	query := db.s.Select("*").From(into.Table())
	for i := 0; i < len(kc); i++ {
		query.Where(dbr.Eq(kc[i], keyVals[i]))
	}

	return query.LoadOne(into)
}

func (db DB) FindBy(into Record, where string, args ...interface{}) error {
	return db.s.Select("*").From(into.Table()).Where(where, args...).LoadOne(into)
}

func (db DB) FindAll(typ Record) ([]Record, error) {
	result := []Record{}
	_, err := db.s.Select("*").From(typ.Table()).Load(dbr.InterfaceLoader(&result, typ))
	return result, err
}

func (db DB) FindByFilter(typ Record, where string) ([]Record, error) {
	result := []Record{}
	cond, err := filter.Parse(where)
	if err != nil {
		return nil, err
	}
	_, err = db.s.Select("*").From(typ.Table()).Where(cond).Load(dbr.InterfaceLoader(&result, typ))
	return result, err
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

type TX struct {
	t *dbr.Tx
}

func (tx TX) Insert(r Record) error {
	if impl, ok := r.(BeforeInserter); ok {
		if err := impl.BeforeInsert(); err != nil {
			return err
		}
	}

	ic, iv := r.Insert()
	result, err := tx.t.InsertInto(r.Table()).Columns(ic...).Values(iv...).Exec()
	if err != nil {
		return err
	}

	if impl, ok := r.(AfterInserter); ok {
		impl.AfterInsert(result)
	}

	return nil
}

func (tx TX) Update(r Record) error {
	if impl, ok := r.(BeforeUpdater); ok {
		if err := impl.BeforeUpdate(); err != nil {
			return err
		}
	}

	uc, uv := r.Update()
	query := tx.t.Update(r.Table())
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

func (tx TX) Delete(r Record) error {
	if impl, ok := r.(BeforeDeleter); ok {
		if err := impl.BeforeDelete(); err != nil {
			return err
		}
	}

	kc, kv := r.Key()
	query := tx.t.DeleteFrom(r.Table())
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

func (tx TX) FindByKey(into Record, keyVals ...interface{}) error {
	kc, _ := into.Key()
	if got, want := len(keyVals), len(kc); got != want {
		return fmt.Errorf("expected %d values in key, got %d", want, got)
	}

	query := tx.t.Select("*").From(into.Table())
	for i := 0; i < len(kc); i++ {
		query.Where(dbr.Eq(kc[i], keyVals[i]))
	}

	return query.LoadOne(into)
}

func (tx TX) FindBy(into Record, where string, args ...interface{}) error {
	return tx.t.Select("*").From(into.Table()).Where(where, args...).LoadOne(into)
}

func (tx TX) FindAll(typ Record) ([]Record, error) {
	result := []Record{}
	_, err := tx.t.Select("*").From(typ.Table()).Load(dbr.InterfaceLoader(&result, typ))
	return result, err
}

func (tx TX) FindByFilter(typ Record, where string) ([]Record, error) {
	result := []Record{}
	cond, err := filter.Parse(where)
	if err != nil {
		return nil, err
	}
	_, err = tx.t.Select("*").From(typ.Table()).Where(cond).Load(dbr.InterfaceLoader(&result, typ))
	return result, err
}
