package main

import "strings"

func opsMatch(item OpsRecord, query OpsQuery) bool {
	if query.Subject != "" && !strings.Contains(strings.ToLower(item.Subject), strings.ToLower(query.Subject)) {
		return false
	}
	if query.Status != "" && item.Status != query.Status {
		return false
	}
	if query.Priority != "" && item.Priority != query.Priority {
		return false
	}
	if query.Owner != "" && item.Owner != query.Owner {
		return false
	}
	return true
}
const opsDefaultPageSize = 20

// opsQueryDefaults fills in sensible defaults and clamps out-of-range values so
// that downstream pagination never panics on empty, zero, or negative inputs.
func opsQueryDefaults(q OpsQuery) OpsQuery {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = opsDefaultPageSize
	}
	return q
}

// opsBounds computes the [start, end) slice window for a page over total items.
// The returned values are always within [0, total] and start <= end, so the
// expression items[start:end] can never panic regardless of the page requested.
func opsBounds(total, page, size int) (int, int) {
	q := opsQueryDefaults(OpsQuery{Page: page, PageSize: size})
	start := (q.Page - 1) * q.PageSize
	if start < 0 {
		start = 0
	}
	if start > total {
		start = total
	}
	end := start + q.PageSize
	if end > total {
		end = total
	}
	return start, end
}
func opsPageCount(total, size int) int {
	if size < 1 || total == 0 {
		return 0
	}
	return (total + size - 1) / size
}
func opsQueryKey(q OpsQuery) string {
	return strings.Join([]string{q.Subject, string(q.Status), string(q.Priority), q.Owner}, "|")
}
// opsClonePage returns a fully independent copy of p: the Items slice and each
// record's Labels map are duplicated, so mutations to the clone never bleed back
// into the original (or vice versa).
func opsClonePage(p OpsPage) OpsPage {
	out := p
	out.Items = make([]OpsRecord, len(p.Items))
	for i, item := range p.Items {
		out.Items[i] = item.Clone()
	}
	return out
}
func opsLastID(p OpsPage) string {
	if len(p.Items) == 0 {
		return ""
	}
	return p.Items[len(p.Items)-1].ID
}
