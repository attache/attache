package main

import (
	"errors"
	"strings"

	"github.com/iancoleman/strcase"
)

type Model struct {
	DefaultKey     bool
	KeyStructField string
	KeyColumn      string

	Fields []Field
}

type Field struct {
	Column      string
	Key         bool
	StructField string
	Type        string

	NoInsert, NoUpdate, NoSelect bool
}

func parseFields(defs []string) (list []Field, hasKey bool, err error) {
	list = []Field{} // ensure list is non-nil

	for _, fdef := range defs {
		var (
			parts = strings.SplitN(fdef, ":", 3)
			f     = Field{}
		)

		f.StructField = strcase.ToCamel(parts[0])
		f.Column = strcase.ToSnake(f.StructField)

		if len(parts) > 1 && parts[1] != "" {
			f.Type = parts[1]
		} else {
			f.Type = "string"
		}

		if len(parts) > 2 {
			f.Key = strings.Contains(parts[2], "key")
			if hasKey && f.Key {
				return nil, false, errors.New("cannot specify multiple keys")
			}

			hasKey = hasKey || f.Key

			f.NoInsert = strings.Contains(parts[2], "noinsert")
			f.NoUpdate = strings.Contains(parts[2], "noupdate")
			f.NoSelect = strings.Contains(parts[2], "noselect")

			// enforce certain flags for keys
			if f.Key {
				f.NoUpdate = true
			}
		}

		list = append(list, f)
	}

	return
}

func buildModel(defs []string) (*Model, error) {
	m := new(Model)

	var err error
	var hasKey bool
	m.Fields, hasKey, err = parseFields(defs)
	if err != nil {
		return nil, err
	}

	if !hasKey {
		m.DefaultKey = true
		m.KeyColumn = "id"
		m.KeyStructField = "ID"

		m.Fields = append(
			append(
				make([]Field, 0, len(m.Fields)+1),
				Field{
					Column:      m.KeyColumn,
					StructField: m.KeyStructField,
					Type:        "int64",
					Key:         true,
					NoUpdate:    true,
					NoInsert:    true,
				},
			),
			m.Fields...,
		)
	} else {
		for _, f := range m.Fields {
			if f.Key {
				m.KeyStructField = f.StructField
				m.KeyColumn = f.Column
				break // only one key, quit once we find it
			}
		}
	}

	return m, nil
}
