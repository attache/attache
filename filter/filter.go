package filter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocraft/dbr"
)

type Relation int

const (
	ILLEGAL Relation = iota
	EQ
	NEQ
	GT
	LT
	GEQ
	LEQ
	NULL
	NOTNULL
	LIKE
	NOTLIKE
	STARTSWITH
	ENDSWITH
	IN
	NOTIN
)

func add(cond dbr.Builder, isOrCondition bool, field string, rel Relation, val interface{}) dbr.Builder {
	var cond2 dbr.Builder
	switch rel {
	case EQ, IN:
		cond2 = dbr.Eq(field, val)
	case NEQ, NOTIN:
		cond2 = dbr.Neq(field, val)
	case GT:
		cond2 = dbr.Gt(field, val)
	case LT:
		cond2 = dbr.Lt(field, val)
	case GEQ:
		cond2 = dbr.Gte(field, val)
	case LEQ:
		cond2 = dbr.Lte(field, val)
	case NULL:
		cond2 = dbr.Eq(field, nil)
	case NOTNULL:
		cond2 = dbr.Neq(field, nil)
	case LIKE:
		cond2 = dbr.Like(field, fmt.Sprint("%", val, "%"))
	case NOTLIKE:
		cond2 = dbr.NotLike(field, fmt.Sprint("%", val, "%"))
	case STARTSWITH:
		cond2 = dbr.Like(field, fmt.Sprint(val, "%"))
	case ENDSWITH:
		cond2 = dbr.Like(field, fmt.Sprint("%", val))
	}

	if cond == nil {
		return cond2
	}

	if isOrCondition {
		return dbr.BuildFunc(func(d dbr.Dialect, buf dbr.Buffer) error {
			if err := cond.Build(d, buf); err != nil {
				return err
			}
			buf.WriteString(" OR ")
			return cond2.Build(d, buf)
		})
	}

	return dbr.BuildFunc(func(d dbr.Dialect, buf dbr.Buffer) error {
		if err := cond.Build(d, buf); err != nil {
			return err
		}
		buf.WriteString(" AND ")
		return cond2.Build(d, buf)
	})
}

var strToRel = map[string]Relation{
	"=":           EQ,
	"!=":          NEQ,
	">":           GT,
	"<":           LT,
	">=":          GEQ,
	"<=":          LEQ,
	"@NULL":       NULL,
	"@NOTNULL":    NOTNULL,
	"@LIKE":       LIKE,
	"@NOTLIKE":    NOTLIKE,
	"@STARTSWITH": STARTSWITH,
	"@ENDSWITH":   ENDSWITH,
	"@IN":         IN,
	"@NOTIN":      NOTIN,
}

var strToBoolOp = map[string]string{
	"^":   "AND",
	"^OR": "OR",
}

func Parse(filter string) (dbr.Builder, error) {
	if len(filter) == 0 {
		return dbr.Expr(""), nil
	}

	var (
		cond                       dbr.Builder
		op                         string
		field, rel, vals, pos, err = parseClause(filter, 0)
	)

	if err != nil {
		goto parseFailure
	}

	cond = add(cond, false, field, rel, firstOrAll(vals))

	for pos < len(filter) {
		if op, err = parseBoolOp(filter, pos); err != nil {
			goto parseFailure
		}

		if op == "AND" {
			pos += 1
		} else {
			pos += 3
		}
		if field, rel, vals, pos, err = parseClause(filter, pos); err != nil {
			goto parseFailure
		}

		cond = add(cond, op == "OR", field, rel, firstOrAll(vals))
	}

	return cond, nil

parseFailure:
	return dbr.Expr("FALSE"), err
}

func parseClause(filter string, start int) (field string, rel Relation, vals []string, end int, err error) {
	end = start
	if field, err = parseIdent(filter, end); err != nil {
		return
	}
	end += len(field)

	var relStr string
	if relStr, err = parseRelation(filter, end); err != nil {
		return
	}
	end += len(relStr)

	rel = strToRel[relStr]
	if rel == NULL || rel == NOTNULL {
		return
	}

	var read int
	vals, read = parseValues(filter, end)
	return field, rel, vals, end + read, nil
}

var (
	fIdent    = regexp.MustCompile(`^[a-zA-Z_][a-z_0-9]*`)
	fRelation = regexp.MustCompile(`^(?:=|!=|@NULL|@NOTNULL|@LIKE|@NOTLIKE|@STARTSWITH|@ENDSWITH|@IN|@NOTIN)`)
	fBoolOp   = regexp.MustCompile(`^\^(?:OR)?`)
	fValues   = regexp.MustCompile(`^[^^]*`)
)

func parseIdent(src string, pos int) (match string, err error) {
	//fmt.Println("finding ident in", src[pos:])
	match = fIdent.FindString(src[pos:])
	if match == "" {
		err = fmt.Errorf("expected field name at %d", pos)
	}
	return match, err
}

func parseRelation(src string, pos int) (match string, err error) {
	//fmt.Println("finding relation in", src[pos:])
	match = fRelation.FindString(src[pos:])
	if match == "" {
		err = fmt.Errorf("expected relation at %d", pos)
	}
	return match, err
}

func parseBoolOp(src string, pos int) (string, error) {
	//fmt.Println("finding bool op in", src[pos:])
	match := fBoolOp.FindString(src[pos:])
	if match == "" {
		return "", fmt.Errorf("expected ^ or ^OR at %d", pos)
	}
	return strToBoolOp[match], nil
}

func parseValues(src string, pos int) (list []string, end int) {
	//fmt.Println("finding values in", src[pos:])
	s := fValues.FindString(src[pos:])
	return strings.Split(s, ","), len(s)
}

func firstOrAll(values []string) interface{} {
	if len(values) == 0 {
		return nil
	}

	if len(values) == 1 {
		return values[0]
	}

	return values
}
