package cmd_gen

import (
	"strings"

	"github.com/iancoleman/strcase"
)

type Model struct {
	Name  string
	Table string

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

func parseFields(defs []string) (list []Field, hasKey bool) {
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
			f.Key = strings.Index(parts[2], "key") > -1
			if hasKey && f.Key {
				panic(ArgumentError{
					Cause: nil,
					Issue: "cannot specify multiple keys",
				})
			}

			hasKey = hasKey || f.Key

			f.NoInsert = strings.Index(parts[2], "noinsert") > -1
			f.NoUpdate = strings.Index(parts[2], "noupdate") > -1
			f.NoSelect = strings.Index(parts[2], "noselect") > -1

			// enforce certain flags for keys
			if f.Key {
				f.NoUpdate = true
			}
		}

		list = append(list, f)
	}

	return
}

func buildModel(name string, table string, defs []string) *Model {
	m := new(Model)

	m.Name = strcase.ToCamel(name)
	if table == "" {
		m.Table = strcase.ToSnake(m.Name)
	} else {
		m.Table = table
	}

	hasKey := false
	m.Fields, hasKey = parseFields(defs)

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

	return m
}
