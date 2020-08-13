package dg

import "strings"

// QB is a simple query builder.
type QB struct {
	query   string
	filters []string
}

// NewQB creates a new query builder.
func NewQB(query string) *QB {
	return &QB{query: query}
}

// Filter adds a filter to the query builder.
func (q *QB) Filter(f string) {
	q.filters = append(q.filters, f)
}

// Query returns the query with all filters applied.
func (q *QB) Query() string {
	var f1 string
	var f2 string
	if len(q.filters) > 0 {
		f := strings.Join(q.filters, " AND ")
		f1 = "@filter(" + f + ")"
		f2 = "AND " + f
	}
	query := strings.ReplaceAll(q.query, "<FILTERS>", f1)
	query = strings.ReplaceAll(query, "<AND_FILTERS>", f2)
	return query
}
