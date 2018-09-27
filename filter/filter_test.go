package filter

import (
	"testing"

	"github.com/gocraft/dbr"

	"github.com/gocraft/dbr/dialect"
)

func TestFilter(t *testing.T) {
	cases := []struct {
		src, want, wantErr string
	}{
		// single field OK
		{"field=value", "`field` = 'value'", ""},
		{"field!=value", "`field` != 'value'", ""},
		{"field@NULL", "`field` IS NULL", ""},
		{"field@NOTNULL", "`field` IS NOT NULL", ""},
		{"field@LIKEvalue", "`field` LIKE '%value%'", ""},
		{"field@NOTLIKEvalue", "`field` NOT LIKE '%value%'", ""},
		{"field@STARTSWITHvalue", "`field` LIKE 'value%'", ""},
		{"field@INa,b,c", "`field` IN ('a','b','c')", ""},
		{"field@ENDSWITHvalue", "`field` LIKE '%value'", ""},
		{"field@NOTINa,b,c", "`field` NOT IN ('a','b','c')", ""},
		{"field=", "`field` = ''", ""},
		{"field@NOTIN", "`field` != ''", ""},

		// single field BAD
		{"1field=", "FALSE", "expected field name at 0"},
		{"=test", "FALSE", "expected field name at 0"},
		{"field@NULLblah", "FALSE", "expected ^ or ^OR at 10"},
		{"field@NOTNULLblah", "FALSE", "expected ^ or ^OR at 13"},
		{"field@NOTNULL^", "FALSE", "expected field name at 14"},

		// multi field OK
		{"field1=^field2=", "`field1` = '' AND `field2` = ''", ""},
		{"field1=123^ORfield2=456", "`field1` = '123' OR `field2` = '456'", ""},
		{"field1@NOTNULL^ORfield2@NOTNULL", "`field1` IS NOT NULL OR `field2` IS NOT NULL", ""},
		{"field1@NOTNULL^field2@NOTNULL", "`field1` IS NOT NULL AND `field2` IS NOT NULL", ""},

		// long
		{
			"field1@NOTNULL^field2@NOTNULL^ORfield3@STARTSWITHswag",
			"`field1` IS NOT NULL AND `field2` IS NOT NULL OR `field3` LIKE 'swag%'",
			"",
		},
	}

	for i, c := range cases {
		result, err := Parse(c.src)
		if err != nil {
			if c.wantErr == "" {
				t.Errorf("case %d: unexpected error %q", i, err)
			} else if err.Error() != c.wantErr {
				t.Errorf("case %d: unexpected error %q (expected %q)", i, err, c.wantErr)
			}
		} else {
			if c.wantErr != "" {
				t.Errorf("case %d: expected error %q, got none", i, c.wantErr)
			}
		}

		buf := dbr.NewBuffer()
		err = result.Build(dialect.MySQL, buf)
		if err != nil {
			t.Errorf("case %d: %v", i, err)
		}

		if real, err := dbr.InterpolateForDialect(buf.String(), buf.Value(), dialect.MySQL); err != nil || real != c.want {
			if err != nil {
				t.Errorf("case %d: %s", i, err)
			} else {
				t.Errorf("case %d: expected %s, got %s", i, c.want, real)
			}
		}
	}
}
