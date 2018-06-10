package database

import "database/sql"

type (
	BeforeInserter interface{ BeforeInsert() error }
	AfterInserter  interface{ AfterInsert(sql.Result) }

	BeforeUpdater interface{ BeforeUpdate() (err error) }
	AfterUpdater  interface{ AfterUpdate(sql.Result) }

	BeforeDeleter interface{ BeforeDelete() (err error) }
	AfterDeleter  interface{ AfterDelete(sql.Result) }
)
