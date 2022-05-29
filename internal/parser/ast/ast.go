package ast

import (
	"fmt"
	"network_go/internal/parser/token"
	"strings"
)

type Attrib interface{}

type Field struct {
	Bucket   string
	Operator string
	Matter   string
}

func (f *Field) Match(compare string) bool {
	switch f.Operator {
	case "=":
		return f.Matter == compare
	case "!=":
		return f.Matter != compare
	case "~":
		return strings.Contains(compare, f.Matter)
	case "!~":
		return !strings.Contains(compare, f.Matter)
	}
	return false
}

func (f *Field) String() string {
	return fmt.Sprintf(`%s %s "%s"`, f.Bucket, f.Operator, f.Matter)
}

func NewField(bucket Attrib, operator Attrib, matter Attrib) (*Field, error) {
	matter_s := string(matter.(*token.Token).Lit)
	matter_s = strings.Trim(matter_s, "\"")
	matter_s = strings.ToLower(matter_s)
	return &Field{string(bucket.(*token.Token).Lit), string(operator.(*token.Token).Lit), matter_s}, nil
}

type Query struct {
	Field    *Field
	AndQuery *Query
	OrQuery  *Query
}

func (q Query) String() string {
	// No additional queries - Field only
	if q.AndQuery == nil && q.OrQuery == nil {
		return q.Field.String()
	}

	// Both Queries
	if q.AndQuery != nil && q.OrQuery != nil {
		return fmt.Sprintf("(%s & %s) | (%s)", q.Field, q.AndQuery, q.OrQuery)
	}

	// And Query
	if q.AndQuery != nil {
		return fmt.Sprintf("(%s & %s)", q.Field, q.AndQuery)
	}

	// Or Query
	if q.OrQuery != nil {
		return fmt.Sprintf("(%s | %s)", q.Field, q.OrQuery)
	}

	return ""
}

func (q Query) NumFields() int64 {
	if q.Field == nil {
		return 0
	}

	if q.AndQuery == nil && q.OrQuery == nil {
		return 1
	}

	// Both Queries
	if q.AndQuery != nil && q.OrQuery != nil {
		return 1 + q.AndQuery.NumFields() + q.OrQuery.NumFields()
	}

	// And Query
	if q.AndQuery != nil {
		return 1 + q.AndQuery.NumFields()
	}

	// Or Query
	if q.OrQuery != nil {
		return 1 + q.OrQuery.NumFields()
	}

	return 0
}

func (f Field) SearchString(input string) bool {
	input = strings.ToLower(input)
	switch f.Operator {
	case "=":
		return input == f.Matter
	case "~":
		return strings.Contains(input, f.Matter)
	case "!=":
		return input != f.Matter
	case "!~":
		return !strings.Contains(input, f.Matter)
	}
	return false
}

func (f Field) SearchStringSlice(input []string) bool {
	for _, s := range input {
		if f.SearchString(s) {
			return true
		}
	}
	return false
}

func NewFieldQuery(t Attrib) (*Query, error) {
	return &Query{Field: t.(*Field)}, nil
}

func NewAndQuery(t1 Attrib, t2 Attrib) (*Query, error) {
	fmt.Printf("%s - %s\n", t1, t2)
	orig := t1.(*Query)
	tmp := orig
	for tmp.AndQuery != nil {
		tmp = tmp.AndQuery
	}
	tmp.AndQuery = t2.(*Query)
	return orig, nil
}

func NewOrQuery(t1 Attrib, t2 Attrib) (*Query, error) {
	fmt.Printf("%s - %s\n", t1, t2)
	orig := t1.(*Query)
	tmp := orig
	for tmp.OrQuery != nil {
		tmp = tmp.OrQuery
	}
	tmp.OrQuery = t2.(*Query)
	return orig, nil
}
